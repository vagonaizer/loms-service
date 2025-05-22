package loms

import (
	"context"
	"log"
	"time"

	loms "github.com/vagonaizer/loms/api/protos/gen/loms"
	"github.com/vagonaizer/loms/internal/domain/models"
	"github.com/vagonaizer/loms/internal/domain/ports"
)

type Service struct {
	loms.UnimplementedLOMSServer
	orderRepo ports.OrderRepository
	stockRepo ports.StockRepository
}

func NewService(orderRepo ports.OrderRepository, stockRepo ports.StockRepository) *Service {
	return &Service{
		orderRepo: orderRepo,
		stockRepo: stockRepo,
	}
}

func (s *Service) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	log.Printf("Received OrderCreate request for user %d with %d items", req.User, len(req.Items))

	order := &models.Order{
		UserID:    req.User,
		Status:    models.OrderStatusNew,
		Items:     make([]models.OrderItem, len(req.Items)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for i, item := range req.Items {
		order.Items[i] = models.OrderItem{
			SKU:   item.Sku,
			Count: item.Count,
		}
	}

	// Try to reserve items
	for _, item := range order.Items {
		log.Printf("Checking stock for SKU %d, requested count: %d", item.SKU, item.Count)
		stock, err := s.stockRepo.Get(ctx, item.SKU)
		if err != nil {
			log.Printf("Failed to get stock for SKU %d: %v", item.SKU, err)
			return nil, err
		}

		if stock == nil || stock.TotalCount-stock.Reserved < uint64(item.Count) {
			log.Printf("Insufficient stock for SKU %d: requested %d, available %d",
				item.SKU, item.Count, stock.TotalCount-stock.Reserved)
			order.Status = models.OrderStatusFailed
			if err := s.orderRepo.Create(ctx, order); err != nil {
				return nil, err
			}
			return &loms.OrderCreateResponse{OrderID: order.ID}, nil
		}

		log.Printf("Reserving %d items for SKU %d", item.Count, item.SKU)
		if err := s.stockRepo.Reserve(ctx, item.SKU, uint64(item.Count)); err != nil {
			log.Printf("Failed to reserve items for SKU %d: %v", item.SKU, err)
			return nil, err
		}
	}

	order.Status = models.OrderStatusAwaitingPayment
	if err := s.orderRepo.Create(ctx, order); err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, err
	}

	log.Printf("Successfully created order %d with status %s", order.ID, order.Status)
	return &loms.OrderCreateResponse{OrderID: order.ID}, nil
}

func (s *Service) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	log.Printf("Received OrderInfo request for order %d", req.OrderID)

	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		log.Printf("Failed to get order %d: %v", req.OrderID, err)
		return nil, err
	}

	if order == nil {
		log.Printf("Order %d not found", req.OrderID)
		return &loms.OrderInfoResponse{}, nil
	}

	items := make([]*loms.Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = &loms.Item{
			Sku:   item.SKU,
			Count: item.Count,
		}
	}

	log.Printf("Returning info for order %d: status=%s, user=%d, items=%v",
		order.ID, order.Status, order.UserID, items)
	return &loms.OrderInfoResponse{
		Status: string(order.Status),
		User:   order.UserID,
		Items:  items,
	}, nil
}

func (s *Service) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	log.Printf("Received OrderPay request for order %d", req.OrderID)

	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		log.Printf("Failed to get order %d: %v", req.OrderID, err)
		return nil, err
	}

	if order == nil || order.Status != models.OrderStatusAwaitingPayment {
		log.Printf("Cannot pay order %d: not found or wrong status (%s)",
			req.OrderID, order.Status)
		return &loms.OrderPayResponse{}, nil
	}

	// Release reserved items
	for _, item := range order.Items {
		log.Printf("Releasing %d items for SKU %d", item.Count, item.SKU)
		if err := s.stockRepo.Release(ctx, item.SKU, uint64(item.Count)); err != nil {
			log.Printf("Failed to release items for SKU %d: %v", item.SKU, err)
			return nil, err
		}
	}

	order.Status = models.OrderStatusPayed
	order.UpdatedAt = time.Now()
	if err := s.orderRepo.Update(ctx, order); err != nil {
		log.Printf("Failed to update order %d: %v", req.OrderID, err)
		return nil, err
	}

	log.Printf("Successfully paid order %d", req.OrderID)
	return &loms.OrderPayResponse{}, nil
}

func (s *Service) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	log.Printf("Received OrderCancel request for order %d", req.OrderID)

	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		log.Printf("Failed to get order %d: %v", req.OrderID, err)
		return nil, err
	}

	if order == nil || order.Status != models.OrderStatusAwaitingPayment {
		log.Printf("Cannot cancel order %d: not found or wrong status (%s)",
			req.OrderID, order.Status)
		return &loms.OrderCancelResponse{}, nil
	}

	// Release reserved items
	for _, item := range order.Items {
		log.Printf("Releasing %d items for SKU %d", item.Count, item.SKU)
		if err := s.stockRepo.Release(ctx, item.SKU, uint64(item.Count)); err != nil {
			log.Printf("Failed to release items for SKU %d: %v", item.SKU, err)
			return nil, err
		}
	}

	order.Status = models.OrderStatusCancelled
	order.UpdatedAt = time.Now()
	if err := s.orderRepo.Update(ctx, order); err != nil {
		log.Printf("Failed to update order %d: %v", req.OrderID, err)
		return nil, err
	}

	log.Printf("Successfully cancelled order %d", req.OrderID)
	return &loms.OrderCancelResponse{}, nil
}

func (s *Service) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	log.Printf("Received StocksInfo request for SKU %d", req.Sku)

	stock, err := s.stockRepo.Get(ctx, req.Sku)
	if err != nil {
		log.Printf("Failed to get stock for SKU %d: %v", req.Sku, err)
		return nil, err
	}

	if stock == nil {
		log.Printf("No stock found for SKU %d", req.Sku)
		return &loms.StocksInfoResponse{Count: 0}, nil
	}

	available := stock.TotalCount - stock.Reserved
	log.Printf("Stock info for SKU %d: total=%d, reserved=%d, available=%d",
		req.Sku, stock.TotalCount, stock.Reserved, available)
	return &loms.StocksInfoResponse{
		Count: available,
	}, nil
}

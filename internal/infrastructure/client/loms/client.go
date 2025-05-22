package loms

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/vagonaizer/loms/api/protos/gen/loms"
)

type Client struct {
	client pb.LOMSClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	log.Printf("Connecting to LOMS service at %s", address)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewLOMSClient(conn)
	log.Printf("Successfully connected to LOMS service at %s", address)
	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	log.Println("Closing LOMS client connection")
	return c.conn.Close()
}

func (c *Client) CreateOrder(ctx context.Context, userID int64, items []*pb.Item) (int64, error) {
	log.Printf("Creating order for user %d with items: %v", userID, items)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.OrderCreate(ctx, &pb.OrderCreateRequest{
		User:  userID,
		Items: items,
	})
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return 0, err
	}

	log.Printf("Successfully created order with ID: %d", resp.OrderID)
	return resp.OrderID, nil
}

func (c *Client) GetStock(ctx context.Context, sku uint32) (uint64, error) {
	log.Printf("Getting stock info for SKU: %d", sku)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.StocksInfo(ctx, &pb.StocksInfoRequest{Sku: sku})
	if err != nil {
		log.Printf("Failed to get stock info: %v", err)
		return 0, err
	}

	log.Printf("Successfully got stock info for SKU %d: %d items available", sku, resp.Count)
	return resp.Count, nil
}

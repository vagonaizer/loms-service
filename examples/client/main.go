package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/vagonaizer/loms/api/protos/gen/loms"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up connection to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLOMSClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test StocksInfo
	log.Println("Testing StocksInfo...")
	stockResp, err := client.StocksInfo(ctx, &pb.StocksInfoRequest{Sku: 773297411})
	if err != nil {
		log.Printf("StocksInfo failed: %v", err)
	} else {
		log.Printf("Available stock: %d", stockResp.Count)
	}

	// Test OrderCreate
	log.Println("\nTesting OrderCreate...")
	orderResp, err := client.OrderCreate(ctx, &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 773297411, Count: 2},
			{Sku: 1002, Count: 1},
		},
	})
	if err != nil {
		log.Printf("OrderCreate failed: %v", err)
	} else {
		log.Printf("Created order ID: %d", orderResp.OrderID)

		// Test OrderInfo
		log.Println("\nTesting OrderInfo...")
		infoResp, err := client.OrderInfo(ctx, &pb.OrderInfoRequest{OrderID: orderResp.OrderID})
		if err != nil {
			log.Printf("OrderInfo failed: %v", err)
		} else {
			log.Printf("Order status: %s", infoResp.Status)
			log.Printf("Order user: %d", infoResp.User)
			log.Printf("Order items: %v", infoResp.Items)
		}

		// Test OrderPay
		log.Println("\nTesting OrderPay...")
		_, err = client.OrderPay(ctx, &pb.OrderPayRequest{OrderID: orderResp.OrderID})
		if err != nil {
			log.Printf("OrderPay failed: %v", err)
		} else {
			log.Println("Order paid successfully")
		}
	}
}

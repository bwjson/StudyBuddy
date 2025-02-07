package grpc

import (
	"context"
	"fmt"
	paypalv1 "github.com/bwjson/Paypal_Proto/gen/go/paypal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api paypalv1.PaypalClient
}

// NewClient создает новый gRPC клиент и подключается к серверу.
func NewClient() (*Client, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:44044", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{api: paypalv1.NewPaypalClient(conn)}, nil
}

// BuySubscription отправляет запрос на покупку подписки.
func (c *Client) BuySubscription(ctx context.Context, email, cardNumber, validUntil, cvv string) (*paypalv1.Response, error) {
	req := &paypalv1.PaymentInfo{
		Email:      email,
		CardNumber: cardNumber,
		ValidUntil: validUntil,
		Cvv:        cvv,
	}

	resp, err := c.api.BuySubscription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling BuySubscription: %w", err)
	}

	return resp, nil
}

// CancelSubscription отменяет подписку.
func (c *Client) CancelSubscription(ctx context.Context, email string) (*paypalv1.Response, error) {
	req := &paypalv1.SubscriptionInfo{
		Email: email,
	}

	resp, err := c.api.CancelSubscription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling CancelSubscription: %w", err)
	}

	return resp, nil
}

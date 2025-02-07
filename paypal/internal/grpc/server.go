package grpc

import (
	"context"
	"github.com/bwjson/Paypal_Microservice/storage/sqlite"
	paypalv1 "github.com/bwjson/Paypal_Proto/gen/go/paypal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type serverAPI struct {
	paypalv1.UnimplementedPaypalServer
	storage *sqlite.Storage
}

func Register(gRPC *grpc.Server, storage *sqlite.Storage) {
	paypalv1.RegisterPaypalServer(gRPC, &serverAPI{storage: storage})
}

func (s *serverAPI) BuySubscription(ctx context.Context, req *paypalv1.PaymentInfo) (*paypalv1.Response, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	validUntilDate, err := time.Parse("01/2006", req.ValidUntil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can not parse the valid until date")
	}

	if validUntilDate.Year() < time.Now().Year() || validUntilDate.Year() == time.Now().Year() && validUntilDate.Month() < time.Now().Month() {
		return nil, status.Error(codes.InvalidArgument, "valid until required")
	}

	if req.GetCvv() == "" || len(req.GetCvv()) != 3 {
		return nil, status.Error(codes.InvalidArgument, "cvv required")
	}

	if req.GetCardNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "card number required")
	}

	cardNumber := req.GetCardNumber()

	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	// Проверка формата (только цифры, 13-19 символов)
	if !regexp.MustCompile(`^\d{13,19}$`).MatchString(cardNumber) {
		return nil, status.Error(codes.InvalidArgument, "invalid card number")
	}

	// Проверяем алгоритм Луна
	if !luhnCheck(cardNumber) {
		return nil, status.Error(codes.InvalidArgument, "invalid card number")
	}

	// Добавление инфы о подписке в бд + проверка

	id, err := s.storage.AddSubscription(ctx, req.Email, req.CardNumber)
	if err != nil {
		return nil, err
	}

	return &paypalv1.Response{
		Response: "You've successfully got the subscription",
		Detail:   "ID:" + strconv.Itoa(int(id)),
	}, nil
}

func luhnCheck(cardNumber string) bool {
	sum := 0
	alt := false
	length := len(cardNumber)

	for i := length - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(cardNumber[i]))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}

	return sum%10 == 0
}

func (s *serverAPI) CancelSubscription(ctx context.Context, req *paypalv1.SubscriptionInfo) (*paypalv1.Response, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	// Checking if the subscription exists and deleting it from DB
	err := s.storage.DeleteSubscription(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &paypalv1.Response{
		Response: "Subscription successfully canceled",
		Detail:   "Subscription Email: " + req.Email,
	}, nil
}

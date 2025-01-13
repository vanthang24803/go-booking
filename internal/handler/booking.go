package handler

import (
	"context"

	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/storage"
	"github.com/may20xx/booking/pkg/log"
)

type BookingService interface {
}

type bookingService struct {
	bookingRepo storage.BookingRepository
	paymentRepo storage.PaymentRepository
}

func NewBookingService() BookingService {
	ctx := context.Background()

	db, err := database.GetDatabase(ctx)
	if err != nil {
		log.Msg.Panic("error getting database connection %s", err)
	}

	return &bookingService{
		bookingRepo: storage.NewBookingRepository(db, ctx),
		paymentRepo: storage.NewPaymentRepository(db, ctx),
	}
}

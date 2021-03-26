package service

import (
	"github.com/joocosta/bloctrial/booking"
	"github.com/joocosta/bloctrial/model"
)

type BookingService struct {
	bookingRepo booking.BookingRepository
}

func NewBookingService(BookingRepos booking.BookingRepository) booking.BookingService {
	return &BookingService{bookingRepo: BookingRepos}
}

func (bs *BookingService) GetSingleBooking(id uint) (*model.Booking, []error) {
	book, errs := bs.bookingRepo.GetSingleBooking(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return book, errs
}

// bookings returns all stored comments
func (bk *BookingService) Bookings(uid uint) ([]model.Booking, []error) {
	bkk, errs := bk.bookingRepo.Bookings(uid)
	if len(errs) > 0 {
		return nil, errs
	}
	return bkk, errs
}

func (bk *BookingService) StoreBooking(booking *model.Booking) (*model.Booking, []error) {
	boks, errs := bk.bookingRepo.StoreBooking(booking)
	if len(errs) > 0 {
		return nil, errs
	}
	return boks, errs
}

func (bk *BookingService) UpdateBooking(Booking *model.Booking) (*model.Booking, []error) {
	book, errs := bk.bookingRepo.UpdateBooking(Booking)
	if len(errs) > 0{
		return nil, errs
	}
	return book, nil
}

func (bk *BookingService) DeleteBooking(id uint) (*model.Booking, []error) {
	books, errs := bk.bookingRepo.DeleteBooking(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return books, errs
}

func (bk *BookingService) BookingExists(userid uint, sched uint) bool {
	exists := bk.bookingRepo.BookingExists(userid, sched)
	return exists
}

package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/booking"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"log"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	bookingService booking.BookingService
}

func NewBookingHandler(BkkService booking.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: BkkService}

}

func (bh *BookingHandler) GetSingleUserBookings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	bookings, errs := bh.bookingService.Bookings(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Can't fetch user books"))
		return
	}

	responses.JSON(w, http.StatusOK, bookings)
	return
}

func (bh *BookingHandler) GetSingleBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	booking, errs := bh.bookingService.GetSingleBooking(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant retrieve movie"))
		return
	}

	responses.JSON(w, http.StatusOK, booking)
	return
}

//adds a hall
func (bh *BookingHandler) PostBooking(w http.ResponseWriter, r *http.Request) {
	var b model.Booking
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant read from request body"))
	}

	log.Println(bh.bookingService.BookingExists(b.UserID, b.ScheduleID))
	if bh.bookingService.BookingExists(b.UserID, b.ScheduleID){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Already booked"))
		return
	}

	book, errs := bh.bookingService.StoreBooking(&b)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant store booking"))
		return
	}

	log.Println(book.UserID)

	responses.JSON(w, http.StatusCreated, book)
	return
}

func (bh *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	var b model.Booking
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, err := strconv.Atoi(id)
	bok, errs := bh.bookingService.GetSingleBooking(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch booking"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	bok, errs = bh.bookingService.UpdateBooking(&b)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update book"))
		return
	}
	responses.JSON(w, http.StatusOK, &bok)
	return
}

func (bh *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	_, errs := bh.bookingService.DeleteBooking(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant delete booking"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}



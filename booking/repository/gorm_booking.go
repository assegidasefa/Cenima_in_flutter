package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/joocosta/bloctrial/booking"
	"github.com/joocosta/bloctrial/model"
	schedrepo "github.com/joocosta/bloctrial/schedule/repository"
	usrRepo "github.com/joocosta/bloctrial/user/repository"
	movRepo "github.com/joocosta/bloctrial/movie/repository"
	"log"
)

// CommentGormRepo implements menu.CommentRepository interface
type BookingGormRepo struct {
	conn *gorm.DB
}

// NewHALLGormRepo returns new object of CommentGormRepo
func NewBookingGormRepo(db *gorm.DB) booking.BookingRepository {
	return &BookingGormRepo{conn: db}
}
func (bkkRepo *BookingGormRepo) Bookings(uid uint) ([]model.Booking, []error) {
	bkk := []model.Booking{}
	errs := bkkRepo.conn.Where("user_id = ?", uid).Find(&bkk).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return bkk, errs
}

func (bkkRepo *BookingGormRepo) GetSingleBooking(id uint) (*model.Booking, []error) {
	book := model.Booking{}
	errs := bkkRepo.conn.First(&book, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &book, errs
}

// StoreComment stores a given customer comment in the database
func (bkkRepo *BookingGormRepo) StoreBooking(booking *model.Booking) (*model.Booking, []error) {
	bkk := booking
	errs := bkkRepo.conn.Create(bkk).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return bkk, errs
}

func (bkkRepo *BookingGormRepo) UpdateBooking(booking *model.Booking) (*model.Booking, []error) {
	bok := booking
	errs := bkkRepo.conn.Save(booking).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return bok, errs
}

func (bkkRepo *BookingGormRepo) DeleteBooking(id uint) (*model.Booking, []error) {
	bok, errs := bkkRepo.GetSingleBooking(id)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = bkkRepo.conn.Delete(bok, id).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return bok, errs
}

func (bkkRepo *BookingGormRepo) BookingExists(userid uint, sched uint) bool {
	MovieRepo := movRepo.NewMovieGormRepo(bkkRepo.conn)
	UserRepo := usrRepo.NewUserGormRepo(bkkRepo.conn)
	usr, err := UserRepo.User(userid)
	if len(err) > 0 {
		return false
	}

	SchedRepo := schedrepo.NewScheduleGormRepo(bkkRepo.conn)
	sch, err := SchedRepo.Schedule(sched)
	if len(err) > 0 {
		return false
	}

	bookings := []model.Booking{}
	movie, _ := MovieRepo.Movie(sch.MovieID)
	log.Println(movie.Title)
	errs := bkkRepo.conn.Where("user_id=? And schedule_id=?", usr.ID, sch.ID).Find(&bookings).GetErrors()

	if len(errs) > 0 {
		return false
	}
	return true
}
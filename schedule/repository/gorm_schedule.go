package repository

import (
	"fmt"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/schedule"
	"log"

	"github.com/jinzhu/gorm"
	schedrepo "github.com/joocosta/bloctrial/cinema/repository"
	movrepo "github.com/joocosta/bloctrial/movie/repository"
)

type ScheduleGormRepo struct {
	conn *gorm.DB
}

func NewScheduleGormRepo(db *gorm.DB) schedule.ScheduleRepository {
	return &ScheduleGormRepo{conn: db}
}

func (scheduleRepo *ScheduleGormRepo) Schedules() ([]model.Schedule, []error) {
	schdls := []model.Schedule{}
	errs := scheduleRepo.conn.Find(&schdls).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs

}
func (scheduleRepo *ScheduleGormRepo) CinemaSchedules(id uint, day string) ([]model.Schedule, []error) {
	CineRepo := schedrepo.NewCinemaGormRepo(scheduleRepo.conn)
	cin, err := CineRepo.Cinema(id)
	if len(err) > 0 {
		return nil, err
	}

	schdls := []model.Schedule{}
	fmt.Printf(day)
	errs := scheduleRepo.conn.Where("cinema_id=? And Day=?", cin.ID, day).Find(&schdls).GetErrors()
	//errs := scheduleRepo.conn.Joins("JOIN schedules on hall_id=halls.id AND day = ?", day).Joins("Join halls on halls.id=cinemas.id").Where("cinemas.id=?", id).Find(&schdls).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}

func (scheduleRepo *ScheduleGormRepo) CinemaSchedulesbyCinema(id uint) ([]model.Schedule, []error) {
	CineRepo := schedrepo.NewCinemaGormRepo(scheduleRepo.conn)
	cin, err := CineRepo.Cinema(id)
	if len(err) > 0 {
		return nil, err
	}

	schdls := []model.Schedule{}
	errs := scheduleRepo.conn.Where("cinema_id=?", cin.ID).Find(&schdls).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, nil

}

func (scheduleRepo *ScheduleGormRepo) StoreSchedule(schedule *model.Schedule) (*model.Schedule, []error) {
	schdl := schedule
	errs := scheduleRepo.conn.Create(schdl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdl, errs
}
func (schRepo *ScheduleGormRepo) UpdateSchedules(schedule *model.Schedule) (*model.Schedule, []error) {
	schdl := schedule
	errs := schRepo.conn.Save(schdl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdl, errs
}
func (schRepo *ScheduleGormRepo) UpdateSchedulesBooked(schedule *model.Schedule, Amount uint) *model.Schedule {
	schdl := schedule
	schRepo.conn.Model(&schdl).UpdateColumn("booked", schdl.Booked+1)

	return schdl
}

// DeleteComment deletes a given customer comment from the database
func (schRepo *ScheduleGormRepo) DeleteSchedules(id uint) (*model.Schedule, []error) {
	fmt.Println("(((((((((((((((((((((((((in delete gorm))))))))))))))))))))))")

	schdl, errs := schRepo.Schedule(id)
	fmt.Println("(((((((((((((((((((((((((in delete gorm))))))))))))))))))))))")
	if len(errs) > 0 {
		return nil, errs
	}

	errs = schRepo.conn.Delete(schdl, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdl, errs
}
func (schRepo *ScheduleGormRepo) Schedule(id uint) (*model.Schedule, []error) {
	schdl := model.Schedule{}
	errs := schRepo.conn.First(&schdl, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &schdl, errs
}

func (scheduleRepo *ScheduleGormRepo) ScheduleExists(cinemaid uint, movieid uint) bool {
	CineRepo := schedrepo.NewCinemaGormRepo(scheduleRepo.conn)
	cin, err := CineRepo.Cinema(cinemaid)
	if len(err) > 0 {
		return false
	}

	MovieRepo := movrepo.NewMovieGormRepo(scheduleRepo.conn)
	mov, err := MovieRepo.Movie(movieid)
	if len(err) > 0 {
		return false
	}

	schdls := []model.Schedule{}
	log.Println(cin.Name)
	log.Println(mov.Title)
	errs := scheduleRepo.conn.Where("cinema_id=? And movie_id=?", cin.ID, mov.Id).Find(&schdls).GetErrors()

	if len(errs) > 0 {
		return false
	}
	return true
}
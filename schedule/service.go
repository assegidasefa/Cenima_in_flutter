package schedule

import "github.com/joocosta/bloctrial/model"

type ScheduleService interface {
	Schedules() ([]model.Schedule, []error)
	StoreSchedule(schedule *model.Schedule) (*model.Schedule, []error)
	CinemaSchedules(id uint, day string) ([]model.Schedule, []error)
	CinemaSchedulesbyCinema(id uint) ([]model.Schedule, []error)
	Schedule(id uint) (*model.Schedule, []error)
	UpdateSchedules(hall *model.Schedule) (*model.Schedule, []error)
	UpdateSchedulesBooked(user *model.Schedule, Amount uint) *model.Schedule
	DeleteSchedules(id uint) (*model.Schedule, []error)
	ScheduleExists(cinemaid uint, movieid uint) bool
}
package service

import (
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/schedule"
)

type ScheduleService struct {
	scheduleRepo schedule.ScheduleRepository
}

func NewScheduleService(schRepo schedule.ScheduleRepository) schedule.ScheduleService {
	return &ScheduleService{scheduleRepo: schRepo}
}

func (s *ScheduleService) Schedules() ([]model.Schedule, []error) {
	schdls, errs := s.scheduleRepo.Schedules()
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}

func (s *ScheduleService) CinemaSchedules(id uint, day string) ([]model.Schedule, []error) {
	schdls, errs := s.scheduleRepo.CinemaSchedules(id, day)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}
func (s *ScheduleService) StoreSchedule(schedule *model.Schedule) (*model.Schedule, []error) {
	schdls, errs := s.scheduleRepo.StoreSchedule(schedule)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}
func (ss *ScheduleService) UpdateSchedules(schedule *model.Schedule) (*model.Schedule, []error) {
	schdls, errs := ss.scheduleRepo.UpdateSchedules(schedule)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}
func (ss *ScheduleService) UpdateSchedulesBooked(schedule *model.Schedule, Amount uint) *model.Schedule {
	schdls := ss.scheduleRepo.UpdateSchedulesBooked(schedule, Amount)

	return schdls
}
func (ss *ScheduleService) DeleteSchedules(id uint) (*model.Schedule, []error) {
	schdls, errs := ss.scheduleRepo.DeleteSchedules(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}

func (ss *ScheduleService) Schedule(id uint) (*model.Schedule, []error) {
	schdls, errs := ss.scheduleRepo.Schedule(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}


func (s *ScheduleService) CinemaSchedulesbyCinema(id uint) ([]model.Schedule, []error) {
	schdls, errs := s.scheduleRepo.CinemaSchedulesbyCinema(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return schdls, errs
}

func (s *ScheduleService) ScheduleExists(cinemaid uint, movieid uint) bool {
	exists := s.scheduleRepo.ScheduleExists(cinemaid, movieid)
	return exists
}


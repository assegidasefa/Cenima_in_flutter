package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/schedule"
	"log"
	"strconv"
	"strings"

	"net/http"
)

type ScheduleHandler struct {
	scheduleService schedule.ScheduleService
}

func NewScheduleHandler(schdlService schedule.ScheduleService) *ScheduleHandler {
	fmt.Println("admin schedule handler created")
	return &ScheduleHandler{scheduleService: schdlService}
}

func (as *ScheduleHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, errs := as.scheduleService.Schedules()
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Can't fetch schedules"))
		return
	}

	responses.JSON(w, http.StatusOK, schedules)
	return
}

func (as *ScheduleHandler) GetSchedulesCinemaDay(w http.ResponseWriter, r *http.Request) {
	var day string
	params := mux.Vars(r)
	cinemaid, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not passed"))
		return
	}
	id, _ := strconv.Atoi(cinemaid)

	if len(r.URL.Query()) > 0 {
		if scheduleCinema := r.URL.Query().Get("day"); scheduleCinema != "" {
			scheduleCinema = strings.ToLower(scheduleCinema)
			day = scheduleCinema
		}
	}

	log.Println(r.URL.Query())
	cinemaSchedules, errs := as.scheduleService.CinemaSchedules(uint(id), day)
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
	}
	responses.JSON(w, http.StatusOK, cinemaSchedules)
	return
}

func (as *ScheduleHandler) PostSchedule(w http.ResponseWriter, r *http.Request) {
	var sh model.Schedule
	err := json.NewDecoder(r.Body).Decode(&sh)
	if err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant read from request body"))
	}

	log.Println(as.scheduleService.ScheduleExists(sh.CinemaID, sh.MovieID))
	if as.scheduleService.ScheduleExists(sh.CinemaID, sh.MovieID){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("schedule already exists"))
		return
	}

	schedule, errs := as.scheduleService.StoreSchedule(&sh)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant store schedule"))
		return
	}

	log.Println(schedule.MovieID)

	responses.JSON(w, http.StatusCreated, schedule)
	return
}
func (as *ScheduleHandler) GetSingleSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	schedule, errs := as.scheduleService.Schedule(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant retrieve schedule"))
		return
	}

	responses.JSON(w, http.StatusOK, schedule)
	return
}

func (as *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	_, errs := as.scheduleService.DeleteSchedules(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant delete schedule"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (as *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	var sh model.Schedule
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}

	idd, err := strconv.Atoi(id)
	sch, errs := as.scheduleService.Schedule(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch schedule"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&sh)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	sch, errs = as.scheduleService.UpdateSchedules(&sh)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update schedule"))
		return
	}
	responses.JSON(w, http.StatusOK, &sch)
	return
}
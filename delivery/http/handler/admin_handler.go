package handler

import (
	"encoding/json"
	"errors"
	"log"

	//"fmt"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/cinema"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/movie"
	"github.com/joocosta/bloctrial/schedule"
	//"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	csrv        cinema.CinemaService
	ssrv        schedule.ScheduleService
	msrv        movie.MovieService
	csrfSignKey []byte
}

func NewAdminHandler(cs cinema.CinemaService, ss schedule.ScheduleService, ms movie.MovieService) *AdminHandler {

	return &AdminHandler{csrv: cs, ssrv: ss, msrv: ms}

}

//CINEMA
func (ah *AdminHandler) AdminCinema(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	cinema, errs := ah.csrv.Cinema(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Cinema"))
	}
	responses.JSON(w, http.StatusOK, cinema)
}

func (ah *AdminHandler) AdminCinemas(w http.ResponseWriter, r *http.Request) {
	cinemas, errs := ah.csrv.Cinemas()
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Cinemas"))
	}
	responses.JSON(w, http.StatusOK, cinemas)
}

func (ah *AdminHandler) AdminDeleteCinema(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	cinema, errs := ah.csrv.DeleteCinema(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Delete Cinema"))
	}
	responses.JSON(w, http.StatusNoContent, cinema)

}
func (ah *AdminHandler) AdminCinemaUpdateList(w http.ResponseWriter, r *http.Request){
	var c model.Cinema
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}

	idd, err := strconv.Atoi(id)
	cin, errs := ah.csrv.Cinema(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch cinema"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	cin, errs = ah.csrv.UpdateCinema(&c)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update cinema"))
		return
	}
	responses.JSON(w, http.StatusOK, cin)
	return
}

func (ah *AdminHandler) AdminCinemaNew(w http.ResponseWriter, r *http.Request){
	var c model.Cinema
	err := json.NewDecoder(r.Body).Decode(&c)
	if err!=nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cinema decoding failed"))
		return
	}
	if ah.csrv.CinemaExists(c.Name){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cinema already exists"))
		return
		//json.NewEncoder(w).Encode(err)
	}

	newCinema, errs := ah.csrv.StoreCinema(&c)
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Adding new Cinema Failed"))
		return
	}

	responses.JSON(w, http.StatusOK, newCinema)
	//json.NewEncoder(w).Encode(newCinema)
}



//SCHEDULES
func (ah *AdminHandler) AdminSchedules(w http.ResponseWriter, r *http.Request){
	var day string
	params := mux.Vars(r)
	cinemaid, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not passed"))
		return
	}
	id, _ := strconv.Atoi(cinemaid)

	if len(r.URL.Query()) == 0{
		log.Println(r.URL.Query())
		cinemaSchedules, errs := ah.ssrv.CinemaSchedulesbyCinema(uint(id))
		if len(errs) > 0{
			responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
		}
		responses.JSON(w, http.StatusOK, cinemaSchedules)
	}else{
		if scheduleCinema := r.URL.Query().Get("day"); scheduleCinema != "" {
			day = scheduleCinema
		}
		log.Println(r.URL.Query())
		cinemaSchedules, errs := ah.ssrv.CinemaSchedules(uint(id), day)
		if len(errs) > 0{
			responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
		}
		responses.JSON(w, http.StatusOK, cinemaSchedules)
	}
}

//func (ah *AdminHandler) AdminSchedules(w http.ResponseWriter, r *http.Request){
//	schedules, errs := ah.ssrv.Schedules()
//	if len(errs) > 0{
//		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Schedules"))
//	}
//	responses.JSON(w, http.StatusOK, schedules)
//}
//
//func (ah *AdminHandler) AdminSchedule(w http.ResponseWriter, r *http.Request){
//	params := mux.Vars(r)
//	id, exists := params["id"]
//	if !exists {
//		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
//		return
//	}
//	idd, _ := strconv.Atoi(id)
//
//	sch, errs := ah.ssrv.Schedule(uint(idd))
//	if len(errs) > 0{
//		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Schedule"))
//	}
//	responses.JSON(w, http.StatusOK, sch)
//}

func (ah *AdminHandler) AdminScheduleNew(w http.ResponseWriter, r *http.Request){
	var sh model.Schedule
	err := json.NewDecoder(r.Body).Decode(&sh)
	if err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant read from request body"))
	}

	log.Println(ah.ssrv.ScheduleExists(sh.CinemaID, sh.MovieID))
	if ah.ssrv.ScheduleExists(sh.CinemaID, sh.MovieID){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("schedule already exists"))
		return
	}

	schedule, errs := ah.ssrv.StoreSchedule(&sh)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant store schedule"))
		return
	}

	log.Println(schedule.MovieID)

	responses.JSON(w, http.StatusCreated, schedule)
	return
}

func (ah *AdminHandler) AdminScheduleUpdate(w http.ResponseWriter, r *http.Request){
	var sh model.Schedule
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}

	idd, err := strconv.Atoi(id)
	sch, errs := ah.ssrv.Schedule(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch schedule"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&sh)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	sch, errs = ah.ssrv.UpdateSchedules(&sh)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update schedule"))
		return
	}
	responses.JSON(w, http.StatusOK, &sch)
	return

}

func (ah *AdminHandler) AdminScheduleDelete(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	sch, errs := ah.ssrv.DeleteSchedules(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Delete schedule"))
	}
	responses.JSON(w, http.StatusNoContent, sch)
}


//MOVIES
func (ah *AdminHandler) AdminMovies(w http.ResponseWriter, r *http.Request) {

	movies, errs := ah.msrv.Movies()
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Movies"))
	}
	responses.JSON(w, http.StatusOK, movies)
}

func (ah *AdminHandler) AdminMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	mov, errs := ah.msrv.Movie(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Movie"))
	}
	responses.JSON(w, http.StatusOK, mov)
}

func (ah *AdminHandler) AdminMovieNew(w http.ResponseWriter, r *http.Request){
	var m model.Movie
	err := json.NewDecoder(r.Body).Decode(&m)
	if err!=nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Movie decoding failed"))
		return
	}
	if ah.msrv.MovieExists(m.Title){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Movie already exists"))
		return
		//json.NewEncoder(w).Encode(err)
	}

	newMovie, errs := ah.msrv.StoreMovie(&m)
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Adding new Movie Failed"))
		return
	}

	responses.JSON(w, http.StatusOK, newMovie)
	//json.NewEncoder(w).Encode(newCinema)
}

func (ah *AdminHandler) AdminDeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	mov, errs := ah.msrv.DeleteMovie(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Delete Movie"))
	}
	responses.JSON(w, http.StatusNoContent, mov)
}

func (ah *AdminHandler) AdminMovieUpdateList(w http.ResponseWriter, r *http.Request){
	var u model.Movie
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, err := strconv.Atoi(id)
	mov, errs := ah.msrv.Movie(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch movie"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	mov, errs = ah.msrv.UpdateMovie(&u)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update movie"))
		return
	}
	responses.JSON(w, http.StatusOK, &mov)
	return
}
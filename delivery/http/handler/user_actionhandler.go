package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/booking"
	"github.com/joocosta/bloctrial/cinema"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/movie"
	"github.com/joocosta/bloctrial/rtoken"
	"github.com/joocosta/bloctrial/schedule"
	"github.com/joocosta/bloctrial/user"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserActionHandler struct {
	csrv   cinema.CinemaService
	ssrv   schedule.ScheduleService
	msrv   movie.MovieService
	usrv   user.UserService
	tsrv   rtoken.Service
	bsrv   booking.BookingService
}

func NewUserActionHandler(cs cinema.CinemaService, ss schedule.ScheduleService, ms movie.MovieService, u user.UserService, t rtoken.Service, b booking.BookingService) *UserActionHandler {
	return &UserActionHandler{csrv: cs, ssrv: ss, msrv: ms, tsrv:t, usrv: u, bsrv: b}
}
func (uah *UserActionHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var u model.User
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, err := strconv.Atoi(id)
	//log.Println(idd)
	usr, errs := uah.usrv.User(uint(idd))

	curid := usr.ID
	pass := usr.Password
	amo := usr.Amount
	rol := usr.RoleID

	log.Println(curid)
	log.Println(pass)
	log.Println(amo)
	log.Println(rol)
	log.Println(usr.FullName)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch user"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	passnew, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Password Encryption  failed"))
		return
	}

	u.Password = string(passnew)
	u.RoleID = rol
	u.Amount = amo

	usr, errs = uah.usrv.UpdateUser(&u)

	log.Println(usr.ID)
	log.Println(usr.Password)
	log.Println(usr.Amount)
	log.Println(usr.RoleID)
	log.Println(usr.FullName)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update user"))
		return
	}
	responses.JSON(w, http.StatusOK, &usr)
	return
}

func (uah *UserActionHandler) UserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not given"))
		return
	}
	idd, _ := strconv.Atoi(id)
	_, errs := uah.usrv.DeleteUser(uint(idd))
	if len(errs) != 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Couldn't delete user"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (uah *UserActionHandler) Movies(w http.ResponseWriter, r *http.Request){
	//var mov []model.Movie
	nowshowingmovies, errs := uah.msrv.Movies()
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Movies"))
	}
	responses.JSON(w, http.StatusOK, nowshowingmovies)
}

func (uah *UserActionHandler) SingleMovie(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, err := strconv.Atoi(id)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Can't get movie ID"))
	}

	curMovie, errs := uah.msrv.Movie(uint(idd))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Movies"))
	}
	responses.JSON(w, http.StatusOK, curMovie)
}

func (uah *UserActionHandler) Search(w http.ResponseWriter, r *http.Request){
	var query string

	if len(r.URL.Query()) > 0 {
		if keyword := r.URL.Query().Get("query"); keyword != "" {
			query = keyword
		}
	}

	results, err := uah.msrv.SearchMovie(query)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	responses.JSON(w, http.StatusOK, results)
}


func (uah *UserActionHandler) Cinemas(w http.ResponseWriter, r *http.Request){
	cinemas, errs := uah.csrv.Cinemas()
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant Fetch Cinemas"))
	}
	responses.JSON(w, http.StatusOK, cinemas)
}
func (uah *UserActionHandler) GetSingleCinema(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	cinema, errs := uah.csrv.Cinema(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant retrieve cinema"))
		return
	}

	responses.JSON(w, http.StatusOK, cinema)
	return
}

func (uah *UserActionHandler) CinemaSchedule(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	cinemaid, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not passed"))
		return
	}
	id, _ := strconv.Atoi(cinemaid)

	//if len(r.URL.Query()) == 0{
	log.Println(r.URL.Query())
	cinemaSchedules, errs := uah.ssrv.CinemaSchedulesbyCinema(uint(id))
	if len(errs) > 0{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
	}
	responses.JSON(w, http.StatusOK, cinemaSchedules)
	//}
	//var day string
	//params := mux.Vars(r)
	//cinemaid, exists := params["id"]
	//if !exists {
	//	responses.ERROR(w, http.StatusBadRequest, errors.New("id not passed"))
	//	return
	//}
	//id, _ := strconv.Atoi(cinemaid)
	//
	//if len(r.URL.Query()) == 0{
	//	log.Println(r.URL.Query())
	//	cinemaSchedules, errs := uah.ssrv.CinemaSchedulesbyCinema(uint(id))
	//	if len(errs) > 0{
	//		responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
	//	}
	//	responses.JSON(w, http.StatusOK, cinemaSchedules)
	//}else{
	//	if scheduleCinema := r.URL.Query().Get("day"); scheduleCinema != "" {
	//		day = scheduleCinema
	//	}
	//	log.Println(r.URL.Query())
	//	cinemaSchedules, errs := uah.ssrv.CinemaSchedules(uint(id), day)
	//	if len(errs) > 0{
	//		responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
	//	}
	//	responses.JSON(w, http.StatusOK, cinemaSchedules)
	//}
}

func (uah *UserActionHandler) Bookings(w http.ResponseWriter, r *http.Request){
	_token := r.Header.Get("Authorization")
	_token = strings.Replace(_token, "Bearer ", "", 1)
	fmt.Println(_token)
	claim, err := uah.tsrv.GetClaims(_token)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	var sm []model.Schedule
	user := claim.User
	log.Println(user.FullName)
	bookings, errs := uah.bsrv.Bookings(user.ID)
	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	for _, book := range bookings {
		s, _ := uah.ssrv.Schedule(book.ScheduleID)
		sm = append(sm, *s)
	}
	responses.JSON(w, http.StatusOK, sm)
}


//func (uah *UserActionHandler) CinemaScheduleByDayAndCinema(w http.ResponseWriter, r *http.Request){
//	var day string
//	params := mux.Vars(r)
//	cinemaid, exists := params["id"]
//	if !exists {
//		responses.ERROR(w, http.StatusBadRequest, errors.New("id not passed"))
//		return
//	}
//	id, _ := strconv.Atoi(cinemaid)
//
//	if len(r.URL.Query()) > 0 {
//		if scheduleCinema := r.URL.Query().Get("day"); scheduleCinema != "" {
//			day = scheduleCinema
//		}
//	}
//
//	log.Println(r.URL.Query())
//	cinemaSchedules, errs := uah.ssrv.CinemaSchedules(uint(id), day)
//	if len(errs) > 0{
//		responses.ERROR(w, http.StatusInternalServerError, errors.New("can't get Cinema Schedules"))
//	}
//	responses.JSON(w, http.StatusOK, cinemaSchedules)
//}

//func (uah *UserActionHandler) BookSchedule(w http.ResponseWriter, r *http.Request){
//
//}




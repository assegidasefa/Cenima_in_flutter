package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/movie"
	"log"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService movie.MovieService
}

func NewMovieHander(mvService movie.MovieService) *MovieHandler {
	return &MovieHandler{movieService: mvService}
}

//Get Single Movie
func (mh *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	movie, errs := mh.movieService.Movie(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant retrieve movie"))
		return
	}

	responses.JSON(w, http.StatusOK, movie)
	return
}

//Get all Movies
func (mh *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {

	movies, errs := mh.movieService.Movies()

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Can't fetch movies"))
		return
	}

	responses.JSON(w, http.StatusOK, movies)
	return
}

func (mh *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) {
	var m model.Movie
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant read from request body"))
	}

	log.Println(mh.movieService.MovieExists(m.Title))
	if mh.movieService.MovieExists(m.Title){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Movie already exists"))
		return
	}

	movie, errs := mh.movieService.StoreMovie(&m)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant store movie"))
		return
	}

	log.Println(movie.Title)

	responses.JSON(w, http.StatusCreated, movie)
	return
}

func (mh *MovieHandler) MovieUpdate(w http.ResponseWriter, r *http.Request) {
	var u model.Movie
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, err := strconv.Atoi(id)
	usr, errs := mh.movieService.Movie(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch movie"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	usr, errs = mh.movieService.UpdateMovie(&u)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update movie"))
		return
	}
	responses.JSON(w, http.StatusOK, &usr)
	return
}

func (mh *MovieHandler) MovieDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	_, errs := mh.movieService.DeleteMovie(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant delete movie"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}


func (mh *MovieHandler) Search(w http.ResponseWriter, r *http.Request){
	var query string

	if len(r.URL.Query()) > 0 {
		if keyword := r.URL.Query().Get("query"); keyword != "" {
			query = keyword
		}
	}

	results, err := mh.movieService.SearchMovie(query)
	
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	responses.JSON(w, http.StatusOK, results)
}

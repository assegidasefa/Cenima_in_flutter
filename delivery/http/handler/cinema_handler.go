package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joocosta/bloctrial/cinema"
	"github.com/joocosta/bloctrial/delivery/http/responses"
	"github.com/joocosta/bloctrial/model"
	"log"
	"net/http"
	"strconv"
)

type CinemaHandler struct {
	cinemaService cinema.CinemaService
}

func NewCinemaHandler(CllService cinema.CinemaService) *CinemaHandler {
	return &CinemaHandler{cinemaService: CllService}

}
func (cc *CinemaHandler) GetCinemas(w http.ResponseWriter, r *http.Request) {
	cinemas, errs := cc.cinemaService.Cinemas()

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Can't fetch cinemas"))
		return
	}

	responses.JSON(w, http.StatusOK, cinemas)
	return

}
func (ach *CinemaHandler) PostCinema(w http.ResponseWriter, r *http.Request) {
	var c model.Cinema
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil{
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant read from request body"))
	}

	log.Println(ach.cinemaService.CinemaExists(c.Name))
	if ach.cinemaService.CinemaExists(c.Name){
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cinema already exists"))
		return
	}

	cinema, errs := ach.cinemaService.StoreCinema(&c)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant store movie"))
		return
	}

	log.Println(cinema.Name)

	responses.JSON(w, http.StatusCreated, cinema)
	return
}

// GetSingleCinema
func (ach *CinemaHandler) GetSingleCinema(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	cinema, errs := ach.cinemaService.Cinema(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant retrieve cinema"))
		return
	}

	responses.JSON(w, http.StatusOK, cinema)
	return
}

func (ach *CinemaHandler) CinemaUpdate(w http.ResponseWriter, r *http.Request) {
	var c model.Cinema
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}

	idd, err := strconv.Atoi(id)
	usr, errs := ach.cinemaService.Cinema(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant fetch cinema"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	usr, errs = ach.cinemaService.UpdateCinema(&c)

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant update cinema"))
		return
	}
	responses.JSON(w, http.StatusOK, &usr)
	return
}

func (ach *CinemaHandler) CinemaDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, exists := params["id"]
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("id not provided"))
		return
	}
	idd, _ := strconv.Atoi(id)

	_, errs := ach.cinemaService.DeleteCinema(uint(idd))

	if len(errs) > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Cant delete movie"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

//package handler
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/joocosta/bloctrial/cinema"
//	"github.com/joocosta/bloctrial/model"
//	"log"
//	"net/http"
//	"strconv"
//
//	"github.com/julienschmidt/httprouter"
//)
//
//type CinemaHandler struct {
//	cinemaService cinema.CinemaService
//}
//
//func NewCinemaHandler(CllService cinema.CinemaService) *CinemaHandler {
//	return &CinemaHandler{cinemaService: CllService}
//
//}
//func (cc *CinemaHandler) GetCinemas(w http.ResponseWriter,
//	r *http.Request, _ httprouter.Params) {
//
//	cinemas, errs := cc.cinemaService.Cinemas()
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	output, err := json.MarshalIndent(cinemas, "", "\t\t")
//
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(output)
//	return
//
//}
//func (ach *CinemaHandler) PostCinema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//	l := r.ContentLength
//	body := make([]byte, l)
//	r.Body.Read(body)
//	cinema := &model.Cinema{}
//
//	err := json.Unmarshal(body, &cinema)
//
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	cinema, errs := ach.cinemaService.StoreCinema(cinema)
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	p := fmt.Sprintf("/cinemas/?id=%d", cinema.ID)
//	log.Println(cinema.Name)
//	w.Header().Set("Location", p)
//	w.WriteHeader(http.StatusCreated)
//	return
//}
//
//// GetSingleCinema
//func (ach *CinemaHandler) GetSingleCinema(w http.ResponseWriter,
//	r *http.Request, ps httprouter.Params) {
//
//	id, err := strconv.Atoi(ps.ByName("id"))
//
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//	log.Println(id)
//
//	cinema, errs := ach.cinemaService.Cinema(uint(id))
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	output, err := json.MarshalIndent(cinema, "", "\t\t")
//	log.Println(cinema)
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(output)
//	return
//}
//
//func (ach *CinemaHandler) CinemaUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	id, err := strconv.Atoi(ps.ByName("id"))
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	cinema, errs := ach.cinemaService.Cinema(uint(id))
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	l := r.ContentLength
//
//	body := make([]byte, l)
//	r.Body.Read(body)
//	for _, b := range body{
//		log.Print(string(b))
//	}
//
//	var cin model.Cinema
//
//	err = json.Unmarshal(body, &cin)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	cinema, errs = ach.cinemaService.UpdateCinema(&cin)
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	output, err := json.MarshalIndent(cinema, "", "\t\t")
//
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(output)
//	return
//	//id, err := strconv.Atoi(ps.ByName("id"))
//	//if err != nil {
//	//	w.Header().Set("Content-Type", "application/json")
//	//	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//	//	return
//	//}
//	//
//	//cinema, errs := ach.cinemaService.Cinema(uint(id))
//	//
//	//if len(errs) > 0 {
//	//	w.Header().Set("Content-Type", "application/json")
//	//	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//	//	return
//	//}
//	//
//	//l := r.ContentLength
//	//
//	//body := make([]byte, l)
//	//r.Body.Read(body)
//	//for _, b := range body{
//	//	log.Print(string(b))
//	//}
//	//
//	//var cin model.Cinema
//	//
//	//err = json.Unmarshal(body, &cin)
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	//
//	//cinema, errs = ach.cinemaService.UpdateCinema(&cin)
//	//
//	//if len(errs) > 0 {
//	//	w.Header().Set("Content-Type", "application/json")
//	//	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//	//	return
//	//}
//	//
//	//output, err := json.MarshalIndent(cinema, "", "\t\t")
//	//
//	//if err != nil {
//	//	w.Header().Set("Content-Type", "application/json")
//	//	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//	//	return
//	//}
//	//
//	//w.Header().Set("Content-Type", "application/json")
//	//w.Write(output)
//	//return
//}
//
//func (ach *CinemaHandler) CinemaDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//	id, err := strconv.Atoi(ps.ByName("id"))
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	_, errs := ach.cinemaService.DeleteCinema(uint(id))
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusNoContent)
//	return
//}

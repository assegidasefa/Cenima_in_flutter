package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joocosta/bloctrial/delivery/http/handler"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/rtoken"

	borepo "github.com/joocosta/bloctrial/booking/repository"
	boserv "github.com/joocosta/bloctrial/booking/service"
	crepo "github.com/joocosta/bloctrial/cinema/repository"
	cserv "github.com/joocosta/bloctrial/cinema/service"

	screpo "github.com/joocosta/bloctrial/schedule/repository"
	scserv "github.com/joocosta/bloctrial/schedule/service"

	usrepo "github.com/joocosta/bloctrial/user/repository"
	userv "github.com/joocosta/bloctrial/user/service"

	mrepo "github.com/joocosta/bloctrial/movie/repository"
	mserv "github.com/joocosta/bloctrial/movie/service"
)

func main() {
	db, err := gorm.Open("postgres", "postgres://postgres:9993@localhost/movieevent4?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&model.Cinema{})
	db.AutoMigrate(&model.Movie{})
	db.AutoMigrate(&model.Schedule{})
	db.AutoMigrate(&model.Booking{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Role{})

	db.AutoMigrate(&model.Role{ID: 1, Name: "USER"})
	db.AutoMigrate(&model.Role{ID: 2, Name: "ADMIN"})

	token := rtoken.Service{}

	CinemaRepo := crepo.NewCinemaGormRepo(db)
	Cinemasr := cserv.NewCinemaService(CinemaRepo)

	MovieRepo := mrepo.NewMovieGormRepo(db)
	Moviesr := mserv.NewMovieService(MovieRepo)

	ScheduleRepo := screpo.NewScheduleGormRepo(db)
	Schedulesr := scserv.NewScheduleService(ScheduleRepo)

	BookingRepo := borepo.NewBookingGormRepo(db)
	Bookingsr := boserv.NewBookingService(BookingRepo)

	UserRepo := usrepo.NewUserGormRepo(db)
	usersr := userv.NewUserService(UserRepo)

	roleRepo := usrepo.NewRoleGormRepo(db)
	rolesr := userv.NewRoleService(roleRepo)

	//ch := handler.NewCinemaHandler(Cinemasr)
	mh := handler.NewMovieHander(Moviesr)
	//sch := handler.NewScheduleHandler(Schedulesr)
	//boh := handler.NewBookingHandler(Bookingsr)

	uh := handler.NewUserHandler(usersr, rolesr, token)
	uah := handler.NewUserActionHandler(Cinemasr, Schedulesr, Moviesr, usersr, token, Bookingsr)
	ah := handler.NewAdminHandler(Cinemasr, Schedulesr, Moviesr)

	router := mux.NewRouter()

	router.HandleFunc("/api/search", uh.Authenticated(mh.Search)).Queries("query", "{query}").Methods("GET")

	router.HandleFunc("/movies", ah.AdminMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", ah.AdminMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", ah.AdminMovieUpdateList).Methods("PUT")
	router.HandleFunc("/movies/{id}", ah.AdminDeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies", ah.AdminMovieNew).Methods("POST")
	
	router.HandleFunc("/cinemas", ah.AdminCinemas).Methods("GET")
	router.HandleFunc("/cinemas/{id}", ah.AdminCinema).Methods("GET")
	router.HandleFunc("/cinemas/{id}", ah.AdminCinemaUpdateList).Methods("PUT")
	router.HandleFunc("/cinemas/{id}", ah.AdminDeleteCinema).Methods("DELETE")
	router.HandleFunc("/cinemas", ah.AdminCinemaNew).Methods("POST")
	
	router.HandleFunc("/schedules/{id}", ah.AdminSchedules).Queries("day", "{day}").Methods("GET")
	router.HandleFunc("/schedules/{id}", ah.AdminScheduleUpdate).Methods("PUT")
	router.HandleFunc("/schedules/{id}", ah.AdminScheduleDelete).Methods("DELETE")
	router.HandleFunc("/schedules", ah.AdminScheduleNew).Methods("POST")
	//////
	//router.HandleFunc("/api/bookings", uh.Authenticated(uah.Bookings)).Methods("GET")
	////router.HandleFunc("/bookings/{id}", uh.Authenticated(boh.GetSingleBook)).Methods("GET")
	//router.HandleFunc("/bookings/{id}", uh.Authenticated(uah.UpdateBooking)).Methods("PUT")
	//router.HandleFunc("/bookings/{id}", uh.Authenticated(uah.DeleteBooking)).Methods("DELETE")
	//router.HandleFunc("/bookings", uh.Authenticated(uah.PostBooking)).Methods("POST")
	//////
	router.HandleFunc("/login", uh.Login).Methods("POST")
	router.HandleFunc("/signup", uh.SignUp).Methods("POST")
	router.HandleFunc("/logout", uh.Authenticated(uh.Logout)).Methods("POST")
	router.HandleFunc("/user/{id}", uh.Authenticated(uah.UserUpdate)).Methods("PUT")
	router.HandleFunc("/user/{id}", uh.Authenticated(uah.UserDelete)).Methods("DELETE")

	http.ListenAndServe(":8181", router)

}

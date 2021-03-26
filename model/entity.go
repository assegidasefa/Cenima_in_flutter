package model

//import "time"

type Movie struct {
	Id          int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Image		string `json:"poster_path"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	IMDBRating  float32`json:"rating"`
}

type Cinema struct{
	ID	uint  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
	Price uint `json:"price"`
	Schedules []Schedule `json:"schedules"`
	VIPPrice uint `json:"vipprice"`
	Capacity uint `json:"capacity"`
	VIPCapacity uint `json:"vipcapacity"`
}

type Schedule struct {
	ID           uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MovieID      uint   `json:"movieid" gorm:"not null"`
	StartingTime string `json:"startingtime" gorm:"type:varchar(255);not null"`
	Dimension    string `json:"dimension" gorm:"type:varchar(255);not null"`
	CinemaID     uint   `json:"cinemaid" gorm:"not null;"`
	Day          string `json:"day" gorm:"type:varchar(255);not null"`
	Booked       uint   `json:"booked" gorm:"DEFAULT:0"`
}

type Booking struct {
	ID         uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID     uint `json:"userid"`
	ScheduleID uint `json:"scheduleid" `
}

type Role struct {
	ID   uint
	Name string `gorm:"type:varchar(255)"`
}
type User struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	FullName string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string `json:"pass" gorm:"type:varchar(255)"`
	RoleID   uint
	Amount   uint `json:"amount" gorm:"DEFAULT:300"`
}
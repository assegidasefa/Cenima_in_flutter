package repository

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/movie"
	"log"
)

type MovieGormRepo struct {
	conn *gorm.DB
}

func NewMovieGormRepo(db *gorm.DB) movie.MovieRepository {
	return &MovieGormRepo{conn: db}

}

func (movieRepo *MovieGormRepo) Movies() ([]model.Movie, []error) {
	mvs := []model.Movie{}
	log.Println("Tried getting")
	errs := movieRepo.conn.Find(&mvs).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return mvs, errs

}

func (movieRepo *MovieGormRepo) StoreMovie(movie *model.Movie) (*model.Movie, []error) {
	mv := movie
	errs := movieRepo.conn.Create(mv).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return mv, errs
}


// Event retrieve a event from the database by its id
func (movieRepo *MovieGormRepo) Movie(id uint) (*model.Movie, []error) {
	ctg := model.Movie{}
	errs := movieRepo.conn.First(&ctg, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &ctg, errs
}

// UpdateEvent updates a given event in the database
func (movieRepo *MovieGormRepo) UpdateMovie(movie *model.Movie) (*model.Movie, []error) {
	eve := movie
	errs := movieRepo.conn.Save(eve).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return eve, errs
}

// DeleteEvent deletes a given event from the database
func (movieRepo *MovieGormRepo) DeleteMovie(id uint) (*model.Movie, []error) {
	fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
	eve, errs := movieRepo.Movie(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = movieRepo.conn.Delete(eve, eve.Id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return eve, errs
}
func (movieRepo *MovieGormRepo) SearchMovie(index string) ([]model.Movie, error) {
	items := []model.Movie{}

	err := movieRepo.conn.Where("title ILIKE ?", "%"+index+"%").Find(&items).GetErrors()

	if len(err) != 0 {
		//return nil, err
		errors.New("Search Product Repo not working")
	}
	return items, nil
}

func (movieRepo *MovieGormRepo) MovieExists(movieName string) bool {
	eve := model.Movie{}
	errs := movieRepo.conn.Find(&eve, "title=?", movieName).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}


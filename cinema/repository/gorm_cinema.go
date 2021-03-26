package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/joocosta/bloctrial/cinema"
	"github.com/joocosta/bloctrial/model"
)

// CommentGormRepo implements menu.CommentRepository interface
type CinemaGormRepo struct {
	conn *gorm.DB
}

// NewHALLGormRepo returns new object of CommentGormRepo
func NewCinemaGormRepo(db *gorm.DB) cinema.CinemaRepository {
	return &CinemaGormRepo{conn: db}
}
func (cllRepo *CinemaGormRepo) Cinemas() ([]model.Cinema, []error) {
	cll := []model.Cinema{}
	errs := cllRepo.conn.Find(&cll).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cll, errs
}

//Cinema retrieves a cinema from the database by its id
func (cllRepo *CinemaGormRepo) Cinema(id uint) (*model.Cinema, []error) {
	cll := model.Cinema{}
	errs := cllRepo.conn.First(&cll, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &cll, errs
}

// StoreComment stores a given customer comment in the database
func (cllRepo *CinemaGormRepo) StoreCinema(cinema *model.Cinema) (*model.Cinema, []error) {
	cll := cinema
	errs := cllRepo.conn.Create(cll).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cll, errs
}

func (cliRepo *CinemaGormRepo) UpdateCinema(cinema *model.Cinema) (*model.Cinema, []error) {
	cli := cinema
	errs := cliRepo.conn.Save(cinema).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cli, errs
}
func (cliRepo *CinemaGormRepo) DeleteCinema(id uint) (*model.Cinema, []error) {
	cli, errs := cliRepo.Cinema(id)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = cliRepo.conn.Delete(cli, id).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return cli, errs
}

func (cliRepo *CinemaGormRepo) CinemaExists(cinema string) bool {
	cli := model.Cinema{}
	errs := cliRepo.conn.Find(&cli, "name=?", cinema).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}
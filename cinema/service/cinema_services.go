package service

import (
	"github.com/joocosta/bloctrial/cinema"
	"github.com/joocosta/bloctrial/model"
)

type CinemaService struct {
	cinemaRepo cinema.CinemaRepository
}

func NewCinemaService(CinemaRepos cinema.CinemaRepository) cinema.CinemaService {

	return &CinemaService{cinemaRepo: CinemaRepos}
}

// HALLs returns all stored comments
func (cs *CinemaService) Cinemas() ([]model.Cinema, []error) {
	cll, errs := cs.cinemaRepo.Cinemas()
	if len(errs) > 0 {
		return nil, errs
	}
	return cll, errs
}

func (cs *CinemaService) Cinema(id uint) (*model.Cinema, []error) {
	cmnts, errs := cs.cinemaRepo.Cinema(id)

	if len(errs) > 0 {
		return nil, errs
	}
	return cmnts, errs
}

func (cs *CinemaService) StoreCinema(cinema *model.Cinema) (*model.Cinema, []error) {
	cmnts, errs := cs.cinemaRepo.StoreCinema(cinema)
	if len(errs) > 0 {
		return nil, errs
	}
	return cmnts, errs
}

func (cs *CinemaService) UpdateCinema(cinema *model.Cinema) (*model.Cinema, []error) {
	cmnts, errs := cs.cinemaRepo.UpdateCinema(cinema)
	if len(errs) > 0 {
		return nil, errs
	}
	return cmnts, errs
}

func (cs *CinemaService) DeleteCinema(id uint) (*model.Cinema, []error) {
	cmnts, errs := cs.cinemaRepo.DeleteCinema(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return cmnts, errs
}

func (cs *CinemaService) CinemaExists(cinemaName string) bool {
	exists := cs.cinemaRepo.CinemaExists(cinemaName)
	return exists
}

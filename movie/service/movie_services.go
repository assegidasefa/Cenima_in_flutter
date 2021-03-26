package service

import (
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/movie"
)

type MovieService struct {
	movieRepo movie.MovieRepository
}

func NewMovieService(mvRepo movie.MovieRepository) movie.MovieService {
	return &MovieService{movieRepo: mvRepo}
}

func (m *MovieService) Movies() ([]model.Movie, []error) {
	mvs, errs := m.movieRepo.Movies()
	if len(errs) > 0 {
		return nil, errs
	}
	return mvs, errs
}

func (m *MovieService) StoreMovie(movie *model.Movie) (*model.Movie, []error) {
	mvs, errs := m.movieRepo.StoreMovie(movie)
	if len(errs) > 0 {
		return nil, errs
	}
	return mvs, errs
}


// Event returns a event object with a given id
func (m *MovieService) Movie(id uint) (*model.Movie, []error) {
	mov, errs := m.movieRepo.Movie(id)

	if len(errs) > 0 {
		return mov, errs
	}

	return mov, nil
}

// UpdateEvent updates a event with new data
func (cs *MovieService) UpdateMovie(movie *model.Movie) (*model.Movie, []error) {

	eve, errs := cs.movieRepo.UpdateMovie(movie)

	if len(errs) > 0 {
		return nil, errs
	}

	return eve, nil
}

// DeleteEvent delete a event by its id
func (cs *MovieService) DeleteMovie(id uint) (*model.Movie, []error) {

	eve, errs := cs.movieRepo.DeleteMovie(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return eve, nil
}
func (cs *MovieService) SearchMovie(index string) ([]model.Movie, error) {
	movies, err := cs.movieRepo.SearchMovie(index)

	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (cs *MovieService) MovieExists(movieName string) bool {
	exists := cs.movieRepo.MovieExists(movieName)
	return exists
}

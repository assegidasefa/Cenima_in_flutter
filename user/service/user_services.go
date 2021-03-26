package service

import (
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/user"
)

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepository user.UserRepository) user.UserService {
	return &UserService{userRepo: userRepository}
}

// User retrieves an application user by its id
func (us *UserService) User(id uint) (*model.User, []error) {
	usr, errs := us.userRepo.User(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// UserByEmail retrieves an application user by its email address
func (us *UserService) UserByEmail(email string) (*model.User, []error) {
	usr, errs := us.userRepo.UserByEmail(email)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

func (us *UserService) UpdateUser(user *model.User) (*model.User, []error) {
	usr, errs := us.userRepo.UpdateUser(user)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

func (us *UserService) UpdateUserAmount(user *model.User, Amount uint) *model.User {
	usr := us.userRepo.UpdateUserAmount(user, Amount)
	return usr
}

//func (us *UserService) UpdateUserAmount(user *model.User, Amount uint) (*model.User, error) {
//	usr, err := us.userRepo.UpdateUserAmount(user, Amount)
//	return usr, nil
//	//usr, err := us.userRepo.UpdateUserAmount(user, Amount)
//	//if err != nil{
//	//	return nil, errors.New("Cant get user amount")
//	//}
//	//return usr, nil
//}

// DeleteUser deletes a given application user
func (us *UserService) DeleteUser(id uint) (*model.User, []error) {
	usr, errs := us.userRepo.DeleteUser(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// StoreUser stores a given application user
func (us *UserService) StoreUser(user *model.User) (*model.User, []error) {
	usr, errs := us.userRepo.StoreUser(user)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}
func (us *UserService) EmailExists(email string) bool {
	exists := us.userRepo.EmailExists(email)
	return exists
}

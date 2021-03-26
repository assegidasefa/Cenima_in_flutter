package user

import "github.com/joocosta/bloctrial/model"

type UserService interface {
	User(id uint) (*model.User, []error)
	UserByEmail(email string) (*model.User, []error)
	UpdateUser(user *model.User) (*model.User, []error)
	DeleteUser(id uint) (*model.User, []error)
	StoreUser(user *model.User) (*model.User, []error)
	UpdateUserAmount(user *model.User, Amount uint) *model.User
	//UpdateUserAmount(user *model.User, Amount uint) (*model.User, error)
	EmailExists(email string) bool

}
type RoleService interface {
	Roles() ([]model.Role, []error)
	Role(id uint) (*model.Role, []error)
	RoleByName(name string) (*model.Role, []error)
	UpdateRole(role *model.Role) (*model.Role, []error)
	DeleteRole(id uint) (*model.Role, []error)
	StoreRole(role *model.Role) (*model.Role, []error)
}

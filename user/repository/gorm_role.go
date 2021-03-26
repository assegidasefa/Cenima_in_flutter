package repository

import (

	"github.com/jinzhu/gorm"
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/user"
)

type RoleGormRepo struct {
	conn *gorm.DB
}

func NewRoleGormRepo(db *gorm.DB) user.RoleRepository {
	return &RoleGormRepo{conn: db}
}

func (roleRepo *RoleGormRepo) Role(id uint) (*model.Role, []error) {
	role := model.Role{}
	errs := roleRepo.conn.First(&role, id).GetErrors()
	return &role, errs
}

func (roleRepo *RoleGormRepo) Roles() ([]model.Role, []error) {
	roles := []model.Role{}
	errs := roleRepo.conn.Find(&roles).GetErrors()
	return roles, errs
}

func (roleRepo *RoleGormRepo) RoleByName(name string) (*model.Role, []error) {
	role := model.Role{}
	errs := roleRepo.conn.Find(&role, "name=?", name).GetErrors()
	return &role, errs
}

func (roleRepo *RoleGormRepo) StoreRole(role *model.Role) (*model.Role, []error) {
	errs := roleRepo.conn.Create(&role).GetErrors()
	return role, errs
}

func (roleRepo *RoleGormRepo) UpdateRole(role *model.Role) (*model.Role, []error) {
	r := role
	errs := roleRepo.conn.Save(r).GetErrors()
	return r, errs
}

func (roleRepo *RoleGormRepo) DeleteRole(id uint) (*model.Role, []error) {
	r, errs := roleRepo.Role(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = roleRepo.conn.Delete(r, id).GetErrors()
	return r, errs
}

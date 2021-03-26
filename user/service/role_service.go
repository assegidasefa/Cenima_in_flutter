package service

import (
	"github.com/joocosta/bloctrial/model"
	"github.com/joocosta/bloctrial/user"
)

type RoleService struct {
	roleRepo user.RoleRepository
}

func NewRoleService(RoleRepo user.RoleRepository) user.RoleService {
	return &RoleService{roleRepo: RoleRepo}
}

func (rs *RoleService) Roles() ([]model.Role, []error) {
	return rs.roleRepo.Roles()
}

func (rs *RoleService) RoleByName(name string) (*model.Role, []error) {
	return rs.roleRepo.RoleByName(name)
}

func (rs *RoleService) Role(id uint) (*model.Role, []error) {
	return rs.roleRepo.Role(id)
}

func (rs *RoleService) UpdateRole(role *model.Role) (*model.Role, []error) {
	return rs.roleRepo.UpdateRole(role)
}

func (rs *RoleService) DeleteRole(id uint) (*model.Role, []error) {
	return rs.roleRepo.DeleteRole(id)
}

func (rs *RoleService) StoreRole(role *model.Role) (*model.Role, []error) {
	return rs.roleRepo.StoreRole(role)
}

package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"math"
	"time"

	"gorm.io/gorm"
)

type RoleUsecase interface {
	GetRoleList(req *domain.RoleListRequest) ([]domain.RoleListItem, *domain.PaginationMeta, error)
	GetRoleByID(id string) (*domain.RoleResponse, error)
	CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error)
	UpdateRole(id string, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	DeleteRole(id string) error
	GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error)
}

type roleUsecase struct {
	roleRepo       repo.RoleRepository
	permissionRepo repo.PermissionRepository
}

func NewRoleUsecase(roleRepo repo.RoleRepository, permissionRepo repo.PermissionRepository) RoleUsecase {
	return &roleUsecase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (u *roleUsecase) GetRoleList(req *domain.RoleListRequest) ([]domain.RoleListItem, *domain.PaginationMeta, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	roles, total, err := u.roleRepo.GetAll(req)
	if err != nil {
		return nil, nil, err
	}

	var roleItems []domain.RoleListItem
	for _, role := range roles {
		permissionsCount := len(role.Permissions)
		roleItems = append(roleItems, *domain.ToRoleListItem(role, permissionsCount))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	return roleItems, meta, nil
}

func (u *roleUsecase) GetRoleByID(id string) (*domain.RoleResponse, error) {
	role, err := u.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, err
	}

	return domain.ToRoleResponse(role), nil
}

func (u *roleUsecase) CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	exists, err := u.roleRepo.IsNameExists(req.Nama)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("role dengan nama tersebut sudah ada")
	}

	invalidIDs, err := u.roleRepo.ValidatePermissions(req.PermissionIDs)
	if err != nil {
		return nil, err
	}

	if len(invalidIDs) > 0 {
		return nil, errors.New("beberapa permission tidak ditemukan atau tidak valid")
	}

	permissions := make([]domain.Permission, len(req.PermissionIDs))
	for i, permissionID := range req.PermissionIDs {
		permissions[i] = domain.Permission{ID: permissionID}
	}

	role := &domain.Role{
		Nama:        req.Nama,
		Deskripsi:   req.Deskripsi,
		Status:      req.Status,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.roleRepo.Create(role); err != nil {
		return nil, err
	}

	createdRole, err := u.roleRepo.GetByID(role.ID)
	if err != nil {
		return nil, err
	}

	return domain.ToRoleResponse(createdRole), nil
}

func (u *roleUsecase) UpdateRole(id string, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	role, err := u.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, err
	}

	if role.Nama != req.Nama {
		exists, err := u.roleRepo.IsNameExists(req.Nama, id)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errors.New("role dengan nama tersebut sudah ada")
		}
	}

	invalidIDs, err := u.roleRepo.ValidatePermissions(req.PermissionIDs)
	if err != nil {
		return nil, err
	}

	if len(invalidIDs) > 0 {
		return nil, errors.New("beberapa permission tidak ditemukan atau tidak valid")
	}

	permissions := make([]domain.Permission, len(req.PermissionIDs))
	for i, permissionID := range req.PermissionIDs {
		permissions[i] = domain.Permission{ID: permissionID}
	}

	role.Nama = req.Nama
	role.Deskripsi = req.Deskripsi
	role.Status = req.Status
	role.Permissions = permissions
	role.UpdatedAt = time.Now()

	if err := u.roleRepo.Update(role); err != nil {
		return nil, err
	}

	updatedRole, err := u.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return domain.ToRoleResponse(updatedRole), nil
}

func (u *roleUsecase) DeleteRole(id string) error {
	_, err := u.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role tidak ditemukan")
		}
		return err
	}

	isUsed, err := u.roleRepo.IsUsedByUsers(id)
	if err != nil {
		return err
	}

	if isUsed {
		return errors.New("role tidak dapat dihapus karena masih digunakan oleh pengguna lain")
	}

	return u.roleRepo.Delete(id)
}

func (u *roleUsecase) GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error) {
	_, err := u.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("role tidak ditemukan")
		}
		return nil, nil, err
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	permissions, total, err := u.roleRepo.GetRolePermissions(roleID, req)
	if err != nil {
		return nil, nil, err
	}

	responses := domain.ToPermissionResponseList(permissions)

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	return responses, meta, nil
}

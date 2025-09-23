package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"math"
	"time"

	"gorm.io/gorm"
)

type PermissionUsecase interface {
	GetPermissionList(req *domain.PermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error)
	GetPermissionByID(id string) (*domain.PermissionResponse, error)
	CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error)
	UpdatePermission(id string, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error)
	DeletePermission(id string) error
}

type permissionUsecase struct {
	permissionRepo repo.PermissionRepository
}

func NewPermissionUsecase(permissionRepo repo.PermissionRepository) PermissionUsecase {
	return &permissionUsecase{
		permissionRepo: permissionRepo,
	}
}

func (u *permissionUsecase) GetPermissionList(req *domain.PermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	permissions, total, err := u.permissionRepo.GetAll(req)
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

func (u *permissionUsecase) GetPermissionByID(id string) (*domain.PermissionResponse, error) {
	permission, err := u.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission tidak ditemukan")
		}
		return nil, err
	}

	return domain.ToPermissionResponse(permission), nil
}

func (u *permissionUsecase) CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error) {
	exists, err := u.permissionRepo.IsNameExists(req.Nama)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("permission dengan nama tersebut sudah ada")
	}

	permission := &domain.Permission{
		Nama:      req.Nama,
		Kategori:  req.Kategori,
		Deskripsi: req.Deskripsi,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return domain.ToPermissionResponse(permission), nil
}

func (u *permissionUsecase) UpdatePermission(id string, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error) {
	permission, err := u.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission tidak ditemukan")
		}
		return nil, err
	}

	if permission.Nama != req.Nama {
		exists, err := u.permissionRepo.IsNameExists(req.Nama, id)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errors.New("permission dengan nama tersebut sudah ada")
		}
	}

	permission.Nama = req.Nama
	permission.Kategori = req.Kategori
	permission.Deskripsi = req.Deskripsi
	permission.UpdatedAt = time.Now()

	if err := u.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return domain.ToPermissionResponse(permission), nil
}

func (u *permissionUsecase) DeletePermission(id string) error {
	_, err := u.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission tidak ditemukan")
		}
		return err
	}

	isUsed, err := u.permissionRepo.IsUsedByRoles(id)
	if err != nil {
		return err
	}

	if isUsed {
		return errors.New("permission tidak dapat dihapus karena masih digunakan oleh role lain")
	}

	return u.permissionRepo.Delete(id)
}

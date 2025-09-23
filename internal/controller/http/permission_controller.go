package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PermissionController struct {
	permissionUsecase usecase.PermissionUsecase
}

func NewPermissionController(permissionUsecase usecase.PermissionUsecase) *PermissionController {
	return &PermissionController{
		permissionUsecase: permissionUsecase,
	}
}

func (ctrl *PermissionController) GetPermissionList(c *fiber.Ctx) error {
	req := domain.NewPermissionListRequest()

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	if kategori := c.Query("kategori"); kategori != "" {
		req.Kategori = &kategori
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			req.PerPage = perPage
		}
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	permissions, meta, err := ctrl.permissionUsecase.GetPermissionList(req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Daftar permission berhasil diambil", fiber.Map{
		"data": permissions,
		"meta": meta,
	})
}

func (ctrl *PermissionController) GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID permission diperlukan", nil)
	}

	permission, err := ctrl.permissionUsecase.GetPermissionByID(id)
	if err != nil {
		if err.Error() == "permission tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail permission berhasil diambil", permission)
}

func (ctrl *PermissionController) CreatePermission(c *fiber.Ctx) error {
	var req domain.CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	permission, err := ctrl.permissionUsecase.CreatePermission(&req)
	if err != nil {
		if err.Error() == "permission dengan nama tersebut sudah ada" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusCreated, "Permission berhasil dibuat", permission)
}

func (ctrl *PermissionController) UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID permission diperlukan", nil)
	}

	var req domain.UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	permission, err := ctrl.permissionUsecase.UpdatePermission(id, &req)
	if err != nil {
		if err.Error() == "permission tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "permission dengan nama tersebut sudah ada" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Permission berhasil diupdate", permission)
}

func (ctrl *PermissionController) DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID permission diperlukan", nil)
	}

	err := ctrl.permissionUsecase.DeletePermission(id)
	if err != nil {
		if err.Error() == "permission tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "permission tidak dapat dihapus karena masih digunakan oleh role lain" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Permission berhasil dihapus", nil)
}

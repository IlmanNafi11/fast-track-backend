package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoleController struct {
	roleUsecase usecase.RoleUsecase
}

func NewRoleController(roleUsecase usecase.RoleUsecase) *RoleController {
	return &RoleController{
		roleUsecase: roleUsecase,
	}
}

func (ctrl *RoleController) GetRoleList(c *fiber.Ctx) error {
	req := domain.NewRoleListRequest()

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
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

	roles, meta, err := ctrl.roleUsecase.GetRoleList(req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendPaginatedResponse(c, fiber.StatusOK, "Daftar role berhasil diambil", roles, *meta)
}

func (ctrl *RoleController) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID role diperlukan", nil)
	}

	role, err := ctrl.roleUsecase.GetRoleByID(id)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail role berhasil diambil", role)
}

func (ctrl *RoleController) CreateRole(c *fiber.Ctx) error {
	var req domain.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	role, err := ctrl.roleUsecase.CreateRole(&req)
	if err != nil {
		if err.Error() == "role dengan nama tersebut sudah ada" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		if err.Error() == "beberapa permission tidak ditemukan atau tidak valid" {
			return helper.SendErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusCreated, "Role berhasil dibuat", role)
}

func (ctrl *RoleController) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID role diperlukan", nil)
	}

	var req domain.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	role, err := ctrl.roleUsecase.UpdateRole(id, &req)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "role dengan nama tersebut sudah ada" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		if err.Error() == "beberapa permission tidak ditemukan atau tidak valid" {
			return helper.SendErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Role berhasil diupdate", role)
}

func (ctrl *RoleController) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID role diperlukan", nil)
	}

	err := ctrl.roleUsecase.DeleteRole(id)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "role tidak dapat dihapus karena masih digunakan oleh pengguna lain" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Role berhasil dihapus", nil)
}

func (ctrl *RoleController) GetRolePermissions(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID role diperlukan", nil)
	}

	req := domain.NewRolePermissionListRequest()

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

	permissions, meta, err := ctrl.roleUsecase.GetRolePermissions(id, req)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Daftar permissions role berhasil diambil", fiber.Map{
		"data": permissions,
		"meta": meta,
	})
}

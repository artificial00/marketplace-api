package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"marketplace-api/internal/models"
	"marketplace-api/internal/service"
	"marketplace-api/pkg/middleware"
	"marketplace-api/pkg/utils"
)

type ListingHandler struct {
	listingService service.ListingServiceInterface
}

func NewListingHandler(listingService service.ListingServiceInterface) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
	}
}

// CreateListing создает новое объявление
// @Summary Создать объявление
// @Description Создает новое объявление для авторизованного пользователя
// @Tags listings
// @Security Bearer
// @Accept json
// @Produce json
// @Param listing body models.CreateListingRequest true "Данные объявления"
// @Success 201 {object} utils.SuccessResponse{data=models.Listing}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings [post]
func (h *ListingHandler) CreateListing(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	var req models.CreateListingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format: "+err.Error())
		return
	}

	listing, err := h.listingService.CreateListing(userID, req)
	if err != nil {
		if err.Error() == "title is required" ||
			err.Error() == "description is required" ||
			err.Error() == "price must be greater than 0" ||
			err.Error() == "invalid image URL format" ||
			err.Error() == "title must be less than 255 characters" ||
			err.Error() == "image URL must be less than 500 characters" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to create listing")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, listing, "Listing created successfully")
}

// GetListings получает список объявлений
// @Summary Получить список объявлений
// @Description Возвращает список объявлений с возможностью фильтрации и пагинации
// @Tags listings
// @Accept json
// @Produce json
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param sort_by query string false "Поле для сортировки" Enums(created_at, price)
// @Param sort_dir query string false "Направление сортировки" Enums(asc, desc)
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} utils.SuccessResponse{data=models.PaginatedListings}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings [get]
func (h *ListingHandler) GetListings(c *gin.Context) {
	var filter models.ListingsFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	var currentUserID *int
	if userID, exists := middleware.GetUserID(c); exists {
		currentUserID = &userID
	}

	listings, err := h.listingService.GetListings(filter, currentUserID)
	if err != nil {
		if err.Error() == "min_price cannot be negative" ||
			err.Error() == "max_price cannot be negative" ||
			err.Error() == "min_price cannot be greater than max_price" ||
			err.Error() == "page must be greater than 0" ||
			err.Error() == "limit must be between 1 and 100" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to get listings")
		return
	}

	utils.SendSuccess(c, http.StatusOK, listings, "")
}

// GetListing получает объявление по ID
// @Summary Получить объявление по ID
// @Description Возвращает детальную информацию об объявлении
// @Tags listings
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} utils.SuccessResponse{data=models.Listing}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings/{id} [get]
func (h *ListingHandler) GetListing(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid listing ID")
		return
	}

	var currentUserID *int
	if userID, exists := middleware.GetUserID(c); exists {
		currentUserID = &userID
	}

	listing, err := h.listingService.GetListingByID(id, currentUserID)
	if err != nil {
		if err.Error() == "listing not found" {
			utils.NotFound(c, "Listing not found")
			return
		}
		if err.Error() == "invalid listing ID" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to get listing")
		return
	}

	utils.SendSuccess(c, http.StatusOK, listing, "")
}

// UpdateListing обновляет объявление
// @Summary Обновить объявление
// @Description Обновляет объявление. Только владелец может редактировать свое объявление
// @Tags listings
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Param listing body models.UpdateListingRequest true "Данные для обновления"
// @Success 200 {object} utils.SuccessResponse{data=models.Listing}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings/{id} [put]
func (h *ListingHandler) UpdateListing(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid listing ID")
		return
	}

	var req models.UpdateListingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format: "+err.Error())
		return
	}

	listing, err := h.listingService.UpdateListing(id, userID, req)
	if err != nil {
		if err.Error() == "listing not found" {
			utils.NotFound(c, "Listing not found")
			return
		}
		if err.Error() == "access denied: you can only edit your own listings" {
			utils.Forbidden(c, "You can only edit your own listings")
			return
		}
		if err.Error() == "invalid listing ID" ||
			err.Error() == "title cannot be empty" ||
			err.Error() == "description cannot be empty" ||
			err.Error() == "price must be greater than 0" ||
			err.Error() == "invalid image URL format" ||
			err.Error() == "title must be less than 255 characters" ||
			err.Error() == "image URL must be less than 500 characters" ||
			err.Error() == "no fields to update" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to update listing")
		return
	}

	utils.SendSuccess(c, http.StatusOK, listing, "Listing updated successfully")
}

// DeleteListing удаляет объявление
// @Summary Удалить объявление
// @Description Удаляет объявление. Только владелец может удалить свое объявление
// @Tags listings
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings/{id} [delete]
func (h *ListingHandler) DeleteListing(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid listing ID")
		return
	}

	err = h.listingService.DeleteListing(id, userID)
	if err != nil {
		if err.Error() == "listing not found" {
			utils.NotFound(c, "Listing not found")
			return
		}
		if err.Error() == "access denied: you can only delete your own listings" {
			utils.Forbidden(c, "You can only delete your own listings")
			return
		}
		if err.Error() == "invalid listing ID" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to delete listing")
		return
	}

	utils.SendSuccess(c, http.StatusOK, nil, "Listing deleted successfully")
}

// GetMyListings получает объявления текущего пользователя
// @Summary Получить мои объявления
// @Description Возвращает список объявлений текущего авторизованного пользователя
// @Tags listings
// @Security Bearer
// @Accept json
// @Produce json
// @Param sort_by query string false "Поле для сортировки" Enums(created_at, price)
// @Param sort_dir query string false "Направление сортировки" Enums(asc, desc)
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} utils.SuccessResponse{data=models.PaginatedListings}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /listings/my [get]
func (h *ListingHandler) GetMyListings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	var filter models.ListingsFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	listings, err := h.listingService.GetUserListings(userID, filter)
	if err != nil {
		if err.Error() == "invalid user ID" ||
			err.Error() == "page must be greater than 0" ||
			err.Error() == "limit must be between 1 and 100" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to get user listings")
		return
	}

	utils.SendSuccess(c, http.StatusOK, listings, "")
}

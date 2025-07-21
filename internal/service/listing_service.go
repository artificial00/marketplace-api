package service

import (
	"fmt"

	"marketplace-api/internal/database/postgres"
	"marketplace-api/internal/models"
	"marketplace-api/pkg/utils"
)

type ListingService struct {
	listingRepo *postgres.ListingRepository
}

func NewListingService(listingRepo *postgres.ListingRepository) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
	}
}

type ListingServiceInterface interface {
	CreateListing(userID int, req models.CreateListingRequest) (*models.Listing, error)
	GetListings(filter models.ListingsFilter, currentUserID *int) (*models.PaginatedListings, error)
	GetListingByID(id int, currentUserID *int) (*models.Listing, error)
	UpdateListing(id int, userID int, req models.UpdateListingRequest) (*models.Listing, error)
	DeleteListing(id int, userID int) error
	GetUserListings(userID int, filter models.ListingsFilter) (*models.PaginatedListings, error)
}

// CreateListing создает новое объявление
func (s *ListingService) CreateListing(userID int, req models.CreateListingRequest) (*models.Listing, error) {
	if err := s.validateCreateListingRequest(req); err != nil {
		return nil, err
	}

	listing, err := s.listingRepo.CreateListing(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	return listing, nil
}

// GetListings получает список объявлений с фильтрацией
func (s *ListingService) GetListings(filter models.ListingsFilter, currentUserID *int) (*models.PaginatedListings, error) {
	filter.SetDefaults()

	if err := s.validateListingsFilter(filter); err != nil {
		return nil, err
	}

	listings, err := s.listingRepo.GetListings(filter, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}

	return listings, nil
}

// GetListingByID получает объявление по ID
func (s *ListingService) GetListingByID(id int, currentUserID *int) (*models.Listing, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid listing ID")
	}

	listing, err := s.listingRepo.GetListingByID(id, currentUserID)
	if err != nil {
		if err.Error() == "listing not found" {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return listing, nil
}

// UpdateListing обновляет объявление
func (s *ListingService) UpdateListing(id, userID int, req models.UpdateListingRequest) (*models.Listing, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid listing ID")
	}

	if err := s.validateUpdateListingRequest(req); err != nil {
		return nil, err
	}

	listing, err := s.listingRepo.UpdateListing(id, userID, req)
	if err != nil {
		if err.Error() == "listing not found" {
			return nil, fmt.Errorf("listing not found")
		}
		if err.Error() == "access denied: not owner" {
			return nil, fmt.Errorf("access denied: you can only edit your own listings")
		}
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	return listing, nil
}

// DeleteListing удаляет объявление
func (s *ListingService) DeleteListing(id, userID int) error {
	if id <= 0 {
		return fmt.Errorf("invalid listing ID")
	}

	err := s.listingRepo.DeleteListing(id, userID)
	if err != nil {
		if err.Error() == "listing not found" {
			return fmt.Errorf("listing not found")
		}
		if err.Error() == "access denied: not owner" {
			return fmt.Errorf("access denied: you can only delete your own listings")
		}
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	return nil
}

// GetUserListings получает объявления пользователя
func (s *ListingService) GetUserListings(userID int, filter models.ListingsFilter) (*models.PaginatedListings, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	filter.SetDefaults()

	if err := s.validateListingsFilter(filter); err != nil {
		return nil, err
	}

	listings, err := s.listingRepo.GetUserListings(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user listings: %w", err)
	}

	return listings, nil
}

// validateCreateListingRequest валидирует запрос на создание объявления
func (s *ListingService) validateCreateListingRequest(req models.CreateListingRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) > 255 {
		return fmt.Errorf("title must be less than 255 characters")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.ImageURL != nil && *req.ImageURL != "" {
		if !utils.ValidateURL(*req.ImageURL) {
			return fmt.Errorf("invalid image URL format")
		}
		if len(*req.ImageURL) > 500 {
			return fmt.Errorf("image URL must be less than 500 characters")
		}
	}
	return nil
}

// validateUpdateListingRequest валидирует запрос на обновление объявления
func (s *ListingService) validateUpdateListingRequest(req models.UpdateListingRequest) error {
	if req.Title != nil {
		if *req.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
		if len(*req.Title) > 255 {
			return fmt.Errorf("title must be less than 255 characters")
		}
	}
	if req.Description != nil && *req.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if req.Price != nil && *req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.ImageURL != nil && *req.ImageURL != "" {
		if !utils.ValidateURL(*req.ImageURL) {
			return fmt.Errorf("invalid image URL format")
		}
		if len(*req.ImageURL) > 500 {
			return fmt.Errorf("image URL must be less than 500 characters")
		}
	}
	return nil
}

// validateListingsFilter валидирует параметры фильтрации
func (s *ListingService) validateListingsFilter(filter models.ListingsFilter) error {
	if filter.MinPrice != nil && *filter.MinPrice < 0 {
		return fmt.Errorf("min_price cannot be negative")
	}
	if filter.MaxPrice != nil && *filter.MaxPrice < 0 {
		return fmt.Errorf("max_price cannot be negative")
	}
	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		return fmt.Errorf("min_price cannot be greater than max_price")
	}
	if filter.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	return nil
}

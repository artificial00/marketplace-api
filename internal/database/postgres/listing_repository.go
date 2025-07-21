package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"marketplace-api/internal/models"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

// CreateListing создает новое объявление
func (r *ListingRepository) CreateListing(userID int, req models.CreateListingRequest) (*models.Listing, error) {
	query := `
		INSERT INTO listings (title, description, image_url, price, user_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, title, description, image_url, price, user_id, created_at, updated_at
	`

	var listing models.Listing
	err := r.db.QueryRow(query, req.Title, req.Description, req.ImageURL, req.Price, userID).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.ImageURL,
		&listing.Price,
		&listing.UserID,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	return &listing, nil
}

// GetListings получает список объявлений с фильтрацией и пагинацией
func (r *ListingRepository) GetListings(filter models.ListingsFilter, currentUserID *int) (*models.PaginatedListings, error) {
	baseQuery := `
		FROM listings l 
		JOIN users u ON l.user_id = u.id
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("l.price >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("l.price <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + " " + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count listings: %w", err)
	}

	selectFields := `
		SELECT l.id, l.title, l.description, l.image_url, l.price, 
		       l.user_id, u.login as user_login, l.created_at, l.updated_at
	`

	orderBy := fmt.Sprintf("ORDER BY l.%s %s", filter.SortBy, filter.SortDir)

	limitOffset := fmt.Sprintf("LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.GetOffset())

	finalQuery := selectFields + " " + baseQuery + " " + whereClause + " " + orderBy + " " + limitOffset

	rows, err := r.db.Query(finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	defer rows.Close()

	var listings []models.Listing
	for rows.Next() {
		var listing models.Listing
		err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.ImageURL,
			&listing.Price,
			&listing.UserID,
			&listing.UserLogin,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}

		if currentUserID != nil && *currentUserID == listing.UserID {
			listing.IsOwner = true
		}

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.PaginatedListings{
		Data:       listings,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetListingByID получает объявление по ID
func (r *ListingRepository) GetListingByID(id int, currentUserID *int) (*models.Listing, error) {
	query := `
		SELECT l.id, l.title, l.description, l.image_url, l.price, 
		       l.user_id, u.login as user_login, l.created_at, l.updated_at
		FROM listings l 
		JOIN users u ON l.user_id = u.id
		WHERE l.id = $1
	`

	var listing models.Listing
	err := r.db.QueryRow(query, id).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.ImageURL,
		&listing.Price,
		&listing.UserID,
		&listing.UserLogin,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	if currentUserID != nil && *currentUserID == listing.UserID {
		listing.IsOwner = true
	}

	return &listing, nil
}

// UpdateListing обновляет объявление
func (r *ListingRepository) UpdateListing(id, userID int, req models.UpdateListingRequest) (*models.Listing, error) {
	checkQuery := "SELECT user_id FROM listings WHERE id = $1"
	var ownerID int
	err := r.db.QueryRow(checkQuery, id).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found")
		}
		return nil, fmt.Errorf("failed to check listing ownership: %w", err)
	}

	if ownerID != userID {
		return nil, fmt.Errorf("access denied: not owner")
	}

	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.ImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, *req.ImageURL)
		argIndex++
	}

	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE listings 
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, image_url, price, user_id, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var listing models.Listing
	err = r.db.QueryRow(query, args...).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.ImageURL,
		&listing.Price,
		&listing.UserID,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	listing.IsOwner = true

	return &listing, nil
}

// DeleteListing удаляет объявление
func (r *ListingRepository) DeleteListing(id, userID int) error {
	checkQuery := "SELECT user_id FROM listings WHERE id = $1"
	var ownerID int
	err := r.db.QueryRow(checkQuery, id).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("listing not found")
		}
		return fmt.Errorf("failed to check listing ownership: %w", err)
	}

	if ownerID != userID {
		return fmt.Errorf("access denied: not owner")
	}

	deleteQuery := "DELETE FROM listings WHERE id = $1"
	result, err := r.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

// GetUserListings получает объявления конкретного пользователя
func (r *ListingRepository) GetUserListings(userID int, filter models.ListingsFilter) (*models.PaginatedListings, error) {
	countQuery := "SELECT COUNT(*) FROM listings WHERE user_id = $1"
	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count user listings: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT l.id, l.title, l.description, l.image_url, l.price, 
		       l.user_id, u.login as user_login, l.created_at, l.updated_at
		FROM listings l 
		JOIN users u ON l.user_id = u.id
		WHERE l.user_id = $1
		ORDER BY l.%s %s
		LIMIT $2 OFFSET $3
	`, filter.SortBy, filter.SortDir)

	rows, err := r.db.Query(query, userID, filter.Limit, filter.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("failed to get user listings: %w", err)
	}
	defer rows.Close()

	var listings []models.Listing
	for rows.Next() {
		var listing models.Listing
		err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.ImageURL,
			&listing.Price,
			&listing.UserID,
			&listing.UserLogin,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user listing: %w", err)
		}

		listing.IsOwner = true

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.PaginatedListings{
		Data:       listings,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

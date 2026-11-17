package models

import "api-tests-template/internal/managers/advertisements/models"

// Favorite описывает запись избранного и вложенное объявление.
type Favorite struct {
	ID            string                      `json:"favorite_id"`
	Advertisement models.AdvertisementDetails `json:"advertisement"`
}

// FavoritesList содержит записи избранного и общее число результатов.
type FavoritesList struct {
	Items []Favorite `json:"items"`
	Total *int       `json:"total"`
}

// BulkAddFavoritesResult описывает, какие объявления удалось добавить.
type BulkAddFavoritesResult struct {
	Added        []string `json:"added"`
	Skipped      []string `json:"skipped"`
	NotFound     []string `json:"not_found"`
	LimitReached bool     `json:"limit_reached"`
}

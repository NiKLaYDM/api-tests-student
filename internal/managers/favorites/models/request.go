package models

// AddFavoriteRequest описывает JSON для добавления объявления в избранное.
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id,omitempty"`
}

// BulkAddFavoritesRequest описывает список объявлений для переноса в избранное.
type BulkAddFavoritesRequest struct {
	AdvertisementIDs []string `json:"advertisement_ids"`
}

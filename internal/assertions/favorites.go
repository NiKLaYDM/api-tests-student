package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// FavoriteCount проверяет, сколько записей избранного относится к объявлению.
func FavoriteCount(t *testing.T, items []models.Favorite, adID string, expected int) {
	t.Helper()
	count := 0
	for _, item := range items {
		if item.Advertisement.ID == adID {
			count++
		}
	}
	require.Equal(t, expected, count, "Количество записей в избранном для объявления %s", adID)
}

// FavoriteAbsent проверяет, что в списке нет ни объявления, ни сохранённого ID записи.
func FavoriteAbsent(t *testing.T, items []models.Favorite, adID, favoriteID string) {
	t.Helper()
	FavoriteCount(t, items, adID, 0)
	for _, item := range items {
		require.NotEqual(t, favoriteID, item.ID, "Запись %s объявления %s осталась в избранном", favoriteID, adID)
	}
}

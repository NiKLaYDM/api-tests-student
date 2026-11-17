package fixtures

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет запись как предусловие и возвращает её ID со ссылкой на объявление.
func AddFavorite(t *testing.T, token, adID string) models.Favorite {
	t.Helper()
	t.Logf("Подготовка: добавляем объявление %s в избранное", adID)
	favorite := favorites.AddFavorite(t, token, adID, http.StatusCreated, FavoriteCleanup(t, token)).Value
	require.NotEmpty(t, favorite.ID, "Сервер не вернул ID записи избранного")
	require.Equal(t, adID, favorite.Advertisement.ID, "Сервер вернул другое объявление")
	return favorite
}

// FavoriteCleanup регистрирует удаление записи после успешной операции добавления.
func FavoriteCleanup(t *testing.T, token string) func(string) {
	t.Helper()
	return func(adID string) {
		t.Cleanup(func() {
			t.Logf("Очистка: убираем объявление %s из избранного", adID)
			favorites.RemoveFavorite(t, token, adID, http.StatusNoContent, http.StatusNotFound)
		})
	}
}

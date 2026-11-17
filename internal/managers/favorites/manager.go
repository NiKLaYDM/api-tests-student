package favorites

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/apiresponse"
	"api-tests-template/internal/managers/favorites/models"
)

// MaxFavorites — максимальное количество записей в избранном одного пользователя.
const MaxFavorites = 50

// AddFavorite передаёт ID объявления для очистки до проверки полного ответа.
func AddFavorite(t *testing.T, token, adID string, expectedStatus int, onAdded func(string)) apiresponse.Result[models.Favorite] {
	t.Helper()
	operation := "Добавление в избранное " + adID
	raw := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: adID})
	result := apiresponse.Read[models.Favorite](t, raw, operation, func(body []byte) {
		if onAdded == nil {
			return
		}
		if adID != "" {
			onAdded(adID)
			return
		}
		var created struct {
			Advertisement struct {
				ID string `json:"id"`
			} `json:"advertisement"`
		}
		if err := json.Unmarshal(body, &created); err == nil && created.Advertisement.ID != "" {
			onAdded(created.Advertisement.ID)
		}
	})
	apiresponse.RequireStatus(t, result, operation, expectedStatus)
	return result
}

// AddFavoritesBulk переносит список объявлений в избранное и проверяет HTTP-статус.
func AddFavoritesBulk(t *testing.T, token string, adIDs []string, expectedStatus int) apiresponse.Result[models.BulkAddFavoritesResult] {
	t.Helper()
	raw := favorites.HttpPostFavoritesBulk(t, token, models.BulkAddFavoritesRequest{AdvertisementIDs: adIDs})
	return apiresponse.ReadChecked[models.BulkAddFavoritesResult](t, raw, "Пакетное добавление в избранное", expectedStatus)
}

// ListFavorites запрашивает полный список избранного и проверяет структуру успешного ответа.
func ListFavorites(t *testing.T, token string, expectedStatus int) apiresponse.Result[models.FavoritesList] {
	t.Helper()
	raw := favorites.HttpGetFavorites(t, token, MaxFavorites)
	result := apiresponse.ReadChecked[models.FavoritesList](t, raw, "Чтение списка избранного", expectedStatus)
	if result.Successful() {
		require.NotNil(t, result.Value.Items, "В ответе избранного отсутствует массив items")
		require.NotNil(t, result.Value.Total, "В ответе избранного отсутствует total")
		require.Len(t, result.Value.Items, *result.Value.Total, "Ответ должен содержать весь список избранного")
	}
	return result
}

// RemoveFavorite удаляет запись избранного и сверяет статус с допустимыми исходами.
func RemoveFavorite(t *testing.T, token, adID string, expectedStatuses ...int) apiresponse.Result[struct{}] {
	t.Helper()
	raw := favorites.HttpDeleteFavorite(t, token, adID)
	return apiresponse.ReadChecked[struct{}](t, raw, "Удаление из избранного "+adID, expectedStatuses...)
}

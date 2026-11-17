package favorites

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite отправляет ID объявления для добавления в избранное.
func HttpPostFavorite(t *testing.T, token string, request models.AddFavoriteRequest) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").JSON(request).Expect(t).End().Response
}

// HttpPostFavoritesBulk отправляет список ID для пакетного добавления в избранное.
func HttpPostFavoritesBulk(t *testing.T, token string, request models.BulkAddFavoritesRequest) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesBulkPath).
		ContentType("application/json").JSON(request).Expect(t).End().Response
}

// HttpGetFavorites получает первую страницу избранного без фильтров.
func HttpGetFavorites(t *testing.T, token string, pageSize int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page_size", strconv.Itoa(pageSize)).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет запись избранного по ID объявления.
func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(fmt.Sprintf(path.FavoritePath, url.PathEscape(adID))).
		Expect(t).End().Response
}

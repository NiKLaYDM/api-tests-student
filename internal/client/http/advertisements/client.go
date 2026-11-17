package advertisements

import (
	"context"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisements отправляет запрос поиска объявлений по названию.
func HttpGetAdvertisements(t *testing.T, ctx context.Context, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		WithContext(ctx).
		Expect(t).
		End().Response
}

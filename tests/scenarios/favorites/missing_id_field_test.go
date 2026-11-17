package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoriteWithMissingAdvertisementID() {
	t := s.T()
	t.Log("Действие: отправляем JSON без поля advertisement_id")
	response := favorites.AddFavorite(t, s.buyer.Token, "", http.StatusBadRequest, fixtures.FavoriteCleanup(t, s.buyer.Token))
	assertions.APIError(t, response.Problem, "invalid_request", "'AdvertisementID' failed on the 'required' tag")
}

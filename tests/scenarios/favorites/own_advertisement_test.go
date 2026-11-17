package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoriteOwnAdvertisement() {
	t := s.T()
	ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	t.Logf("Действие: добавляем собственное объявление %s", ad.ID)
	response := favorites.AddFavorite(t, s.owner.Token, ad.ID, http.StatusForbidden, fixtures.FavoriteCleanup(t, s.owner.Token))
	assertions.APIError(t, response.Problem, "forbidden", "Cannot add own advertisement to favorites")
	s.checkFavoriteCount(s.owner.Token, ad.ID, 0)
}

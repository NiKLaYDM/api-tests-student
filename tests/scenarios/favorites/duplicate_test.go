package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestDuplicateFavorite() {
	t := s.T()
	ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	fixtures.AddFavorite(t, s.buyer.Token, ad.ID)

	t.Logf("Действие: повторно добавляем объявление %s", ad.ID)
	response := favorites.AddFavorite(t, s.buyer.Token, ad.ID, http.StatusConflict, nil)
	assertions.APIError(t, response.Problem, "already_exists", "Advertisement is already in favorites")
	s.checkFavoriteCount(s.buyer.Token, ad.ID, 1)
}

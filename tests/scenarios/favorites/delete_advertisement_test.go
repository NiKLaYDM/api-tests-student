package favorites

import (
	"net/http"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestDeletingAdvertisementRemovesItFromFavorites() {
	t := s.T()
	ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	added := fixtures.AddFavorite(t, s.buyer.Token, ad.ID)
	s.checkFavoriteCount(s.buyer.Token, ad.ID, 1)

	t.Logf("Действие: удаляем само объявление %s от имени владельца", ad.ID)
	advertisements.DeleteAdvertisement(t, s.owner.Token, ad.ID, http.StatusNoContent)
	advertisements.GetAdvertisement(t, ad.ID, http.StatusNotFound)
	s.checkFavoriteRemoved(s.buyer.Token, ad.ID, added.ID)
}

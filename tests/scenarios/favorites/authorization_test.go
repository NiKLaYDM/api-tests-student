package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesWithoutAuthorization() {
	s.Run("Чтение избранного без токена", func() {
		t := s.T()
		t.Log("Читаем избранное без заголовка Authorization")
		response := favorites.ListFavorites(t, "", http.StatusUnauthorized)
		assertions.MissingAuthorization(t, response.Problem)
	})

	s.Run("Добавление в избранное без токена", func() {
		t := s.T()
		ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
		t.Log("Добавляем существующее объявление без заголовка Authorization")
		response := favorites.AddFavorite(t, "", ad.ID, http.StatusUnauthorized, fixtures.FavoriteCleanup(t, s.buyer.Token))
		assertions.MissingAuthorization(t, response.Problem)
	})

	s.Run("Удаление из избранного без токена", func() {
		t := s.T()
		ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
		fixtures.AddFavorite(t, s.buyer.Token, ad.ID)
		s.checkFavoriteCount(s.buyer.Token, ad.ID, 1)

		t.Log("Удаляем существующую запись без заголовка Authorization")
		response := favorites.RemoveFavorite(t, "", ad.ID, http.StatusUnauthorized)
		assertions.MissingAuthorization(t, response.Problem)
		s.checkFavoriteCount(s.buyer.Token, ad.ID, 1)
	})

	s.Run("Пакетное добавление в избранное без токена", func() {
		t := s.T()
		t.Log("Отправляем POST /favorites/bulk без заголовка Authorization")
		response := favorites.AddFavoritesBulk(t, "", []string{missingAdvertisementID}, http.StatusUnauthorized)
		assertions.MissingAuthorization(t, response.Problem)
	})
}

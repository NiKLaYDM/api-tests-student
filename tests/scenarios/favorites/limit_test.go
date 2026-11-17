package favorites

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesLimit() {
	t := s.T()
	initial := favorites.ListFavorites(t, s.buyer.Token, http.StatusOK)
	require.Zero(t, *initial.Value.Total, "Для проверки лимита нужен отдельный тестовый аккаунт с пустым избранным")

	t.Logf("Подготовка: создаём %d записей в избранном", favorites.MaxFavorites-1)
	for range favorites.MaxFavorites - 1 {
		ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
		fixtures.AddFavorite(t, s.buyer.Token, ad.ID)
	}

	t.Logf("Проверка: добавление объявления номер %d проходит", favorites.MaxFavorites)
	lastAllowed := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	fixtures.AddFavorite(t, s.buyer.Token, lastAllowed.ID)

	candidate := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	t.Logf("Действие: пытаемся добавить объявление номер %d", favorites.MaxFavorites+1)
	response := favorites.AddFavorite(t, s.buyer.Token, candidate.ID,
		http.StatusBadRequest, fixtures.FavoriteCleanup(t, s.buyer.Token))
	assertions.APIError(t, response.Problem, "limit_reached", fmt.Sprintf("Favorites limit of %d reached", favorites.MaxFavorites))

	after := favorites.ListFavorites(t, s.buyer.Token, http.StatusOK)
	require.Equal(t, favorites.MaxFavorites, *after.Value.Total, "Размер избранного после достижения лимита")
	assertions.FavoriteCount(t, after.Value.Items, lastAllowed.ID, 1)
	assertions.FavoriteCount(t, after.Value.Items, candidate.ID, 0)
}

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
	buyer := fixtures.NewFavoritesUser(t)
	initial := favorites.ListFavorites(t, buyer.Token, http.StatusOK)
	require.Zero(t, *initial.Value.Total, "У нового пользователя избранное должно быть пустым")

	t.Cleanup(func() {
		remaining := favorites.ListFavorites(t, buyer.Token, http.StatusOK)
		require.Zero(t, *remaining.Value.Total, "Избранное нового пользователя не очистилось после удаления объявлений")
	})

	t.Logf("Подготовка: создаём %d объявлений и добавляем их в избранное одним запросом", favorites.MaxFavorites-1)
	adIDs, removePrepared := fixtures.CreateAdvertisements(t, s.owner.Token, favorites.MaxFavorites-1)
	prepared := favorites.AddFavoritesBulk(t, buyer.Token, adIDs, http.StatusOK).Value
	require.ElementsMatch(t, adIDs, prepared.Added, "Пакетное добавление должно подготовить все объявления")
	require.Empty(t, prepared.Skipped, "При подготовке не должно быть пропущенных объявлений")
	require.Empty(t, prepared.NotFound, "При подготовке не должно быть несуществующих объявлений")
	require.False(t, prepared.LimitReached, "Лимит сработал до 50-й записи")

	beforeBoundary := favorites.ListFavorites(t, buyer.Token, http.StatusOK)
	require.Equalf(t, favorites.MaxFavorites-1, *beforeBoundary.Value.Total, "Перед проверкой границы должно быть %d записей", favorites.MaxFavorites-1)

	t.Logf("Проверка: добавление объявления номер %d проходит", favorites.MaxFavorites)
	lastAllowed := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	fixtures.AddFavorite(t, buyer.Token, lastAllowed.ID)

	candidate := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())
	t.Logf("Действие: пытаемся добавить объявление номер %d", favorites.MaxFavorites+1)
	response := favorites.AddFavorite(t, buyer.Token, candidate.ID,
		http.StatusBadRequest, fixtures.FavoriteCleanup(t, buyer.Token))
	assertions.APIError(t, response.Problem, "limit_reached", fmt.Sprintf("Favorites limit of %d reached", favorites.MaxFavorites))

	after := favorites.ListFavorites(t, buyer.Token, http.StatusOK)
	require.Equal(t, favorites.MaxFavorites, *after.Value.Total, "Размер избранного после достижения лимита")
	assertions.FavoriteCount(t, after.Value.Items, lastAllowed.ID, 1)
	assertions.FavoriteCount(t, after.Value.Items, candidate.ID, 0)
	removePrepared()
}

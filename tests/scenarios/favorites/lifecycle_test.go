package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoriteLifecycle() {
	t := s.T()
	ad := fixtures.CreateAdvertisement(t, s.owner.Token, fixtures.AdvertisementRequest())

	t.Logf("Действие: добавляем объявление %s в избранное", ad.ID)
	added := favorites.AddFavorite(t, s.buyer.Token, ad.ID, http.StatusCreated, fixtures.FavoriteCleanup(t, s.buyer.Token))
	require.NotEmpty(t, added.Value.ID, "Сервер не вернул ID записи избранного")
	require.Equal(t, ad.ID, added.Value.Advertisement.ID, "В ответе другое объявление")
	s.checkFavoriteCount(s.buyer.Token, ad.ID, 1)

	t.Logf("Действие: удаляем запись для объявления %s из избранного", ad.ID)
	favorites.RemoveFavorite(t, s.buyer.Token, ad.ID, http.StatusNoContent)
	s.checkFavoriteRemoved(s.buyer.Token, ad.ID, added.Value.ID)
}

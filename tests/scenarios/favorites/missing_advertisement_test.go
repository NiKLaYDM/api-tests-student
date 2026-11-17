package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

const missingAdvertisementID = "00000000-0000-4000-8000-000000000000"

func (s *TestSuite) TestNonexistentAdvertisement() {
	s.Run("Добавление несуществующего объявления", func() {
		t := s.T()
		checkAdvertisementDoesNotExist(t)
		t.Log("Действие: добавляем отсутствующее объявление в избранное")
		response := favorites.AddFavorite(t, s.buyer.Token, missingAdvertisementID,
			http.StatusNotFound, fixtures.FavoriteCleanup(t, s.buyer.Token))
		assertions.APIError(t, response.Problem, "not_found", "Advertisement not found")
	})

	s.Run("Удаление несуществующего объявления из избранного", func() {
		t := s.T()
		checkAdvertisementDoesNotExist(t)
		t.Log("Действие: удаляем из избранного отсутствующее объявление")
		response := favorites.RemoveFavorite(t, s.buyer.Token, missingAdvertisementID, http.StatusNotFound)
		assertions.APIError(t, response.Problem, "not_found", "Favorite entry not found")
	})
}

// checkAdvertisementDoesNotExist подтверждает предусловие об отсутствии объявления на стенде.
func checkAdvertisementDoesNotExist(t *testing.T) {
	t.Helper()
	t.Logf("Подготовка: проверяем, что объявление %s отсутствует", missingAdvertisementID)
	advertisements.GetAdvertisement(t, missingAdvertisementID, http.StatusNotFound)
}

package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	owner authModels.LoginOkResponse
	buyer authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизуем владельца объявлений и пользователя избранного")
	s.owner = fixtures.LoginOwner(s.T())
	s.buyer = fixtures.LoginFavoritesUser(s.T())
	require.NotEmpty(s.T(), s.owner.User.Id, "ID владельца")
	require.NotEmpty(s.T(), s.buyer.User.Id, "ID пользователя избранного")
	require.NotEqual(s.T(), s.owner.User.Id, s.buyer.User.Id, "Для избранного нужны два разных аккаунта")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// checkFavoriteCount проверяет число записей для объявления в полном избранном.
func (s *TestSuite) checkFavoriteCount(token, adID string, expected int) {
	t := s.T()
	t.Helper()
	t.Logf("Проверка: в избранном должно быть %d записей для объявления %s", expected, adID)
	result := favorites.ListFavorites(t, token, http.StatusOK)
	assertions.FavoriteCount(t, result.Value.Items, adID, expected)
}

// checkFavoriteRemoved проверяет оба ID в полном списке после удаления.
func (s *TestSuite) checkFavoriteRemoved(token, adID, favoriteID string) {
	t := s.T()
	t.Helper()
	t.Logf("Проверка: в полном списке нет объявления %s и записи %s", adID, favoriteID)
	result := favorites.ListFavorites(t, token, http.StatusOK)
	assertions.FavoriteAbsent(t, result.Value.Items, adID, favoriteID)
}

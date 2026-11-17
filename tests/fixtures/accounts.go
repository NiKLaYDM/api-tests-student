package fixtures

import (
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
)

// LoginOwner авторизует аккаунт, который создаёт тестовые объявления.
func LoginOwner(t *testing.T) models.LoginOkResponse {
	t.Helper()
	return auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// LoginFavoritesUser авторизует второй аккаунт, который работает с избранным.
func LoginFavoritesUser(t *testing.T) models.LoginOkResponse {
	t.Helper()
	return auth.Login(t, os.Getenv("TEST_FAVORITES_LOGIN"), os.Getenv("TEST_FAVORITES_PASSWORD"))
}

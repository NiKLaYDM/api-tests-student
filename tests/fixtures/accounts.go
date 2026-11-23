package fixtures

import (
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
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

// NewFavoritesUser создаёт отдельного пользователя для проверки лимита избранного.
func NewFavoritesUser(t *testing.T) models.LoginOkResponse {
	t.Helper()
	t.Log("Подготовка: регистрируем отдельного пользователя для проверки лимита избранного")
	request := models.RegisterRequest{
		Email:     "lab3-limit-" + utils.RandomString(20) + "@example.com",
		Password:  utils.RandomString(32) + "Aa1!",
		FirstName: "FavoritesLimit",
	}
	return auth.Register(t, request)
}

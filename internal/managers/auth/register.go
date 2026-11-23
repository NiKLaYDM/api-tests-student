package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/apiresponse"
	"api-tests-template/internal/managers/auth/models"
)

// Register создаёт тестового пользователя и возвращает его токен и профиль.
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	t.Helper()
	raw := auth.HttpPostAuthRegister(t, request)
	response := apiresponse.ReadChecked[models.LoginOkResponse](t, raw, "Регистрация тестового пользователя", http.StatusCreated)
	require.NotEmpty(t, response.Value.Token, "В ответе регистрации отсутствует токен")
	require.NotEmpty(t, response.Value.User.Id, "В ответе регистрации отсутствует ID пользователя")
	require.Equal(t, strings.ToLower(request.Email), response.Value.User.Email, "Зарегистрирован другой email")
	return response.Value
}

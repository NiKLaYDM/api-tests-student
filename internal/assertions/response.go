package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/helpers/apiresponse"
)

// APIError проверяет код ошибки и ожидаемую часть сообщения.
func APIError(t *testing.T, problem apiresponse.ErrorResponse, code, message string) {
	t.Helper()
	require.Equal(t, code, problem.Error, "Код ошибки для сообщения %q", problem.Message)
	require.Contains(t, problem.Message, message, "Сообщение ошибки %q", code)
}

// MissingAuthorization проверяет ошибку запроса без заголовка Authorization.
func MissingAuthorization(t *testing.T, problem apiresponse.ErrorResponse) {
	t.Helper()
	APIError(t, problem, "unauthorized", "Authorization header required")
}

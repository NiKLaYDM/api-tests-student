package auth

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/auth/models"
)

// HttpPostAuthRegister отправляет запрос регистрации пользователя.
func HttpPostAuthRegister(t *testing.T, request models.RegisterRequest) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Post(path.AuthRegisterPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

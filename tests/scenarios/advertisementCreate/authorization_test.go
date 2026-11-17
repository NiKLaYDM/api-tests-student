package advertisementCreate

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestCreateAdvertisementWithoutAuthorization() {
	s.Run("Заголовок Authorization отсутствует", func() {
		s.T().Log("Отправляем валидную форму без авторизации")
		response := advertisements.CreateAdvertisement(s.T(), "", fixtures.AdvertisementRequest(),
			http.StatusUnauthorized, fixtures.AdvertisementCleanup(s.T(), s.owner.Token))
		assertions.MissingAuthorization(s.T(), response.Problem)
	})

	s.Run("Токен неверный", func() {
		s.T().Log("Отправляем валидную форму с неверным токеном")
		const invalidToken = "incorrect_token"
		response := advertisements.CreateAdvertisement(s.T(), invalidToken, fixtures.AdvertisementRequest(),
			http.StatusUnauthorized, fixtures.AdvertisementCleanup(s.T(), s.owner.Token))
		assertions.APIError(s.T(), response.Problem, "unauthorized", "Invalid or expired token")
	})
}

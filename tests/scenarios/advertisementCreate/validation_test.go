package advertisementCreate

import (
	"fmt"
	"net/http"
	"strings"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

const maxTitleLength = 50

func (s *TestSuite) TestCreateAdvertisementInvalidFields() {
	emptyTitle := ""
	negativePrice := -1.0
	longTitle := strings.Repeat("A", maxTitleLength+1)
	const titleIsRequired = "title is required"

	cases := []struct {
		name    string
		change  func(*models.CreateAdvertisementRequest)
		message string
	}{
		{"Поле title пропущено", func(request *models.CreateAdvertisementRequest) { request.Title = nil }, titleIsRequired},
		{"Поле title пустое", func(request *models.CreateAdvertisementRequest) { request.Title = &emptyTitle }, titleIsRequired},
		{"Поле description пропущено", func(request *models.CreateAdvertisementRequest) { request.Description = nil }, "description is required"},
		{"Фотографии пропущены", func(request *models.CreateAdvertisementRequest) { request.Photos = nil }, "at least 1 photo is required"},
		{"Цена отрицательная", func(request *models.CreateAdvertisementRequest) { request.Price = &negativePrice }, "price must be between 0 and 1000000 inclusive"},
		{"Заголовок длиннее предела", func(request *models.CreateAdvertisementRequest) { request.Title = &longTitle }, fmt.Sprintf("title must not exceed %d characters", maxTitleLength)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := fixtures.AdvertisementRequest()
			tc.change(&request)
			s.T().Log("Отправляем форму с одним нарушением")
			response := advertisements.CreateAdvertisement(s.T(), s.owner.Token, request,
				http.StatusBadRequest, fixtures.AdvertisementCleanup(s.T(), s.owner.Token))
			assertions.APIError(s.T(), response.Problem, "invalid_request", tc.message)
		})
	}
}

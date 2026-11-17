package advertisementSearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	loginData       models.LoginOkResponse
	existingAdTitle string
	searchBeginning string
	searchMiddle    string
	searchEnd       string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя и получение токена")
	s.loginData = fixtures.LoginOwner(s.T())

	base.Precondition("Создаём объявление для поиска подстрок")
	suffix := utils.RandomString(5)
	s.searchBeginning = "Windranger " + suffix
	s.searchMiddle = suffix + " Compass"
	s.searchEnd = "Gale " + suffix
	s.existingAdTitle = s.searchBeginning + " " + fixtures.WindrangerArcanaName + " " + suffix
	request := fixtures.AdvertisementRequest()
	request.Title = &s.existingAdTitle
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request)
	advList.WaitForAdvertisementInSearch(s.T(), s.loginData.Token, s.existingAdTitle, ad.ID)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestSearchByTitle проверяет поиск существующего объявления по части названия.
func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, s.searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items, "Поиск не вернул ни одного объявления")
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
	})
}

// TestSearchByTitle_Substring проверяет положение подстроки и регистр поиска.
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", s.searchBeginning},
		{"середина", s.searchMiddle},
		{"конец", s.searchEnd},
		{"нижний регистр", strings.ToLower(s.searchMiddle)},
		{"верхний регистр", strings.ToUpper(s.searchMiddle)},
		{"смешанный регистр", utils.MixCase(s.searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
		})
	}
}

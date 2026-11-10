package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

// Объявление: https://testboard.avito.com/advertisements/dd8b00cf-d6c6-40cb-b1da-954c40b7ffe7
const (
	existingAdTitle = `Аркана на Shadow Fiend "Demon Eater"`
	searchBeginning = "Аркана"
	searchMiddle    = "Shadow Fiend"
	searchEnd       = `Demon Eater"`
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что константы для поиска заданы")
	require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle — полное название существующего объявления")
	require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning — подстроку из начала названия")
	require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle — подстроку из середины названия")
	require.NotEmpty(s.T(), searchEnd, "задайте searchEnd — подстроку из конца названия")

	base.Precondition("Авторизация пользователя и получение токена")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestSearchByTitle проверяет поиск существующего объявления по части названия.
func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items, "Поиск не вернул ни одного объявления")
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

// TestSearchByTitle_Substring проверяет положение подстроки и регистр поиска.
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}

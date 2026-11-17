package myAdvertisement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"

	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
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

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = fixtures.LoginOwner(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest())
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
	})

	s.Run("В списке только объявления текущего пользователя", func() {
		items := gjson.Get(advertisementsBody, "items").Array()
		require.NotEmpty(s.T(), items, "Список собственных объявлений пуст")
		for _, item := range items {
			id := item.Get("id").String()
			require.Equal(s.T(), s.loginData.User.Id, item.Get("user_id").String(), "Другой владелец у объявления %s", id)
		}
	})

	s.Run("В списке есть подготовленное объявление со своими полями и фото", func() {
		item := gjson.Get(advertisementsBody, fmt.Sprintf(`items.#(id==%q)`, ad.ID))
		require.True(s.T(), item.Exists(), "Созданное объявление %s отсутствует в списке собственных", ad.ID)
		require.NotEmpty(s.T(), item.Get("created_at").String(), "Дата создания пустая")
		require.NotEmpty(s.T(), item.Get("updated_at").String(), "Дата обновления пустая")
		require.Equal(s.T(), ad.Title, item.Get("title").String(), "Название объявления")
		require.Equal(s.T(), ad.Description, item.Get("description").String(), "Описание объявления")
		require.NotEmpty(s.T(), item.Get("photos").Array(), "Фото подготовленного объявления отсутствуют")
	})
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}

package advertisementCreate

import (
	"bytes"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	authModels "api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	owner authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизуем владельца тестовых объявлений")
	s.owner = fixtures.LoginOwner(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestCreateAdvertisementWithAllFields() {
	t := s.T()
	request := fixtures.AdvertisementRequest()

	t.Logf("Действие: создаём объявление %q со всеми полями", *request.Title)
	response := advertisements.CreateAdvertisement(t, s.owner.Token, request, http.StatusCreated, fixtures.AdvertisementCleanup(t, s.owner.Token))
	created := response.Value
	require.NotEmpty(t, created.ID, "Сервер не вернул ID нового объявления")

	s.Run("Чтение по ID возвращает сохранённые поля", func() {
		actual := advertisements.GetAdvertisement(s.T(), created.ID, http.StatusOK).Value
		s.Equal(created.ID, actual.ID, "ID объявления")
		s.Equal(*request.Title, actual.Title, "Заголовок")
		s.Equal(*request.Description, actual.Description, "Описание")
		s.Equal(*request.Price, actual.Price, "Цена")
		s.Equal(*request.Quantity, actual.Quantity, "Количество")
		s.Equal(created.Photos, actual.Photos, "Фотографии объявления")
		s.Equal(s.owner.User.Id, actual.UserID, "Владелец объявления")
	})

	s.Run("Поиск по полному названию находит созданное объявление", func() {
		list := advertisements.WaitForAdvertisementInSearch(s.T(), s.owner.Token, *request.Title, created.ID)
		expected := models.Advertisement{ID: created.ID, Title: *request.Title}
		require.Contains(s.T(), list.Items, expected,
			"Поиск по %q не вернул объявление %s с точным названием", *request.Title, created.ID)
	})

	s.Run("Файл фотографии доступен по URL из списка", func() {
		t := s.T()
		photos := advertisements.GetAdvertisementPhotos(t, s.owner.Token, created.ID)
		require.Len(t, created.Photos, len(request.Photos), "Количество фото в ответе создания")
		require.Len(t, photos, len(request.Photos), "Количество фото в списке")
		require.NotEmpty(t, created.Photos[0].ID, "Отсутствует ID загруженного фото")
		require.NotEmpty(t, photos[0].URL, "Отсутствует URL загруженного фото")
		require.Equal(t, created.Photos[0], photos[0], "Фото в списке отличается от загруженного")

		original, err := png.Decode(bytes.NewReader(request.Photos[0].Data))
		require.NoError(t, err, "Тестовый PNG должен декодироваться")
		downloaded := advertisements.DownloadPhoto(t, photos[0].URL)
		assertions.PhotoEqual(t, original, downloaded)
	})
}

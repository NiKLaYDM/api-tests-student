package advertisements

import (
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/apiresponse"
	"api-tests-template/internal/managers/advertisements/models"
)

// CreateAdvertisement разбирает ответ создания и перед проверкой статуса сообщает ID для очистки.
func CreateAdvertisement(
	t *testing.T, token string, request models.CreateAdvertisementRequest, expectedStatus int, onCreated func(string),
) apiresponse.Result[models.AdvertisementDetails] {
	t.Helper()
	const operation = "Создание объявления"
	raw := advertisements.HttpPostAdvertisement(t, token, request)
	result := apiresponse.Read[models.AdvertisementDetails](t, raw, operation, func(body []byte) {
		if onCreated == nil {
			return
		}
		var created struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body, &created); err == nil && created.ID != "" {
			onCreated(created.ID)
		}
	})
	apiresponse.RequireStatus(t, result, operation, expectedStatus)
	return result
}

// GetAdvertisement получает данные объявления по ID и проверяет ожидаемый статус.
func GetAdvertisement(t *testing.T, id string, expectedStatus int) apiresponse.Result[models.AdvertisementDetails] {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, id)
	return apiresponse.ReadChecked[models.AdvertisementDetails](t, raw, "Чтение объявления "+id, expectedStatus)
}

// GetAdvertisementPhotos возвращает фотографии существующего объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	return apiresponse.ReadChecked[[]models.Photo](t, raw, "Чтение фотографий объявления "+id, http.StatusOK).Value
}

// DownloadPhoto скачивает PNG по URL и проверяет, что тело ответа декодируется.
func DownloadPhoto(t *testing.T, photoURL string) image.Image {
	t.Helper()
	const maxPhotoBytes = 5 << 20
	raw := advertisements.HttpGetPhotoFile(t, photoURL)
	defer raw.Body.Close()
	require.Equal(t, http.StatusOK, raw.StatusCode, "Получение файла фотографии по URL")
	photo, err := png.Decode(io.LimitReader(raw.Body, maxPhotoBytes))
	require.NoError(t, err, "Не удалось декодировать скачанный PNG")
	return photo
}

// DeleteAdvertisement удаляет объявление и сверяет статус с допустимыми исходами.
func DeleteAdvertisement(t *testing.T, token, id string, expectedStatuses ...int) int {
	t.Helper()
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	return apiresponse.ReadChecked[struct{}](t, raw, "Удаление объявления "+id, expectedStatuses...).StatusCode
}

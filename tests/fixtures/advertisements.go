package fixtures

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// WindrangerArcanaName — название арканы, используемой в данных сценариев.
const WindrangerArcanaName = "Compass of the Rising Gale"

//go:embed testdata/compass-of-the-rising-gale.png
var windrangerArt []byte

// AdvertisementRequest готовит валидную форму с уникальным заголовком и артом арканы Windranger.
func AdvertisementRequest() models.CreateAdvertisementRequest {
	title := "Windranger: " + WindrangerArcanaName + " " + utils.RandomString(8)
	description := "Тестовое объявление для лабораторной работы. Аркана Windranger из Dota 2."
	price := 125.50
	quantity := 2
	return models.CreateAdvertisementRequest{
		Title:       &title,
		Description: &description,
		Price:       &price,
		Quantity:    &quantity,
		Photos:      []models.PhotoUpload{{Name: "compass-of-the-rising-gale.png", Data: windrangerArt}},
	}
}

// CreateAdvertisement создаёт объявление как предусловие и регистрирует его очистку.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.AdvertisementDetails {
	t.Helper()
	t.Logf("Подготовка: создаём объявление %q", *request.Title)
	response := advertisements.CreateAdvertisement(t, token, request, http.StatusCreated, AdvertisementCleanup(t, token))
	require.NotEmpty(t, response.Value.ID, "Сервер не вернул ID объявления")
	return response.Value
}

// AdvertisementCleanup регистрирует удаление после получения ID созданного объявления.
func AdvertisementCleanup(t *testing.T, token string) func(string) {
	t.Helper()
	return func(id string) {
		t.Cleanup(func() {
			t.Logf("Очистка: удаляем объявление %s", id)
			status := advertisements.DeleteAdvertisement(t, token, id, http.StatusNoContent, http.StatusNotFound, http.StatusForbidden)
			if status == http.StatusForbidden {
				// TestBoard может вернуть 403 при повторном удалении отсутствующего объявления.
				advertisements.GetAdvertisement(t, id, http.StatusNotFound)
			}
		})
	}
}

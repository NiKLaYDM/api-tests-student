package advertisements

import (
	"context"
	"net/http"
	"slices"
	"testing"
	"time"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/apiresponse"
	"api-tests-template/internal/managers/advertisements/models"
)

const (
	searchTimeout  = 5 * time.Second
	searchInterval = 200 * time.Millisecond
)

// GetAdvertisementsList ищет объявления и возвращает типизированный ответ после проверки HTTP 200.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	return getAdvertisementsList(t, t.Context(), token, text)
}

// getAdvertisementsList выполняет поиск с заданным дедлайном HTTP-запроса.
func getAdvertisementsList(t *testing.T, ctx context.Context, token, text string) models.AdvertisementsList {
	t.Helper()
	raw := advertisements.HttpGetAdvertisements(t, ctx, token, text)
	return apiresponse.ReadChecked[models.AdvertisementsList](t, raw, "Поиск объявлений по "+text, http.StatusOK).Value
}

// WaitForAdvertisementInSearch повторяет успешный поиск до появления объявления или истечения таймаута.
func WaitForAdvertisementInSearch(t *testing.T, token, text, id string) models.AdvertisementsList {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), searchTimeout)
	defer cancel()
	ticker := time.NewTicker(searchInterval)
	defer ticker.Stop()

	for {
		list := getAdvertisementsList(t, ctx, token, text)
		if slices.ContainsFunc(list.Items, func(item models.Advertisement) bool { return item.ID == id }) {
			return list
		}

		select {
		case <-ctx.Done():
			t.Fatalf("Объявление %s не появилось в поиске по %q за %s", id, text, searchTimeout)
		case <-ticker.C:
		}
	}
}

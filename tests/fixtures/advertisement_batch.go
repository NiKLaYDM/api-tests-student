package fixtures

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

// CreateAdvertisements готовит объявления по три одновременно и возвращает ID и функцию их удаления.
func CreateAdvertisements(t *testing.T, token string, count int) ([]string, func()) {
	t.Helper()

	const concurrentRequests = 3
	adIDs := make([]string, count)
	removed := make([]bool, count)
	slots := make(chan struct{}, concurrentRequests)

	prepared := t.Run("Подготовка объявлений", func(group *testing.T) {
		for i := range count {
			group.Run(fmt.Sprintf("Объявление %d", i+1), func(child *testing.T) {
				child.Parallel()
				slots <- struct{}{}
				defer func() { <-slots }()

				onCreated := func(id string) {
					adIDs[i] = id
					t.Cleanup(func() {
						if !removed[i] {
							cleanupAdvertisement(t, token, id)
						}
					})
				}
				response := advertisements.CreateAdvertisement(child, token, AdvertisementRequest(), http.StatusCreated, onCreated)
				require.NotEmpty(child, response.Value.ID, "Сервер не вернул ID подготовленного объявления")
			})
		}
	})
	require.True(t, prepared, "Не удалось подготовить объявления для проверки лимита")

	cleanup := func() {
		t.Helper()
		cleaned := t.Run("Очистка подготовленных объявлений", func(group *testing.T) {
			for i, id := range adIDs {
				group.Run(fmt.Sprintf("Объявление %d", i+1), func(child *testing.T) {
					child.Parallel()
					slots <- struct{}{}
					defer func() { <-slots }()

					cleanupAdvertisement(child, token, id)
					removed[i] = true
				})
			}
		})
		require.True(t, cleaned, "Не удалось удалить подготовленные объявления")
	}
	return adIDs, cleanup
}

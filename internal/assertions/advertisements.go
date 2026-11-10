package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет наличие объявления с точным названием.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	titles := make([]string, 0, len(items))
	for _, item := range items {
		titles = append(titles, item.Title)
	}

	require.Contains(t, titles, title, "Среди найденных объявлений нет объявления с названием %q", title)
}

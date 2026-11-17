package advertisements

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var photoClient = &http.Client{Timeout: 15 * time.Second}

// HttpGetPhotoFile скачивает файл по публичному URL отдельным клиентом без авторизации TestBoard.
func HttpGetPhotoFile(t *testing.T, photoURL string) *http.Response {
	t.Helper()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, photoURL, nil)
	require.NoError(t, err, "Не удалось собрать запрос скачивания фото")
	response, err := photoClient.Do(request)
	require.NoError(t, err, "Не удалось скачать файл фотографии")
	return response
}

package advertisements

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpPostAdvertisement отправляет multipart-форму создания объявления.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	if request.Title != nil {
		require.NoError(t, form.WriteField("title", *request.Title))
	}
	if request.Description != nil {
		require.NoError(t, form.WriteField("description", *request.Description))
	}
	if request.Price != nil {
		require.NoError(t, form.WriteField("price", strconv.FormatFloat(*request.Price, 'f', -1, 64)))
	}
	if request.Quantity != nil {
		require.NoError(t, form.WriteField("quantity", strconv.Itoa(*request.Quantity)))
	}
	for _, photo := range request.Photos {
		headers := make(textproto.MIMEHeader)
		headers.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{
			"name":     "photos",
			"filename": photo.Name,
		}))
		headers.Set("Content-Type", http.DetectContentType(photo.Data))
		part, err := form.CreatePart(headers)
		require.NoError(t, err, "Не удалось создать поле фотографии %q", photo.Name)
		_, err = part.Write(photo.Data)
		require.NoError(t, err, "Не удалось записать фотографию %q", photo.Name)
	}
	require.NoError(t, form.Close(), "Не удалось завершить multipart-форму")

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(form.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement получает объявление по ID в query-параметре.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisementPhotos запрашивает фотографии объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, url.PathEscape(id))).
		Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление от имени владельца.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

package apiresponse

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ErrorResponse описывает стандартное тело ошибки TestBoard.
type ErrorResponse struct {
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

// Result хранит HTTP-статус и разобранное тело ответа.
type Result[T any] struct {
	StatusCode int
	Value      T
	Problem    ErrorResponse
}

// Successful показывает, вернул ли сервер успешный HTTP-статус.
func (r Result[T]) Successful() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// Read читает и закрывает тело ответа. Для 2xx onSuccess вызывается до разбора модели.
func Read[T any](t *testing.T, raw *http.Response, operation string, onSuccess func([]byte)) Result[T] {
	t.Helper()

	require.NotNil(t, raw, "%s: сервер не вернул HTTP-ответ", operation)
	defer raw.Body.Close()
	result := Result[T]{StatusCode: raw.StatusCode}

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "%s: не удалось прочитать тело ответа HTTP %d", operation, raw.StatusCode)

	if result.Successful() && onSuccess != nil {
		onSuccess(body)
	}
	switch {
	case raw.StatusCode == http.StatusNoContent:
		require.Empty(t, body, "%s: ответ 204 должен быть без тела", operation)
	case result.Successful():
		require.NoError(t, json.Unmarshal(body, &result.Value), "%s: не удалось разобрать тело ответа HTTP %d", operation, raw.StatusCode)
	default:
		require.NoError(t, json.Unmarshal(body, &result.Problem), "%s: не удалось разобрать тело ошибки HTTP %d", operation, raw.StatusCode)
	}
	return result
}

// RequireStatus проверяет ожидаемый статус и добавляет текст ошибки сервера в отчёт.
func RequireStatus[T any](t *testing.T, response Result[T], operation string, expected ...int) {
	t.Helper()
	detail := response.Problem
	const message = "%s: HTTP %d, error=%q, message=%q, details=%s"
	if len(expected) == 1 {
		require.Equalf(t, expected[0], response.StatusCode, message,
			operation, response.StatusCode, detail.Error, detail.Message, detail.Details)
		return
	}
	require.Containsf(t, expected, response.StatusCode, message,
		operation, response.StatusCode, detail.Error, detail.Message, detail.Details)
}

// ReadChecked разбирает ответ и сверяет его статус до передачи результата сценарию.
func ReadChecked[T any](t *testing.T, raw *http.Response, operation string, expected ...int) Result[T] {
	t.Helper()
	result := Read[T](t, raw, operation, nil)
	RequireStatus(t, result, operation, expected...)
	return result
}

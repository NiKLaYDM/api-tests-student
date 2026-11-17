package assertions

import (
	"bytes"
	"image"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/require"
)

// PhotoEqual сверяет размеры и цвета пикселей декодированных изображений.
func PhotoEqual(t *testing.T, expected, actual image.Image) {
	t.Helper()

	bounds := expected.Bounds()
	require.Equal(t, bounds, actual.Bounds(), "Размер скачанного изображения")
	expectedPixels := image.NewNRGBA64(bounds)
	actualPixels := image.NewNRGBA64(bounds)
	draw.Draw(expectedPixels, bounds, expected, bounds.Min, draw.Src)
	draw.Draw(actualPixels, bounds, actual, bounds.Min, draw.Src)
	require.True(t, bytes.Equal(expectedPixels.Pix, actualPixels.Pix), "Пиксели скачанного изображения отличаются от исходного PNG")
}

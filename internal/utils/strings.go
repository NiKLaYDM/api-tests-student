package utils

import "unicode"

// MixCase чередует верхний и нижний регистр букв, начиная с верхнего.
// Пробелы, цифры и знаки препинания остаются без изменений.
func MixCase(text string) string {
	characters := []rune(text)
	upper := true
	for i, character := range characters {
		if !unicode.IsLetter(character) {
			continue
		}

		if upper {
			characters[i] = unicode.ToUpper(character)
		} else {
			characters[i] = unicode.ToLower(character)
		}
		upper = !upper
	}

	return string(characters)
}

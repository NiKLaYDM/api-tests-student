package models

// CreateAdvertisementRequest содержит поля multipart-формы создания объявления.
// nil позволяет пропустить поле при проверке валидации.
type CreateAdvertisementRequest struct {
	Title       *string
	Description *string
	Price       *float64
	Quantity    *int
	Photos      []PhotoUpload
}

// PhotoUpload содержит имя и содержимое файла для поля photos.
type PhotoUpload struct {
	Name string
	Data []byte
}

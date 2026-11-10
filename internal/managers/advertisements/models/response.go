package models

// AdvertisementsList содержит список объявлений из ответа API.
type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

// Advertisement содержит поля объявления, необходимые для проверки поиска.
type Advertisement struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

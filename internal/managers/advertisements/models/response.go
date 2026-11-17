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

// AdvertisementDetails описывает объявление в ответах создания и чтения по ID.
type AdvertisementDetails struct {
	Advertisement
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	Photos      []Photo `json:"photos"`
}

// Photo описывает фотографию, сохранённую вместе с объявлением.
type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

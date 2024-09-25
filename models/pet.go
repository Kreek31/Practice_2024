package models

// Pet структура для примера
type Pet struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Species string `json:"species"`
	Breed   string `json:"breed"`
}

package web

type CreateUrlDTO struct {
	LongURL string `json:"long_url"`
}

type ResponseCreateDTO struct {
	ShortURL string `json:"short_url"`
}

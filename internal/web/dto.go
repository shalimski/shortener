package web

type CreateURLDTO struct {
	LongURL string `json:"long_url"`
}

type ResponseCreateDTO struct {
	ShortURL string `json:"short_url"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

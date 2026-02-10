package schemas

type ShortenUrlRequest struct {
	OriginalUrl string `json:"original_url"`
}

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type GetUrlResponse struct {
	OriginalUrl string `json:"original_url"`
}

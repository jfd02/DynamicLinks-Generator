package models

type ShortLinkResponse struct {
	ShortLink string    `json:"shortLink"`
	Warnings  []Warning `json:"warnings"`
}

type LongLinkResponse struct {
	LongLink string `json:"longLink"`
}

type LinkResponse struct {
	ShortLink   string `json:"shortLink"`
	PreviewLink string `json:"previewLink,omitempty"`
}

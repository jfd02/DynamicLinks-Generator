package models

type ShortenLinkRequest struct {
	LongDynamicLink string `json:"longDynamicLink"`
}

type ExchangeShortLinkRequest struct {
	RequestedLink string `json:"requestedLink"`
}

type CreateDynamicLinkRequest struct {
	DynamicLinkInfo DynamicLinkInfo `json:"dynamicLinkInfo"`
	Suffix          Suffix          `json:"suffix,omitempty"`
}

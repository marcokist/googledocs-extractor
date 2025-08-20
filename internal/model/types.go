package model

type ImageData struct {
	URL    string `json:"url"`
	Base64 string `json:"base64"`
}

type ResponseData struct {
	Status  string      `json:"status"`
	Images  []ImageData `json:"images,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ContentBlock struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	URL     string `json:"url,omitempty"`
	Base64  string `json:"base64,omitempty"`
}

type FullContentResponseData struct {
	Status  string         `json:"status"`
	Content []ContentBlock `json:"content,omitempty"`
	Message string         `json:"message,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

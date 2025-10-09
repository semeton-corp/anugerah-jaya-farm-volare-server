package dto

import "io"

type FileResponse struct {
	URL string `json:"url"`
}

type FileDetailResponse struct {
	Metadata map[string]string
	Body     io.Reader
}

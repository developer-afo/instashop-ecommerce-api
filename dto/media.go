package dto

import "io"

type GetMediaDTO struct {
	Body          io.ReadCloser `json:"body"`
	ContentType   *string       `json:"content_type"`
	ContentLength *int64        `json:"content_length"`
}

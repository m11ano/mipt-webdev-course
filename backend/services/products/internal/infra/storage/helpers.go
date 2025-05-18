package storage

import "github.com/gabriel-vasile/mimetype"

// Return mime and extension
func DetectMimeByBytes8KB(data []byte) (string, string) {
	mime := mimetype.Detect(data)
	ext := ""
	if mime.String() != "application/octet-stream" {
		ext = mime.Extension()
	}
	return mime.String(), ext
}

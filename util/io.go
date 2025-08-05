package util

import (
	"io"
	"mime/multipart"
	"os"
)

const (
	path_ticket_image_dir = "assets/image/ticket/"
)

func SaveTicketImage(file multipart.File, filename string) error {
	dst, err := os.Create(path_ticket_image_dir + filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, file)
	return err
}

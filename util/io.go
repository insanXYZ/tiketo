package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

const (
	PathTicketImageDir = "assets/image/ticket/"
)

func SaveTicketImage(file multipart.File, filename string) error {
	dst, err := os.Create(PathTicketImageDir + filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, file)
	return err
}

func DeleteTicketImage(filename string) error {
	return os.Remove(PathTicketImageDir + filename)
}

func GenerateFilenameTicket(id, ext string) string {
	return fmt.Sprintf("%s-%s%s", id, time.Now().Format(time.DateTime), ext)
}

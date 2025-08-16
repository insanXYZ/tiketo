package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	PathTicketImageDir = "assets/image/ticket/"
)

func SaveTicketImage(file multipart.File, filename string) {
	dst, err := os.Create(filepath.Join(PathTicketImageDir, filename))
	if err == nil {
		io.Copy(dst, file)
	}
}

func DeleteTicketImage(filename string) error {
	return os.Remove(PathTicketImageDir + filename)
}

func GenerateFilenameTicket(id, ext string) string {
	return fmt.Sprintf("%s-%s%s", id, time.Now().Format(time.DateTime), ext)
}

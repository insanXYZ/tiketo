package main

import (
	"fmt"
	"os"
)

func buildDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv())
}

func main() {
	// sql.Open("postgres")
}

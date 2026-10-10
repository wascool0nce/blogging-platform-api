package main

import (
	"fmt"
	"os"

	"github.com/wascool0nce/blogging-platform-api/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

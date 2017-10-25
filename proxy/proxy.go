package proxy

import (
	"fmt"
	"os"
)

func Set()  {
	fmt.Println("\x1b[33;1m[HP] Set https proxy\x1b[0m")
	os.Setenv("HTTPS_PROXY", "https://180.247.205.229:8080")
}
package proxy

import (
	"fmt"
	"os"
)

func Set()  {
	fmt.Print("Set https proxy \n\n")
	os.Setenv("HTTPS_PROXY", "https://193.37.152.6:3128")
}
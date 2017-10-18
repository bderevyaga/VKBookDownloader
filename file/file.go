package file

import (
	"io"
	"net/http"
	"os"
)

const pach = "./buffer/"

func Download(filepath string, url string) (err error) {
	out, err := os.Create(pach + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

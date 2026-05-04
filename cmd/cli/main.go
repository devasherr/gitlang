package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadBinary(url, path string) error {
	// TODO: path will be gloabl, saved on PATH
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status)
	}

	// create dir if not exist
	if err = os.Mkdir(".gbin", 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = os.Chmod(path, 0755); err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	return err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments provided")
	}

	switch os.Args[1] {
	case "init":
		_, err := os.ReadDir(".git")
		if err != nil {
			log.Fatal("unable to locate .git, make sure current project is tracked by git")
		}

		if err := downloadBinary("https://github.com/devasherr/gitlang/releases/download/v0.1.0/gitlang-dispatcher-linux-amd64", "./.gbin/dispatcher"); err != nil {
			log.Fatalf("failed to download dispatcher: %s", err.Error())
		}
	default:
		log.Fatalf("unknown argument: %s", os.Args[1])
	}
}

package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/zc310/utils"

	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println("file\tMD5\tSHA1")
	if len(os.Args) == 1 {
		hashdir(".")
		return
	}
	for _, file := range os.Args[1:] {
		if ok, err := utils.IsDirectory(file); err == nil && ok {
			hashdir(file)
		} else {
			hashfile(file)
		}
	}
}
func hashdir(a string) {
	filepath.Walk(a, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(path, "\t", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		hashfile(path)
		return nil
	})
}
func hashfile(a string) {
	fmt.Print(a)
	fmt.Print("\t")
	f, err := os.Open(a)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	h1 := md5.New()
	h2 := sha1.New()

	if _, err := io.Copy(io.MultiWriter(h1, h2), f); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(hex.EncodeToString(h1.Sum(nil)))
	fmt.Print("\t")
	fmt.Println(hex.EncodeToString(h2.Sum(nil)))

}

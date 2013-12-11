package main

import (
	"crypto/rand"
	"io"
	"os"
	"log"
	"io/ioutil"
	"path/filepath"
	"encoding/base64"
)

const basePath = "files/"
const baseUrl = "static.whooplist.com/"

func writeFileBase64(filename, dataEncoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(dataEncoded)
	if err != nil {
		return "", err
	}
	return writeFile(filename, data)
}

func writeFile(filename string, data []byte) (path string, err error) {
        randomData := make([]byte, 24)
        n, err := io.ReadFull(rand.Reader, randomData)

        if n != len(randomData) || err != nil {
                return
        }

        randString := base64.StdEncoding.EncodeToString(randomData)

        err = os.Mkdir(filepath.Join(basePath, randString), 0666)
        if err != nil {
                return
        }

        err = ioutil.WriteFile(filepath.Join(basePath, randString, filename), data, 0666)

        if err != nil {
                return
        }

        path = filepath.Join(baseUrl, randString, filename)
        log.Print("Saved: ", path)

        return
}


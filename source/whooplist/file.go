package whooplist

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const basePath = "files/"
const baseUrl = "static.whooplist.com/"

const infraName = "assets/"
const userCreatedName = "user/"

const stickyVersion = byte(1)

func WriteFileBase64(filename string, dataEncoded *string, infrastructure bool) (string, error) {
	data, err := base64.StdEncoding.DecodeString(*dataEncoded)
	if err != nil {
		return "", err
	}
	return WriteFile(filename, data, infrastructure)
}

func WriteFile(filename string, data []byte, infra bool) (path string, err error) {
	nameData := make([]byte, 16)
	if infra {
		hasher := sha512.New()
		hasher.Write([]byte(filename))
		hasher.Write([]byte(data))
		nameData = hasher.Sum(nil)
		nameData = nameData[:16]
	} else {
		n, err := io.ReadFull(rand.Reader, nameData)
		if n != len(nameData) || err != nil {
			return "", err
		}
	}

	uniqueDir := base32.StdEncoding.EncodeToString(nameData)
	uniqueDir = uniqueDir[:len(uniqueDir)-6]

	var subDir string
	if infra {
		subDir = infraName
	} else {
		subDir = userCreatedName
	}

	dir := filepath.Join(basePath, subDir, uniqueDir)
	file := filepath.Join(dir, filename)

	err = os.Mkdir(dir, 0744)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0644)

	if err != nil {
		return
	}

	log.Print("Saved: ", file)

	return

}

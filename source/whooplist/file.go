package whooplist

import (
	"io"
	"os"
	"log"
	"strings"
	"io/ioutil"
	"path/filepath"
	"encoding/base64"
	"crypto/rand"
	"crypto/sha512"
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

func WriteFile(filename string, data []byte, infrastructure bool) (path string, err error) {
	hasher := sha512.New()
	hasher.Write([]byte(filename))
	hasher.Write([]byte(data))
	filenameHash := hasher.Sum(nil)

        secondaryData := make([]byte, 24)

	if infrastructure {
		secondaryData[0] = stickyVersion
	} else {
		n, err := io.ReadFull(rand.Reader, secondaryData)
		if n != len(secondaryData) || err != nil {
			return "", err
		}
	}

        randString := base64.StdEncoding.EncodeToString(append(filenameHash, secondaryData...))
	randString = strings.Replace(randString, "/", "-", -1)

	var subdir string
	if infrastructure {
		subdir = infraName
	} else  {
		subdir = userCreatedName
	}

        err = os.Mkdir(filepath.Join(basePath, subdir, randString), 0700)
        if err != nil {
                return
        }

        err = ioutil.WriteFile(filepath.Join(basePath, subdir, randString, filename), data, 0600)

        if err != nil {
                return
        }



        path = filepath.Join(baseUrl, subdir, randString, filename)
        log.Print("Saved: ", path)

        return
}


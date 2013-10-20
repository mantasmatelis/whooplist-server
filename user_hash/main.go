package main

import (
	"encoding/base64"
	"os"
	"strings"
	"fmt"
	"bufio"
	"code.google.com/p/go.crypto/scrypt"
)


func hash(username, password string) (hash string, err error) {
        secret := "aYdZYlE9ybGXn5CldvQ3f/shKxNshtAOvqDlaw/wbUBHwc5r9zBal9hf9CDkGxSgddAMtNm+uz1G"
        secretData, err := base64.StdEncoding.DecodeString(secret)

        if err != nil {
                return
        }

        salt := make([]byte, len(secretData) + len(username))

	fmt.Println(salt)

        copy(salt, secretData)
        copy(salt[len(secretData):], []byte(strings.ToLower(username)))

	fmt.Println(salt)

        hash_data, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)

        hash = base64.StdEncoding.EncodeToString(hash_data)

        return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Println("Enter username")
	username_raw, _ := in.ReadString('\n')
	username := strings.TrimSpace(username_raw)

	fmt.Println("Enter password")
	password_raw, _ := in.ReadString('\n')
	password := strings.TrimSpace(password_raw)

	fmt.Println("Your hash is:")
	fmt.Println(hash(username, password))
}

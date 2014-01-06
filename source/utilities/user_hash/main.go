package main

import (
	"os"
	"strings"
	"fmt"
	"bufio"
	"../../whooplist"
)


func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Println("Enter username")
	username_raw, _ := in.ReadString('\n')
	username := strings.TrimSpace(username_raw)

	fmt.Println("Enter password")
	password_raw, _ := in.ReadString('\n')
	password := strings.TrimSpace(password_raw)

	fmt.Println("Your hash is:")

	hash, _ := whooplist.Hash(username, password)

	fmt.Println(hash)
}

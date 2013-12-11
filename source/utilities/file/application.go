package main

func main() {
	baseUrl := "https://files.whooplist.com/"
	basePath := "files/"

	rpcServe(baseUrl, basePath)
	httpServe(baseUrl, basePath)
}

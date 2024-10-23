package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	var clientID string
	fmt.Print("Please enter your Github Client ID: ")
	fmt.Scan(&clientID)
	fmt.Println(clientID)

	data := url.Values{}

	githubUrl := "https://github.com/login/device/code"

	http.PostForm(githubUrl, data)

}

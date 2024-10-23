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
	data.Set("client_id", clientID)
	data.Set("scope", "repo")

	githubUrl := "https://github.com/login/device/code"
	resp, err := http.PostForm(githubUrl, data)

	if err != nil {
		fmt.Println("Error sending request", err)
		return
	}

	http.PostForm(githubUrl, data)

}

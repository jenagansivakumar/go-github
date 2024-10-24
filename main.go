package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type DeviceResponse struct {
}

func main() {
	var clientID string
	fmt.Print("Please enter your Github Client ID: ")
	fmt.Scan(&clientID)
	fmt.Println(clientID)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", "repo")
	fmt.Println("Form Data:", data.Encode())

	githubUrl := "https://github.com/login/device/code"
	resp, err := http.PostForm(githubUrl, data)
	fmt.Printf("Response Status: %s \n", resp.Status)

	if err != nil {
		fmt.Println("Error sending request", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Response received from GitHub")

}

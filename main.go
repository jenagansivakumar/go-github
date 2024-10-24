package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DeviceResponse struct {
	DeviceCode      string
	UserCode        string
	VerificationUri string
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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(string(body))

	fmt.Println("Response received from GitHub")

}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DeviceResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
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
	if err != nil {
		fmt.Println("Error sending request", err)
		return
	}
	fmt.Printf("Response Status: %s \n", resp.Status)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var deviceResponse DeviceResponse
	err = json.Unmarshal(body, &deviceResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Response received from GitHub")
}

package main

import (
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

func pollForAccessTokens(deviceCode string, clientID string) {
	data := url.Values{}
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientID)

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

	fmt.Println("Raw Response Body:", string(body))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	var deviceResponse DeviceResponse
	deviceResponse.DeviceCode = values.Get("device_code")
	deviceResponse.UserCode = values.Get("user_code")
	deviceResponse.VerificationUri = values.Get("verification_uri")

	fmt.Printf("Device Code: %s \n", deviceResponse.DeviceCode)
	fmt.Printf("User Code: %s \n", deviceResponse.UserCode)
	fmt.Printf("Verification URI: %s \n", deviceResponse.VerificationUri)
	fmt.Println("Response received from GitHub")
}

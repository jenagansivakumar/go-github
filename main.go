package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DeviceResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
}

type RepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
}

func extractAccessToken(body string) (string, error) {
	values, err := url.ParseQuery(body)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}
	accessToken := values.Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("access token not found in response")
	}
	return accessToken, nil
}

func pollForAccessTokens(deviceCode string, clientID string) (string, error) {
	data := url.Values{}
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientID)

	maxAttempts := 10
	interval := 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
		if err != nil {
			return "", fmt.Errorf("error posting form")
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("error retrieving body: %v", err)
		}
		fmt.Printf("Polling attempt %d: %s\n", attempt, string(body))

		if strings.Contains(string(body), "access_token") {
			accessToken, err := extractAccessToken(string(body))
			if err != nil {
				return "", err
			}
			return accessToken, nil
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return "", fmt.Errorf("access token not received within the maximum attempts")
}

func createRepo(token, repoName, description string, private bool) error {
	repoRequest := RepoRequest{
		Name:        repoName,
		Description: description,
		Private:     private,
	}

	jsonData, err := json.Marshal(repoRequest)
	if err != nil {
		fmt.Println("Error parsing json: ", err)
	}

	resp, err := http.Post("https://api.github.com/user/repos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error posting to endpoint: ", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error retrieving body: %v ", err)
	}
	fmt.Println(body)
	return nil
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
		fmt.Printf("Error sending request: %v", err)
		return
	}
	fmt.Printf("Response Status: %s \n", resp.Status)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
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

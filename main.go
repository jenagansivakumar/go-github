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
	DeviceCode      string
	UserCode        string
	VerificationUri string
	Interval        int
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

func pollForAccessTokens(deviceCode, clientID string, interval int) (string, error) {
	data := url.Values{}
	data.Set("device_code", deviceCode)
	data.Set("client_id", clientID)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	maxAttempts := 10

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
		return fmt.Errorf("error parsing json: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/user/repos", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error posting to endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error retrieving body: %v", err)
	}
	fmt.Println("Response: ", string(body))
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

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	deviceResponse := DeviceResponse{
		DeviceCode:      values.Get("device_code"),
		UserCode:        values.Get("user_code"),
		VerificationUri: values.Get("verification_uri"),
		Interval:        5,
	}

	fmt.Printf("Device Code: %s \n", deviceResponse.DeviceCode)
	fmt.Printf("User Code: %s \n", deviceResponse.UserCode)
	fmt.Printf("Verification URI: %s \n", deviceResponse.VerificationUri)
	fmt.Println("Please go to the above URL and enter the user code to authenticate.")

	accessToken, err := pollForAccessTokens(deviceResponse.DeviceCode, clientID, deviceResponse.Interval)
	if err != nil {
		fmt.Println("Error retrieving access token:", err)
		return
	}

	err = createRepo(accessToken, "Test-repo", "This is a test repo", true)
	if err != nil {
		fmt.Println("Error creating repo:", err)
	}
	fmt.Println("Successfully created repo!")
}

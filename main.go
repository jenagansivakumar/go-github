package main

import (
	"bytes"
	"encoding/json"
	"flag"
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

	maxAttempts := 15
	interval := 2

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
	url := "https://api.github.com/user/repos"

	repoData := RepoRequest{
		Name:        repoName,
		Description: description,
		Private:     private,
	}
	jsonData, err := json.Marshal(repoData)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create repo: %s", resp.Status)
	}

	fmt.Println("Repository created successfully.")
	return nil
}

func main() {
	var clientID string
	fmt.Print("Please enter your Github Client ID: ")
	fmt.Scan(&clientID)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", "repo")

	githubUrl := "https://github.com/login/device/code"
	resp, err := http.PostForm(githubUrl, data)
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

	accessToken, err := pollForAccessTokens(deviceResponse.DeviceCode, clientID)
	if err != nil {
		fmt.Println("Error receiving access token:", err)
		return
	}

	repoName := flag.String("name", "", "Name of the repository")
	description := flag.String("desc", "", "Description of the repository")
	private := flag.Bool("private", false, "Set repository as private")
	flag.Parse()

	err = createRepo(accessToken, *repoName, *description, *private)
	if err != nil {
		fmt.Println("Error creating repository:", err)
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	response, err := http.Get("https://api.github.com")
	if err != nil {
		fmt.Println("Error!")
	}

	defer response.Body.Close()

	fmt.Printf(response.Status)

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error fetching body")
	}

	fmt.Println(string(body))
}

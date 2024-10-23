package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	response, err := http.Get("https://api.github.com")
	if err != nil {
		fmt.Println("Error!", err)
	}

	defer response.Body.Close()

	fmt.Println(response.Status)

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error fetching body: ", err)
	}

	fmt.Println(string(body))
}

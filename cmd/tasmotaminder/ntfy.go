package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func notify(title string, tags string, body string) error {
	token, url := getNtfyConfig()

	if url == "" {
		return nil
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Title", title)
	req.Header.Set("Tags", tags)

	response, responseErr := http.DefaultClient.Do(req)
	if responseErr != nil {
		return responseErr
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Println("Error closing response body: %s", closeErr)
		}
	}()

	if response.StatusCode >= 400 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
		} else {
			log.Printf("Response Status: %d, Response Body: %s", response.StatusCode, string(responseBody))
		}
	}

	return nil
}

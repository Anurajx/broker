package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func TSproxyUpstream(targetURL string, apiKey string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("Upstream request creation error: %v", err)
		// writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Upstream request execution error: %v", err)
		// writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusBadGateway, fmt.Errorf("failed to execute upstream request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Upstream body read error: %v", err)
		// writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusBadGateway, fmt.Errorf("failed to execute upstream request: %w", err)
	}

	return body, resp.StatusCode, nil

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(resp.StatusCode)
	// w.Write(body)
}

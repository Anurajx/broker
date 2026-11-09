package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	BaseURL   = "https://backend.delhimetrorail.com/api/v2/en"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/139.0.0.0 Safari/537.36"
	Referer   = "https://www.delhimetrorail.com/"
	Origin    = "https://www.delhimetrorail.com"
)

// Helper to write JSON error messages
// func writeJSON(statusCode int, data interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(data)
// }

// MetroUpstreamProxy creates the HTTP client inline, attaches Delhi Metro headers,
// executes the request with a context timeout, and writes the response directly to w.
func MetroUpstreamProxy(targetURL string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("Upstream request creation error: %v", err)
		// writeJSON( http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request %w", err)
	}

	// Set required Delhi Metro API headers
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", Referer)
	req.Header.Set("Origin", Origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Upstream request execution error: %v", err)
		// writeJSON( http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to execute request %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Upstream body read error: %v", err)
		// writeJSON( http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response %w", err)
	}

	return body, resp.StatusCode, nil
	// w.Header().Set("Content-Type", "application/json") //ADD THIS TO  FUNCTION CALL
	// w.WriteHeader(resp.StatusCode)
	// w.Write(body)
}

// -----------------------------------------------------------------------------
// Route Handlers (HTTP Handlers forwarding to MetroUpstreamProxy)
// -----------------------------------------------------------------------------

// HandleGetLineList proxies requests for the line list.
// func HandleGetLineList(r *http.Request) {
// 	targetURL := fmt.Sprintf("%s/line_list", BaseURL)
// 	MetroUpstreamProxy(targetURL)
// }

// HandleGetLineInfo proxies line info requests (expects lineCode, e.g., "LN1").
// func HandleGetLineInfo(lineCode string) {
// 	targetURL := fmt.Sprintf("%s/metro_line/%s", BaseURL, lineCode)
// 	MetroUpstreamProxy(targetURL)
// }

// HandleGetStationsByLine proxies station lists by line (e.g., "LN1").
// func HandleGetStationsByLine(lineCode string) {
// 	targetURL := fmt.Sprintf("%s/station_by_line_linepage/%s", BaseURL, lineCode)
// 	MetroUpstreamProxy(targetURL)
// }

// HandleGetStationInfo proxies station information requests (e.g., "NBAA").
// func HandleGetStationInfo(stationCode string) {
// 	targetURL := fmt.Sprintf("%s/station/%s", BaseURL, stationCode)
// 	MetroUpstreamProxy(targetURL)
// }

// HandleGetRoute proxies route lookup requests.
// Example ISO timestamp format: 2026-01-03T20:26:54.799
// func HandleGetRoute(fromCode, toCode, routeType string, departureTime time.Time) {
// 	timeFormatted := departureTime.Format("2006-01-02T15:04:05.000")
// 	targetURL := fmt.Sprintf("%s/station_route/%s/%s/%s/%s", BaseURL, fromCode, toCode, routeType, timeFormatted)
// 	MetroUpstreamProxy(targetURL)
// }

// // HandleGetFareWithRoute proxies fare calculation requests.
// func HandleGetFareWithRoute(fromCode, toCode, routeType string) {
// 	targetURL := fmt.Sprintf("%s/new_fare_with_route/%s/%s/%s/", BaseURL, fromCode, toCode, routeType)
// 	MetroUpstreamProxy(targetURL)
// }

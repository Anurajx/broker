package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	busEtaKey := os.Getenv("BUS_ETA_KEY")
	journeyPlannerKey := os.Getenv("JOURNEY_PLANNER_KEY")
	upstreamLink := os.Getenv("UPSTREAM_URL_KEY")

	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})

	// Schedule API
	mux.HandleFunc("GET /api/schedule", func(w http.ResponseWriter, r *http.Request) {
		stopID := r.URL.Query().Get("stop_id")
		if stopID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing stop_id parameter in request"})
			return
		}

		if busEtaKey == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/eta-calculator/stop/etas?stop_id=%s", upstreamLink, url.QueryEscape(stopID))
		proxyUpstream(w, upstreamURL, busEtaKey)
	})

	// Get Stops API
	mux.HandleFunc("GET /api/get_stops", func(w http.ResponseWriter, r *http.Request) {
		mode := r.URL.Query().Get("mode")
		if mode == "" {
			mode = "bus"
		}

		if journeyPlannerKey == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/journey-planner/get_stops?mode=%s", upstreamLink, url.QueryEscape(mode))
		proxyUpstream(w, upstreamURL, journeyPlannerKey)
	})

	// Route API
	mux.HandleFunc("GET /api/route", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		source := query.Get("src")
		sourceType := query.Get("src_type")
		if sourceType == "" {
			sourceType = "place"
		}

		destination := query.Get("dst")
		destinationType := query.Get("dst_type")
		if destinationType == "" {
			destinationType = "place"
		}

		transitMode := query.Get("mode")
		timeVal := query.Get("time")

		if source == "" || destination == "" || transitMode == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing required parameters in request"})
			return
		}

		if journeyPlannerKey == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf(
			"%s/journey-planner/multi_modal?src=%s&src_type=%s&dst=%s&dst_type=%s&mode=%s&time=%s",
			upstreamLink,
			url.QueryEscape(source),
			url.QueryEscape(sourceType),
			url.QueryEscape(destination),
			url.QueryEscape(destinationType),
			url.QueryEscape(transitMode),
			url.QueryEscape(timeVal),
		)

		proxyUpstream(w, upstreamURL, journeyPlannerKey)
	})

	// Wrap router with CORS middleware
	handler := corsMiddleware(mux)

	addr := "0.0.0.0:" + port
	fmt.Printf("Transit broker running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func proxyUpstream(w http.ResponseWriter, targetURL string, apiKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("Upstream request creation error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return
	}

	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Upstream request execution error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Upstream body read error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Upstream request encountered error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

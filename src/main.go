package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

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
		body, respCode, err := TSproxyUpstream(upstreamURL, busEtaKey)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
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
		body, respCode, err := TSproxyUpstream(upstreamURL, journeyPlannerKey)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	// Get metro station Info
	// mux.HandleFunc("GET /api/metro_station", func(w http.ResponseWriter, r *http.Request) {
	// 	station_id := r.URL.Query().Get("id")
	// 	if station_id == "" {
	// 		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "station id not "})
	// 	}

	// 	// upstreamURL := fmt.Sprintf("%s/journey-planner/get_stops?mode=%s", upstreamLink, url.QueryEscape(mode))
	// 	// TSproxyUpstream( upstreamURL, journeyPlannerKey)
	// 	HandleGetStationInfo(w, station_id)

	// })

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

		body, respCode, err := TSproxyUpstream(upstreamURL, journeyPlannerKey)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	// Wrap router with CORS middleware
	handler := corsMiddleware(mux)

	addr := "0.0.0.0:" + port
	fmt.Printf("Transit broker running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

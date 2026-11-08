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
	TS_Base_URL := os.Getenv("TS_BASE_URL")
	Metro_Base_URL := os.Getenv("DMRC_BASE_URL")
	// UserAgent := os.Getenv("USER_AGENT")
	// Referer := os.Getenv("REFERER")
	// Origin := os.Getenv("ORIGIN")

	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})

	// Schedule API
	mux.HandleFunc("GET /api/bus_schedule", func(w http.ResponseWriter, r *http.Request) {
		stopID := r.URL.Query().Get("stop_id")
		if stopID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing stop_id parameter in request"})
			return
		}

		if busEtaKey == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/eta-calculator/stop/etas?stop_id=%s", TS_Base_URL, url.QueryEscape(stopID))
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
	mux.HandleFunc("GET /api/get_bus_stops", func(w http.ResponseWriter, r *http.Request) {
		mode := r.URL.Query().Get("mode")
		if mode == "" {
			mode = "bus"
		}

		if journeyPlannerKey == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/journey-planner/get_stops?mode=%s", TS_Base_URL, url.QueryEscape(mode))
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

	// 	// upstreamURL := fmt.Sprintf("%s/journey-planner/get_stops?mode=%s", TS_Base_URL, url.QueryEscape(mode))
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
			TS_Base_URL,
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

	mux.HandleFunc("GET /api/get_metro_lines", func(w http.ResponseWriter, r *http.Request) {
		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/line_list", Metro_Base_URL)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	mux.HandleFunc("GET /api/metro_line", func(w http.ResponseWriter, r *http.Request) {
		lineCode := r.URL.Query().Get("line")

		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/metro_line/%s", Metro_Base_URL, lineCode)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	mux.HandleFunc("GET /api/station_by_line", func(w http.ResponseWriter, r *http.Request) {
		lineCode := r.URL.Query().Get("line")

		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/station_by_line_linepage/%s", Metro_Base_URL, lineCode)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	mux.HandleFunc("GET /api/station_info", func(w http.ResponseWriter, r *http.Request) { //PROVIDE STATION CODE NOT ID
		stationCode := r.URL.Query().Get("code")

		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/station/%s", Metro_Base_URL, stationCode)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	mux.HandleFunc("GET /api/DMRC_route", func(w http.ResponseWriter, r *http.Request) { //PROVIDE STATION CODE NOT ID
		fromCode := r.URL.Query().Get("fromCode")
		toCode := r.URL.Query().Get("toCode")
		routeType := r.URL.Query().Get("routeType")
		departureTime := r.URL.Query().Get("time") //FIX TIME PARSING

		// timeFormatted := departureTime.Format("2006-01-02T15:04:05.000")

		// BaseURL, fromCode, toCode, routeType, timeFormatted

		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/station_route/%s/%s/%s/%s", BaseURL, fromCode, toCode, routeType, departureTime)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
		if err != nil {
			writeJSON(w, respCode, map[string]string{"error": "Upstream request encountered error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respCode)
		w.Write(body)
	})

	mux.HandleFunc("GET /api/DMRC_route_fare", func(w http.ResponseWriter, r *http.Request) { //PROVIDE STATION CODE NOT ID
		fromCode := r.URL.Query().Get("fromCode")
		toCode := r.URL.Query().Get("toCode")
		routeType := r.URL.Query().Get("routeType")

		if Metro_Base_URL == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong with server ENV"})
			return
		}

		upstreamURL := fmt.Sprintf("%s/new_fare_with_route/%s/%s/%s/", BaseURL, fromCode, toCode, routeType)
		body, respCode, err := MetroUpstreamProxy(upstreamURL)
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

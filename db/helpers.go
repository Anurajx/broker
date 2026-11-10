package main

import (
	"database/sql"
	"encoding/json"
)

func InsertStation(db *sql.DB, s MetroStation) error {
	// Marshal nested struct slices into raw JSON strings
	linesJSON, _ := json.Marshal(s.MetroLines)
	prevNextJSON, _ := json.Marshal(s.PrevNextStations)
	facilitiesJSON, _ := json.Marshal(s.StationFacility)
	gatesJSON, _ := json.Marshal(s.Gates)
	liftsJSON, _ := json.Marshal(s.Lifts)
	platformsJSON, _ := json.Marshal(s.Platforms)
	stnFacilitiesJSON, _ := json.Marshal(s.StationsFacilities)
	parkingsJSON, _ := json.Marshal(s.Parkings)

	query := `
	INSERT INTO metro_stations (
		id, station_code, station_name, station_description, station_type,
		interchange, latitude, longitude, mobile, landline, opening_time, closing_time,
		metro_lines, prev_next_stations, station_facility, gates, lifts, platforms, stations_facilities, parkings
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
	station_code        = excluded.station_code,
    station_name        = excluded.station_name,
    station_description = excluded.station_description,
    station_type        = excluded.station_type,
    interchange         = excluded.interchange,
    latitude            = excluded.latitude,
    longitude           = excluded.longitude,
    mobile              = excluded.mobile,
    landline            = excluded.landline,
    opening_time        = excluded.opening_time,
    closing_time        = excluded.closing_time,
    metro_lines         = excluded.metro_lines,
    prev_next_stations  = excluded.prev_next_stations,
    station_facility    = excluded.station_facility,
    gates               = excluded.gates,
    lifts               = excluded.lifts,
    platforms           = excluded.platforms,
    stations_facilities = excluded.stations_facilities,
    parkings            = excluded.parkings;
	`

	_, err := db.Exec(query,
		s.ID, s.StationCode, s.StationName, s.StationDescription, s.StationType,
		s.Interchange, s.Latitude, s.Longitude, s.Mobile, s.Landline, s.OpeningTime, s.ClosingTime,
		string(linesJSON), string(prevNextJSON), string(facilitiesJSON),
		string(gatesJSON), string(liftsJSON), string(platformsJSON),
		string(stnFacilitiesJSON), string(parkingsJSON),
	)
	return err
}

func GetStationByID(db *sql.DB, id int) (*MetroStation, error) {
	query := `
	SELECT id, station_code, station_name, station_description, station_type,
	       interchange, latitude, longitude, mobile, landline, opening_time, closing_time,
	       metro_lines, prev_next_stations, station_facility, gates, lifts, platforms, stations_facilities, parkings
	FROM metro_stations WHERE id = ?`

	var s MetroStation
	var linesRaw, prevNextRaw, facRaw, gatesRaw, liftsRaw, platformsRaw, stnFacRaw, parkingsRaw string

	err := db.QueryRow(query, id).Scan(
		&s.ID, &s.StationCode, &s.StationName, &s.StationDescription, &s.StationType,
		&s.Interchange, &s.Latitude, &s.Longitude, &s.Mobile, &s.Landline, &s.OpeningTime, &s.ClosingTime,
		&linesRaw, &prevNextRaw, &facRaw, &gatesRaw, &liftsRaw, &platformsRaw, &stnFacRaw, &parkingsRaw,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON strings directly back into your struct fields
	json.Unmarshal([]byte(linesRaw), &s.MetroLines)
	json.Unmarshal([]byte(prevNextRaw), &s.PrevNextStations)
	json.Unmarshal([]byte(facRaw), &s.StationFacility)
	json.Unmarshal([]byte(gatesRaw), &s.Gates)
	json.Unmarshal([]byte(liftsRaw), &s.Lifts)
	json.Unmarshal([]byte(platformsRaw), &s.Platforms)
	json.Unmarshal([]byte(stnFacRaw), &s.StationsFacilities)
	json.Unmarshal([]byte(parkingsRaw), &s.Parkings)

	return &s, nil
}

func InsertMetroLine(db *sql.DB, line MetroLine) error {
	startStationJSON, _ := json.Marshal(line.StartStation)
	endStationJSON, _ := json.Marshal(line.EndStation)

	query := `
	INSERT INTO metro_lines (
		id, name, line_code, line_color_verbose, line_color_family,
		primary_color_code, secondary_color_code, station_count,
		start_station, end_station, junction_count, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name                 = excluded.name,
		line_code            = excluded.line_code,
		line_color_verbose   = excluded.line_color_verbose,
		line_color_family    = excluded.line_color_family,
		primary_color_code   = excluded.primary_color_code,
		secondary_color_code = excluded.secondary_color_code,
		station_count        = excluded.station_count,
		start_station        = excluded.start_station,
		end_station          = excluded.end_station,
		junction_count       = excluded.junction_count,
		status               = excluded.status;
	`

	_, err := db.Exec(query,
		line.ID, line.Name, line.LineCode, line.LineColorVerbose, line.LineColorFamily,
		line.PrimaryColorCode, line.SecondaryColorCode, line.StationCount,
		string(startStationJSON), string(endStationJSON), line.JunctionCount, line.Status,
	)
	return err
}

func GetMetroLineByID(db *sql.DB, id int) (*MetroLine, error) {
	query := `
	SELECT id, name, line_code, line_color_verbose, line_color_family,
	       primary_color_code, secondary_color_code, station_count,
	       start_station, end_station, junction_count, status
	FROM metro_lines WHERE id = ?`

	var line MetroLine
	var startStationRaw, endStationRaw string

	err := db.QueryRow(query, id).Scan(
		&line.ID, &line.Name, &line.LineCode, &line.LineColorVerbose, &line.LineColorFamily,
		&line.PrimaryColorCode, &line.SecondaryColorCode, &line.StationCount,
		&startStationRaw, &endStationRaw, &line.JunctionCount, &line.Status,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON strings directly back into nested struct fields
	json.Unmarshal([]byte(startStationRaw), &line.StartStation)
	json.Unmarshal([]byte(endStationRaw), &line.EndStation)

	return &line, nil
}

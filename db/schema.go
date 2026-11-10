package main

import (
	"time"
)

// METRO ---
type MetroStation struct {
	ID                 int     `json:"id"`
	StationCode        string  `json:"station_code"`
	StationName        string  `json:"station_name"`
	StationDescription string  `json:"station_description"`
	StationType        string  `json:"station_type"`
	Interchange        bool    `json:"interchange"`
	Latitude           string  `json:"latitude"`
	Longitude          string  `json:"longitude"`
	Mobile             string  `json:"mobile"`
	Landline           *string `json:"landline"` // Pointer handles null values
	OpeningTime        string  `json:"opening_time"`
	ClosingTime        string  `json:"closing_time"`

	// Metro Lines serving this station
	MetroLines []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		LineColor    string `json:"line_color"`
		LineCode     string `json:"line_code"`
		StartStation string `json:"start_station"`
		EndStation   string `json:"end_station"`
		Status       string `json:"status"`
	} `json:"metro_lines"`

	// Dynamic map for "Line 3", "Line 6", etc.
	PrevNextStations []map[string][]struct {
		LineID      int    `json:"line_id"`
		LineColor   string `json:"line_color"`
		PrevStation struct {
			ID          int    `json:"id"`
			StationName string `json:"station_name"`
			StationCode string `json:"station_code"`
		} `json:"prev_station"`
		NextStation struct {
			ID          int    `json:"id"`
			StationName string `json:"station_name"`
			StationCode string `json:"station_code"`
		} `json:"next_station"`
	} `json:"prev_next_stations"`

	StationFacility []struct {
		Name string `json:"name"`
	} `json:"station_facility"`

	Gates []struct {
		GateName        string `json:"gate_name"`
		GateCode        string `json:"gate_code"`
		Location        string `json:"location"`
		DivyangFriendly bool   `json:"divyang_friendly"`
		Status          string `json:"status"`
	} `json:"gates"`

	Lifts []struct {
		LiftType               string `json:"lift_type"`
		Name                   string `json:"name"`
		DescriptionLocation    string `json:"description_location"`
		Code                   string `json:"code"`
		AvailableOutsideInside string `json:"available_outside_inside"`
		DivyangFriendly        bool   `json:"divyang_friendly"`
		Status                 bool   `json:"status"`
	} `json:"lifts"`

	Platforms []struct {
		PlatformName string `json:"platform_name"`
		PlatformCode string `json:"platform_code"`
		TrainTowards struct {
			ID          int    `json:"id"`
			StationName string `json:"station_name"`
			StationCode string `json:"station_code"`
		} `json:"train_towards"`
	} `json:"platforms"`

	StationsFacilities []struct {
		Kind       string `json:"kind"`
		DetailList []struct {
			FacilityName        string `json:"facility_name"`
			Purpose             string `json:"purpose"`
			LocationDescription string `json:"location_description"`
		} `json:"detail_list"`
	} `json:"stations_facilities"`

	Parkings []struct {
		Provider           string `json:"provider"`
		CapacityCar        int    `json:"capacity_car"`
		CapacityMotorcycle int    `json:"capacity_motorcycle"`
		CapacityCycle      int    `json:"capacity_cycle"`
		ParkingCode        string `json:"parking_code"`
		Location           string `json:"location"`
	} `json:"parkings"`
}

type MetroLine struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	LineColorVerbose   string `json:"line_color"` //Extensions of line can have diff names
	LineCode           string `json:"line_code"`
	LineColorFamily    string `json:"class_primary"` //Color family remains same across same line
	PrimaryColorCode   string `json:"primary_color_code"`
	SecondaryColorCode string `json:"secondary_color_code"`
	StationCount       int    `json:"station_count"`
	StartStation       struct {
		ID          int    `json:"id"`
		StationName string `json:"station_name"`
		StationCode string `json:"station_code"`
	} `json:"first_station"`
	EndStation struct {
		ID          int    `json:"id"`
		StationName string `json:"station_name"`
		StationCode string `json:"station_code"`
	} `json:"last_station"`
	JunctionCount int    `json:"junction"`
	Status        string `json:"status"`
}

type LineStatus struct { //FROM METRO LINE ENDPOINT
	StationUpdate string        `json:"stations_updates"`
	LineUpdate    []interface{} `json:"line_updates"`
}

type MetroRoute struct {
	totalStation int
	from         string
	to           string
	fromStatus   bool
	toStatus     bool
	timeTaken    time.Time
	fare         int
	route        []MetroPath
	message      string
	// routeType    string // Shortest OR LeastInterchange
}

type MetroPath struct { //USE THIS FOR -> Metro lines and Interchanges
	line               string
	interchange        bool
	line_no            int
	path_time          time.Time
	startStationName   string
	endStationName     string
	towardsStationName string
	intermediate       []string
	platform           int
}

// BUS ---

type BusStop struct {
	ID        int
	latitude  float64
	longitude float64
	name      string
	schedule  []BusStopETA
}

type BusStopETA struct {
	routeShortName string
	routeLongName  string
	routeID        int
	startStop      string
	endStop        string
	busArrival     struct {
		AC          bool
		agency      string
		color       string
		ETA         int
		numberPlate string
	}
}

type TSroute struct {
	route []TSpath
	time  int
}

type TSpath struct {
	pathType string   //Walk, Bus, Metro
	routes   []string //List of route numbers for pathType
	agency   string
	towards  string
	color    string
	stops    []string
}

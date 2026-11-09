package main

import (
	"time"
)

// METRO ---

type metroStation struct {
	id                 string
	stationCode        string
	stationName        string
	stationDescription string
	elevated           bool
	interchange        bool
	latitude           float64
	longitude          float64
	mobile             int
	landline           int
	stationFacility    map[string]bool
	gates              map[int]string
	lift               map[int]string
	escalator          map[int]string
	platform           map[int]string //Platform no. -> Towards stationCode
	parking            map[string]int
	openingTime        time.Time
	closingTime        time.Time
	status             string
	lines              map[int]metroLine //LINE ID -> metroLine
	prevStation        map[int]string    //LINE ID -> stationCode

}

type metroLine struct {
	id                 int
	name               string
	lineColorVerbose   string //Extensions of line can have diff names
	lineCode           string
	lineColorFamily    string //Color family remains same across same line
	primaryColorCode   string
	secondaryColorCode string
	stationCount       int
	startStation       string //Store station Code
	endStation         string
	junctionCount      int
	status             string
	stationUpdate      string
	lineUpdate         string
}

type metroRoute struct {
	totalStation int
	from         string
	to           string
	fromStatus   bool
	toStatus     bool
	timeTaken    time.Time
	fare         int
	route        []metroPath
	message      string
	// routeType    string // Shortest OR LeastInterchange
}

type metroPath struct { //USE THIS FOR -> Metro lines and Interchanges
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

type busStop struct {
	id        int
	latitude  float64
	longitude float64
	name      string
	schedule  []busStopETA
}

type busStopETA struct {
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

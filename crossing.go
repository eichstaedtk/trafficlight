/*
 * crossing.go
 *
 * A program to represent a traffic crossing with 4 traffic lights as goroutines.
 *
 * author: Konrad Eichstädt
 */
package main

import (
	"log"
	"time"
)

var n, s, e, w TrafficLight
var coloursNorthsouth = make(chan TrafficLightColours)
var coloursEastwest = make(chan TrafficLightColours)

var axisSwitchChannelNorthEast = make(chan AxisSwitch)
var axisSwitchChannelSouthWest = make(chan AxisSwitch)

var quitChannel = make(chan string)

type Message struct {
	activeLight TrafficLightColours
	activeAxis  DirectionAxis
}

var syncChanel = make(chan Message)

/*
* Main Routine to start the application
 */

func main() {

	log.Printf("*** Activation of Crossing of Traffic Lights ***")

	startTrafficLights()

	//waiting time Duration of Active Crossing

	time.Sleep(10 * time.Millisecond)

	quitChannel <- "stop"

	log.Printf("*** Finishing of  Crossing of Traffic Light ***") // wait for a quit signal
}

/*
* Start function for construction and activation of all traffic lights
 */

func startTrafficLights() {

	// Construction of Crossing Traffic Lights

	n := n.NewTrafficLight(North, Red, coloursNorthsouth, axisSwitchChannelNorthEast, quitChannel)
	s := s.NewTrafficLight(South, Red, coloursNorthsouth, axisSwitchChannelSouthWest, quitChannel)
	e := e.NewTrafficLight(East, Red, coloursEastwest, axisSwitchChannelNorthEast, quitChannel)
	w := w.NewTrafficLight(West, Red, coloursEastwest, axisSwitchChannelSouthWest, quitChannel)

	// Starting Traffic Lights GotRoutines

	go n.Start()
	go s.Start()
	go e.Start()
	go w.Start()

	//north-south axis starts, east-west axis has to wait to take over control - Sending Signal to Axis Channel for AxisSwitch
}

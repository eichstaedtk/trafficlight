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
	"sync"
	"time"
)

var n,s,e,w TrafficLight
var coloursNorthsouth = make(chan TrafficLightColours)
var coloursEastwest = make(chan TrafficLightColours)

var axisChannel = make(chan AxisSwitch)
var quitChannel = make(chan string)

var wgNorthSouth sync.WaitGroup
var wgNEastWest sync.WaitGroup

/*
* Main Routine to start the application
*/

func main() {

	log.Printf("*** Activation of Crossing of Traffic Lights ***")

	startTrafficLights()

	//waiting time Duration of Active Crossing

	time.Sleep(10 * time.Millisecond)

	//sending stop signal to quitChannel

	quitChannel <- "stop"

	log.Printf("*** Finishing of  Crossing of Traffic Light ***") // wait for a quit signal
}

/*
* Start function for construction and activation of all traffic lights
*/

func startTrafficLights() {

	// Construction of Crossing Traffic Lights

	n := n.NewTrafficLight(North,Red, coloursNorthsouth,axisChannel,quitChannel,&wgNorthSouth)
	s := s.NewTrafficLight(South,Red, coloursNorthsouth,axisChannel,quitChannel,&wgNorthSouth)
	e := e.NewTrafficLight(East,Red, coloursEastwest,axisChannel,quitChannel,&wgNEastWest)
	w := w.NewTrafficLight(West,Red, coloursEastwest,axisChannel,quitChannel,&wgNEastWest)

	// Starting Traffic Lights GotRoutines

	go n.Start()
	go s.Start()
	go e.Start()
	go w.Start()

	//north-south axis starts, east-west axis has to wait to take over control - Sending Signal to Axis Channel for AxisSwitch

	axisChannel <- AxisSwitch{ActiveAxis: North_South}
}
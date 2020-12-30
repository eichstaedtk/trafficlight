/*
 * trafficlight.go
 *
 * A program to represent a traffic light.
 *
 * author: Konrad Eichstädt
 */
package main

import (
	"log"
	"sync"
)

type TrafficLightColours int
type Direction int
type DirectionAxis int

const (
	Red TrafficLightColours = iota
	Green
	Yellow
)
const (
	North Direction = iota
	East
	South
	West
)
const (
	East_West DirectionAxis = iota
	North_South
)

type AxisSwitch struct{ ActiveAxis DirectionAxis }
type TrafficLight struct {
	Direction   Direction
	Colour      TrafficLightColours
	Colours     chan TrafficLightColours
	Axis        chan AxisSwitch
	activeAxis  DirectionAxis
	QuitChannel chan string
	WaitGroup   *sync.WaitGroup
}

func (t TrafficLight) NewTrafficLight(direction Direction, colour TrafficLightColours, colours chan TrafficLightColours, axis chan AxisSwitch, quitChannel chan string, waitGroup *sync.WaitGroup) TrafficLight {
	t.Direction = direction
	t.Colour = colour
	t.Colours = colours
	t.Axis = axis
	t.QuitChannel = quitChannel
	t.WaitGroup = waitGroup
	return t
}

func (t TrafficLight) Start() {

loop:

	for {

		select {

		case switchAxis := <-t.Axis:

			t.activeAxis = switchAxis.ActiveAxis

			//Synchronization of one Axis check if the traffic light is on active axis

			if switchAxis.ActiveAxis == axis(t.Direction) {

				t.Show()
				t.WaitGroup.Add(1)

				//Sending Active Colour and wait for synchronisation of one axis

				t.Colours <- t.Colour
				t.WaitGroup.Wait()

				//Green
				t.switchLight()
				//Yellow
				t.switchLight()
				//Red
				t.switchLight()

				//Sending Signal for Switching Crossing Control
				if axis(t.Direction) == North_South {
					t.Axis <- AxisSwitch{ActiveAxis: East_West}
					log.Printf("Switch to EAST_WEST")
				} else {
					t.Axis <- AxisSwitch{ActiveAxis: North_South}
					log.Printf("Switch to NORTH_SOUTH")
				}

			} else {
				t.Axis <- switchAxis
			}

		case color := <-t.Colours:
			{

				t.Colour = color
				t.Show()
				t.WaitGroup.Done()

			}

		case <-t.QuitChannel:
			break loop
		}

	}

}

/**
* Function to switch the light of traffic lights
 */

func (t *TrafficLight) switchLight() {
	t.Colour = next(t.Colour)
	t.Show()
	t.WaitGroup.Add(1)
	t.Colours <- t.Colour
	t.WaitGroup.Wait()
}

/**
* Function to show the traffic lights
 */

func (t TrafficLight) Show() {
	log.Printf("Traffic Light Direction %s and Colour %s", t.Direction.String(), t.Colour.String())
}

/**
 * Function to switch the traffic lights
 */

func next(colour TrafficLightColours) TrafficLightColours {

	if colour == Red {
		return Green
	}
	if colour == Green {
		return Yellow
	}
	if colour == Yellow {
		return Red
	}

	return -1
}

/**
 * Function to convert direction into a direction axis
 */

func axis(direction Direction) DirectionAxis {

	if direction == North || direction == South {
		return North_South
	} else if direction == East || direction == West {
		return East_West
	}

	return -1
}

func (c TrafficLightColours) String() string {
	return [...]string{"Red", "Green", "Yellow"}[c]
}

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

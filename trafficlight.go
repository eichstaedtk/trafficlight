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
	Direction         Direction
	ActiveAxis        DirectionAxis
	Colour            TrafficLightColours
	AxisSwitchChannel chan AxisSwitch
	ColoursChannel    chan TrafficLightColours
	QuitChannel       chan string
}

func (t TrafficLight) NewTrafficLight(direction Direction, colour TrafficLightColours, coloursChannel chan TrafficLightColours, axisSwitchChannel chan AxisSwitch, quitChannel chan string) TrafficLight {
	t.Direction = direction
	t.Colour = colour
	t.AxisSwitchChannel = axisSwitchChannel
	t.ColoursChannel = coloursChannel
	t.QuitChannel = quitChannel
	return t
}

func (t TrafficLight) Start() {

loop:

	for {

		t.Show()

		select {

		case t.ColoursChannel <- t.Colour:

		case colour := <-t.ColoursChannel:
			t.Colour = colour

		case <-t.QuitChannel:
			goto loop

		}

		if t.Colour == Red {
			if t.ActiveAxis == North_South {
				t.ActiveAxis = East_West
			} else {
				t.ActiveAxis = North_South
			}

			select {
			case sc := <-t.AxisSwitchChannel:
				t.ActiveAxis = sc.ActiveAxis

			case t.AxisSwitchChannel <- AxisSwitch{t.ActiveAxis}:
			}
		}

		if axis(t.Direction) == t.ActiveAxis {
			t.Colour = next(t.Colour)
		}
	}
}

/**
* Function to show the traffic lights
 */

func (t TrafficLight) Show() {
	log.Printf("Traffic Light Direction %s and Colour %s  ActiveAxis %s", t.Direction.String(), t.Colour.String(), t.ActiveAxis)
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

func (d DirectionAxis) String() string {
	return [...]string{"EAST_WEST", "NORTH_SOUTH"}[d]
}

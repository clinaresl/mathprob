// -*- coding: utf-8 -*-
// rectangle.go
//
// Description: Definition of a rectangle using the lower-left and upper-right
// coordinates with an optional color
//
// -----------------------------------------------------------------------------
//
// Started on <dom 23-05-2021 19:22:32.823217719 (1621790552)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// This package provides a number of reusable components that can be used for
// creating TikZ drawings
package components

import (
	"errors"
	"fmt"
	"log"

	"github.com/clinaresl/mathprob/helpers"
)

// constants
// ----------------------------------------------------------------------------

// TikZ code to generate a rectangle: it just simply create a rectangle with the
// specified color from the lower-left to the upper-right coordinates
const tikzRectangle = `\draw [{{.GetColor}}] ({{.GetCoordinate0}}) rectangle ({{.GetCoordinate1}});`

// types
// ----------------------------------------------------------------------------

// A rectangle requires two coordinates (either explicit or formulas or a
// combination of both) given to specify the lower-left (coordinate0) and
// upper-right (coordinate1) corners. Additionally, a color can be given
type Rectangle struct {
	coord0, coord1 Coordinate
	color          string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a rectangle given two coordinates. Note that the
// color is specified through a dedicated service
func NewRectangle(coord0, coord1 Coordinate) Rectangle {
	return Rectangle{coord0: coord0,
		coord1: coord1,
	}
}

// return a valid specification of a rectangle with no error if all the keys
// given in dict are correct for defining a rectangle. Otherwise, return an
// error. If an error is returned, the contents of the rectangle are
// undetermined
//
// A dictionary is correct if and only if it correctly defines a rectangle,
// i.e., the lower-left and upper-right corners should be correctly specified as
// Coordinates. These are the only mandatory arguments. In addition, it is also
// possible to specify a color as a string
func VerifyRectangleDict(dict map[string]interface{}) (Rectangle, error) {

	// first of all, ensure that all mandatory parameters are given and that
	// they are of the correct type. Create slices for both mandatory and all
	// arguments
	all := []string{"coord0", "coord1", "color"}
	mandatory := []string{"coord0", "coord1"}

	// verify that all mandatory arguments are given in the dictionary
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then immediately raise
		// an error
		if _, ok := dict[key]; !ok {
			return Rectangle{}, fmt.Errorf("Mandatory key '%v' for defining a rectangle not found", key)
		}
	}

	// now ensure that the mandatory parameters are of the right type
	var err error
	var coord0, coord1 Coordinate
	if coord0, err = VerifyCoordinateDict(dict); err != nil {

		return Rectangle{}, fmt.Errorf("Processing coord0 of a rectangle raised the error: %v", err)
	}
	if coord1, err = VerifyCoordinateDict(dict); err != nil {
		return Rectangle{}, fmt.Errorf("Processing coord1 of a rectangle raised the error: %v", err)
	}

	// now, perform the same operation with the optional parameters
	var color string
	if _, ok := dict["color"]; ok {
		if _, ok := dict["color"].(string); !ok {
			return Rectangle{}, errors.New("The color of a rectangle should be given as a string")
		}
		color = dict["color"].(string)
	}

	// in case any other arguments were given, but they are not acknoweldged,
	// issue a warning
	for key, _ := range dict {
		if !helpers.Find(key, all) {
			log.Printf("The parameter '%v' is not acknowledged for creating a rectangle and it will be ignored")
		}
	}

	// At this point, the dictionary is correct, return a valid box
	return Rectangle{coord0: coord0,
		coord1: coord1,
		color:  color,
	}, nil
}

// methods
// ----------------------------------------------------------------------------

// Set the color of a rectangle
func (rect *Rectangle) SetColor(color string) {
	rect.color = color
}

// Local Variables:
// mode:go
// fill-column:80
// End:

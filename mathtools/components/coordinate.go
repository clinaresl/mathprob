/*
  coordinate.go

  Description: Definition of coordinates as reusable components to be used in
               TikZ drawings

  -----------------------------------------------------------------------------

  Started on  <Mon Jul 10 09:31:10 2017 >
  Last update <>
  -----------------------------------------------------------------------------
  Made by Carlos Linares López
  Login <carlos.linares@uc3m.es>
*/

package components

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"
)

// constants
// ----------------------------------------------------------------------------

// TikZ code to generate a coordinate
const tikzCoordinate = `\coordinate {{.GetLabel}} at {{.GetPosition}};`

// types
// ----------------------------------------------------------------------------

// Any element that can be positioned is by definition a position
type Position interface {
	Position() string
}

// The very first component is a point, which is given as a pair of
// floating-point numbers
type Point struct {
	X, Y float64
}

// Also formulas can be created. They are defined as strings which shall be
// given as valid TikZ expressions
type Formula string

// A coordinate is created by providing a label to either a point or a formula
// and, in general, to any component that is positionable
type Coordinate struct {
	Position
	label string
}

// functions
// ----------------------------------------------------------------------------

// Create a new coordinate given a label and a positionable element
func NewCoordinate(position Position, label string) Coordinate {
	return Coordinate{Position: position, label: label}
}

// return a valid Point and no error if the keywords "x" and "y" are given in
// the dictionary. Otherwise, an error is returned. If an error is returned the
// contents of the Point are undetermined.
//
// The values of "x" and "y" must be floating-point numbers. If the keywords
// exist but the type assertion fails it returns false the same
func verifyPointDict(dict map[string]interface{}) (Point, error) {

	// traverse the entire dictionary and get the values of "x" and "y" in case
	// they are present
	var ok bool
	var fvalue float64
	coords := make(map[string]float64)
	for key, value := range dict {
		if key == "x" {

			// if the value of x is not given as a floating-point number
			// immediately return false
			if fvalue, ok = value.(float64); !ok {
				return Point{}, errors.New("The x coordinate of a Point must be given as a floating-point number")
			}
			coords[key] = fvalue
		}
		if key == "y" {

			// if the value of y is not given as a floating-point number
			// immediately return false
			if fvalue, ok = value.(float64); !ok {
				return Point{}, errors.New("The y coordinate of a Point must be given as a floating-point number")
			}
			coords[key] = fvalue
		}
	}

	// if either x or y is missing, then return false
	if len(coords) < 2 {
		return Point{}, errors.New("Both 'x' and 'y' coordinates must be given for defining a Point")
	}

	// at this point, two coordinates were correctly provided so that return the
	// corresponding Point and no error
	return Point{X: coords["x"], Y: coords["y"]}, nil
}

// return a valid Formula and no error if the keyword "formula" is given in the
// dictionary. Otherwise, an error is returned. If an error is returned the
// contents of the Formula are undetermined.
//
// The value of the keyword "formula" must be a string. If the keyword exist but
// the type assertion fails, or an empty string is given, it returns false the
// same
func verifyFormulaDict(dict map[string]interface{}) (Formula, error) {

	// traverse the entire dictionary and get the value of "formula" in case
	// it is present
	var ok bool
	svalue := ""
	for key, value := range dict {
		if key == "formula" {

			// if the value of x is not given as a string immediately return
			// false
			if svalue, ok = value.(string); !ok {
				return Formula(""), errors.New("Formulas have to be given as strings")
			}
		}
	}

	// if no formula was given, or it was empty return false
	if svalue == "" {
		return Formula(""), errors.New("Either a formula was not given or it is the empty string")
	}

	// at this point, a valid formula has been specified
	return Formula(svalue), nil
}

// return a valid coordinate and no error if all the keys given in dict are
// correct for defining a coordinate. Otherwise, return an error. If an error is
// returned the contents of the Coordinate are undetermined
//
// A dictionary is correct if and only if all the mandatory arguments have been
// given, and all parameters given are compatible.
//
// A coordinate can be specified either with a Formula (using the keyword
// "formula") which should consist of a string, or as a pair of floating-point
// numbers x and y which create a Point, using the keywords "x" and "y"
func VerifyCoordinateDict(dict map[string]interface{}) (Coordinate, error) {

	// the mandatory keys are given next
	mandatory := []string{"label"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then raise an error and
		// exit
		if _, ok := dict[key]; !ok {
			return Coordinate{}, fmt.Errorf("Mandatory key '%v' for defining a coordinate not found", key)
		}
	}

	// secondly, verify that a point and a formula haven not been simultaneously
	// given
	point, errp := verifyPointDict(dict)
	formula, errf := verifyFormulaDict(dict)
	if errp == nil && errf == nil {
		return Coordinate{}, errors.New("Either a 'position' or 'formula' have to be given, but not both")
	}

	// likewise, if neither a point nor a formula were given, then return an
	// error
	if errp != nil && errf != nil {
		return Coordinate{}, errors.New("Either a 'position' or a 'formula' have to be given")
	}

	// otherwise, the dictionary is correct and a coordinate has been correctly
	// speciied. First, if the coordinate was given relative to a point
	if errp == nil {
		return Coordinate{Position: point, label: dict["label"].(string)}, nil
	}

	// if not, then return the coordinate defined with a formula
	return Coordinate{Position: formula, label: dict["label"].(string)}, nil
}

// methods
// ----------------------------------------------------------------------------

// -- Point

// Points are of course positionable and, as such, they return a string that
// represents their location as a valid TikZ representation
func (p Point) Position() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

// -- Formula

// Formulas are also positionable as their resolution results in a unique point
// and, as such, they return a string that can be used to compute their location
// as a valid TikZ representation
func (f Formula) Position() string {
	return fmt.Sprintf("($%v$)", string(f))
}

// -- Coordinate

// Return the label of this coordinate
func (c Coordinate) GetLabel() string {
	return fmt.Sprintf("(%s)", c.label)
}

// Return the position of this coordinate.
func (c Coordinate) GetPosition() string {
	return c.Position.Position()
}

// return a TikZ representation of a coordinate
func (c Coordinate) String() string {

	// create a template with the TikZ code for showing a coordinate
	tpl, err := template.New("coordinate").Parse(tikzCoordinate)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, c); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

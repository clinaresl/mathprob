/*
  components.go
  Description: Definition of different reusable components to be used in TikZ
               drawings
  -----------------------------------------------------------------------------

  Started on  <Mon Jul 10 09:31:10 2017 >
  Last update <>
  -----------------------------------------------------------------------------
  Made by Carlos Linares López
  Login <carlos.linares@uc3m.es>
*/

// This package provides a number of reusable components that can be used for
// creating exercises automatically.
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

// Unless it is truly obvious or simplistic, TikZ code is generated via text
// templates

// TikZ code to generate a coordinate
const tikzCoordinate = `\coordinate {{.GetLabel}} at {{.GetPosition}};`

// TikZ code to generate text
const tikzText = `\draw {{.GetLabel}} node {{.GetText}};`

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

// Text is written centered in a given coordinate identified by a label and can
// hold any string (including other LaTeX commands that can affect the
// look&feel)
type Text struct {
	label string
	text  string
}

// functions
// ----------------------------------------------------------------------------

// return true and a valid Point if the keywords "x" and "y" are given in the
// dictionary, and false otherwise. The values of "x" and "y" must be
// floating-point numbers. If the keywords exist but the type assertion fails it
// returns false the same
func verifyPointDict(dict map[string]interface{}) (bool, Point) {

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
				return false, Point{0, 0}
			}
			coords[key] = fvalue
		}
		if key == "y" {

			// if the value of y is not given as a floating-point number
			// immediately return false
			if fvalue, ok = value.(float64); !ok {
				return false, Point{0, 0}
			}
			coords[key] = fvalue
		}
	}

	// if either x or y is missing, then return false
	if len(coords) < 2 {
		return false, Point{0, 0}
	}

	// at this point, two coordinates were correctly provided so that return the
	// corresponding Point
	return true, Point{X: coords["x"], Y: coords["y"]}
}

// return true and a valid Formula if the keyword "formula" is given in the
// dictionary, and false otherwise. The values of "string" must be a string. If
// the keyword exist but the type assertion fails, or an empty string is given,
// it returns false the same
func verifyFormulaDict(dict map[string]interface{}) (bool, Formula) {

	// traverse the entire dictionary and get the values of "x" and "y" in case
	// they are present
	var ok bool
	svalue := ""
	for key, value := range dict {
		if key == "formula" {

			// if the value of x is not given as a string immediately return
			// false
			if svalue, ok = value.(string); !ok {
				return false, Formula("")
			}
		}
	}

	// if no formula was given, or it was empty return false
	if svalue == "" {
		return false, Formula("")
	}

	// at this point, a formula was given, so return it
	return true, Formula(svalue)
}

// return true if all the keys given in dict are correct for defining a
// coordinate.
//
// A dictionary is correct if and only if all the mandatory arguments have been
// given, and all parameters given are compatible. If not, false and an error is
// returned.
//
// A coordinate can be specified either with a Formula which should consist of a
// string, or as a pair of floating-point numbers x and y which create a Point
func VerifyCoordinateDict(dict map[string]interface{}) (bool, error) {

	// the mandatory keys are given next
	mandatory := []string{"label"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then raise an error and
		// exit
		if _, ok := dict[key]; !ok {
			return false, fmt.Errorf("Mandatory key '%v' for defining a coordinate not found", key)
		}
	}

	// secondly, verify that a point and a formula haven not been simultaneously
	// given
	okp, _ := verifyPointDict(dict)
	okf, _ := verifyFormulaDict(dict)
	if okp && okf {
		return false, errors.New("Either a 'position' or 'formula' have to be given, but not both")
	}

	// likewise, if neither a point nor a formula were given, then return an error
	if !okp && !okf {
		return false, errors.New("Either a 'position' or a 'formula' have to be given")
	}

	// otherwise, the dictionary is correct
	return true, nil
}

// return true if all the keys given in dict are correct for defining a
// text box
//
// A dictionary is correct if and only if all the mandatory arguments have been
// given. If not, false and an error is returned.
//
// A text box is specified providing the label of a coordinate and also the text
// to show with the keywords "label" and "text"
func VerifyTextDict(dict map[string]interface{}) (bool, error) {

	// the mandatory keys are given next
	mandatory := []string{"label", "text"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then raise an error and
		// exit
		if _, ok := dict[key]; !ok {
			return false, fmt.Errorf("Mandatory key '%v' for defining a text box not found", key)
		}
	}

	// make also sure that parameters are given with the right type
	if _, ok := dict["label"].(string); !ok {
		return false, errors.New("the label of a text box should be given as a string")
	}
	if _, ok := dict["text"].(string); !ok {
		return false, errors.New("the text to show in a text box should be given as a string")
	}

	// otherwise, the dictionary is correct
	return true, nil
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

// Create a new coordinate given a label and a positionable element
func NewCoordinate(position Position, label string) Coordinate {
	return Coordinate{Position: position, label: label}
}

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

// -- Text

// Create a new instance of a text box given a coordinate (identified by its
// label) and the text to show
func NewText(label, text string) Text {
	return Text{label: label, text: text}
}

// Return the label identifying the coordinate of this text box
func (t Text) GetLabel() string {
	return fmt.Sprintf("(%v)", t.label)
}

// Return the text to show of this text box
func (t Text) GetText() string {
	return fmt.Sprintf("{%v}", t.text)
}

// return a TikZ representation of a text box
func (t Text) String() string {

	// create a template with the TikZ code for showing a text box
	tpl, err := template.New("text").Parse(tikzText)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, t); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

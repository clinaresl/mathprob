// -*- coding: utf-8 -*-
// rectangle.go
//
// Description: Definition of a rectangle using the lower-left and upper-right
// coordinates with additional options
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
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"

	"github.com/clinaresl/mathprob/helpers"
)

// constants
// ----------------------------------------------------------------------------

// TikZ code to generate a rectangle: it just simply create a rectangle with the
// specified options from the lower-left to the upper-right references. Note
// that rectangles use references, either the names of coordinates (i.e.,
// labels) or formulas explicitly given
const tikzRectangle = `\draw [{{.GetOptions}}] ({{.GetReference0}}) rectangle ({{.GetReference1}});`

// TikZ code to generate a rectangle with coordinates: it consists of an
// ordinary rectangle created using labels but preceded of the coordinates used
// for creating them. Note that .GetPosition0 and .GetPosition1 return
// coordinates and thus they draw themselves with all the necessary information
// automatically
const tikzCoordinatedRectangle = `{{.GetPosition0}}
{{.GetPosition1}}
\draw [{{.GetOptions}}] {{.GetLabel0}} rectangle {{.GetLabel1}};`

// types
// ----------------------------------------------------------------------------

// Any rectangle has options which consist of a comma-separated list of options
// in a string
type BaseRectangle struct {
	options string
}

// A rectangle requires two references (either the name of labels or formulas
// explicitly given or a combination of both) given to specify the lower-left
// and upper-right coordinate. Additionally, an arbitrary number of options can
// be given as a comma-separated string
type Rectangle struct {
	ref0, ref1 string
	BaseRectangle
}

// A coordinated rectangle requires two coordinates to be explicitly created.
// These coordinates can be either given explicitly or using formulas
type CoordinatedRectangle struct {
	coord0, coord1 Coordinate
	BaseRectangle
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a rectangle given two references. Note that the
// options are specified through a dedicated service
func NewRectangle(ref0, ref1 string) Rectangle {
	return Rectangle{
		ref0: ref0,
		ref1: ref1,
	}
}

// Create a new instance of a coordinated rectangle given two coordinates. Note
// that the options are specified through a dedicated service
func NewCoordinatedRectangle(coord0, coord1 Coordinate) CoordinatedRectangle {
	return CoordinatedRectangle{
		coord0: coord0,
		coord1: coord1,
	}
}

// return a valid specification of a rectangle with no error if all the keys
// given in dict are correct for defining a rectangle. Otherwise, return an
// error. If an error is returned, the contents of the rectangle are
// undetermined
//
// A dictionary is correct if and only if it correctly defines a rectangle,
// i.e., the lower-left and upper-right references should be correctly specified
// as strings. These are the only mandatory arguments. In addition, it is also
// possible to specify arbitrary options as a string
func VerifyRectangleDict(dict map[string]interface{}) (Rectangle, error) {

	// first of all, ensure that all mandatory parameters are given and that
	// they are of the correct type. Create slices for both mandatory and all
	// arguments
	all := []string{"ref0", "ref1", "options"}
	mandatory := []string{"ref0", "ref1"}

	// verify that all mandatory arguments are given in the dictionary
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then immediately raise
		// an error
		if _, ok := dict[key]; !ok {
			return Rectangle{}, fmt.Errorf("Mandatory key '%v' for defining a rectangle not found", key)
		}
	}

	// now ensure that the mandatory parameters are of the right type
	var ok bool
	var ref0, ref1 string
	if ref0, ok = dict["ref0"].(string); !ok {
		return Rectangle{}, errors.New("The lower-left corner of a rectangle should be given as a string")
	}
	if ref1, ok = dict["ref1"].(string); !ok {
		return Rectangle{}, errors.New("The upper-right corner of a rectangle should be given as a string")
	}

	// now, perform the same operation with the optional parameters
	var options string
	if _, ok := dict["options"]; ok {
		if _, ok := dict["options"].(string); !ok {
			return Rectangle{}, errors.New("The options of a rectangle should be given as a string")
		}
		options = dict["options"].(string)
	}

	// in case any other arguments were given, but they are not acknoweldged,
	// issue a warning
	for key, _ := range dict {
		if !helpers.Find(key, all) {
			log.Printf("The parameter '%v' is not acknowledged for creating a rectangle and it will be ignored")
		}
	}

	// At this point, the dictionary is correct, return a valid box
	boptions := BaseRectangle{options: options}
	return Rectangle{ref0: ref0,
		ref1:          ref1,
		BaseRectangle: boptions,
	}, nil
}

// methods
// ----------------------------------------------------------------------------

// --BaseRectangle

// Set the options of a rectangle
func (rect *BaseRectangle) SetOptions(options string) {
	rect.options = options
}

// Get the options used
func (rect BaseRectangle) GetOptions() string {
	return rect.options
}

// --Rectangle

// Return the reference of the lower-left corner of the rectangle
func (rect Rectangle) GetReference0() string {

	return rect.ref0
}

// Return the reference of the upper-right corner of the rectangle
func (rect Rectangle) GetReference1() string {

	return rect.ref1
}

// Finally, rectangles are stringers and these are the means provided for
// automatically reusing this component
func (rect Rectangle) String() string {

	// create a template with the TikZ code for showing a rectangle
	tpl, err := template.New("rectangle").Parse(tikzRectangle)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, rect); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// --CoordinatedRectangle

// Return the label of the lower-left corner of the rectangle
func (rect CoordinatedRectangle) GetLabel0() string {

	return fmt.Sprintf("%v", rect.coord0.GetLabel())
}

// Return the label of the upper-right corner of the rectangle
func (rect CoordinatedRectangle) GetLabel1() string {

	return fmt.Sprintf("%v", rect.coord1.GetLabel())
}

// Return the reference of the lower-left corner of the rectangle
func (rect CoordinatedRectangle) GetPosition0() string {

	return fmt.Sprintf("%v", rect.coord0)
}

// Return the reference of the upper-right corner of the rectangle
func (rect CoordinatedRectangle) GetPosition1() string {

	return fmt.Sprintf("%v", rect.coord1)
}

// Finally, rectangles are stringers and these are the means provided for
// automatically reusing this component
func (rect CoordinatedRectangle) String() string {

	// create a template with the TikZ code for showing a rectangle
	tpl, err := template.New("rectangle").Parse(tikzCoordinatedRectangle)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, rect); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// Local Variables:
// mode:go
// fill-column:80
// End:

/*
  box.go
  Description: Definition of boxes as reusable components to be used in TikZ
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

// TikZ code to generate a box: first, the coordinate is created; next, the text
// is shown centered on the given coordinate with the width and height specified
const tikzBox = `{{.GetCoord}}
\draw {{.GetLabel}} node [rounded corners, rectangle, minimum width={{.GetMinWidth}}, minimum height = {{.GetMinHeight}}{{.GetDraw}}] {{.GetText}};`

// types
// ----------------------------------------------------------------------------

// A box is centered in a coordinate and can show any text (including the empty
// string). Additionally, the minimum width and height have to be given as
// strings, the reason being that they can consist of TikZ formulae. A box can
// be drawn or not. When not drawn, it serves to properly locate items in a
// bigger picture
type Box struct {
	Coordinate
	minWidth, minHeight string
	draw                bool
	text                string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a box given a coordinate and the text to show
func NewBox(coord Coordinate, draw bool, minWidth, minHeight, text string) Box {
	return Box{Coordinate: coord,
		minWidth:  minWidth,
		minHeight: minHeight,
		draw:      draw,
		text:      text}
}

// return a valid specification of a box with no error if all the keys given in
// dict are correct for defining a box. Otherwise, return an error. If an error
// is returned, the contents of the box are undetermined
//
// A dictionary is correct if and only if it correctly defines a text box (see
// VerifyTextDict) and also provides values from the minimum width and height
// with the keywords "minwidth" and "minheight" which should be given as strings
// as they can consist of LaTeX formulae, and an integer value (interpreted as a
// boolean) for "draw" indicating whether the box should be displayed or not
func VerifyBoxDict(dict map[string]interface{}) (Box, error) {

	// first of all, verify that this dictionary correctly provides information
	// for creating a text box
	var err error
	var text Text
	if text, err = VerifyTextDict(dict); err != nil {
		return Box{}, fmt.Errorf("A text box was not properly defined while creating a box: %v", err)
	}

	// now, beyond the definition of a text box, the mandatory keys are given
	// next
	mandatory := []string{"minwidth", "minheight", "draw"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then raise an error and
		// exit
		if _, ok := dict[key]; !ok {
			return Box{}, fmt.Errorf("Mandatory key '%v' for defining a box not found", key)
		}
	}

	// make also sure that parameters are given with the right type
	if _, ok := dict["minwidth"].(string); !ok {
		return Box{}, errors.New("the minimum width of a box should be given as a string")
	}
	if _, ok := dict["minheight"].(string); !ok {
		return Box{}, errors.New("the minimum height of a box should be given as a string")
	}
	if _, ok := dict["draw"].(int); !ok {
		return Box{}, errors.New("the value for 'draw'ing a box should be given as an integer which is then interpreted as a bool value, where non-null values represent true")
	}

	// otherwise, the dictionary is correct
	return Box{Coordinate: text.Coordinate,
		minWidth:  dict["minwidth"].(string),
		minHeight: dict["minheight"].(string),
		draw:      dict["draw"].(int) != 0,
		text:      dict["text"].(string)}, nil
}

// -- Box

// Return the coordinate of this box
func (b Box) GetCoord() string {
	return fmt.Sprintf("%v", b.Coordinate)
}

// Return the label of the coordinate of this box
func (b Box) GetLabel() string {
	return fmt.Sprintf("(%v)", b.label)
}

// Return the minimum width of this box
func (b Box) GetMinWidth() string {
	return fmt.Sprintf("%v", b.minWidth)
}

// Return the minimum height of this box
func (b Box) GetMinHeight() string {
	return fmt.Sprintf("%v", b.minHeight)
}

// Return the string ", draw" if and only if this box has to be drawn and ""
// otherwise
func (b Box) GetDraw() string {
	if b.draw {
		return ", draw"
	}
	return ""
}

// Return the text to show in this box
func (b Box) GetText() string {
	return fmt.Sprintf("{%v}", b.text)
}

// return a TikZ representation of a box
func (b Box) String() string {

	// create a template with the TikZ code for showing a box
	tpl, err := template.New("box").Parse(tikzBox)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, b); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

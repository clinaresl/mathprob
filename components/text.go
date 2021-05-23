/*
  text.go

  Description: Definition of text as reusable components to be used in TikZ
			   drawings

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

// TikZ code to generate text: first, the coordinate is created; next, the text
// is shown centered on the coordinate
const tikzText = `{{.GetCoord}}
\draw {{.GetLabel}} node {{.GetText}};`

// types
// ----------------------------------------------------------------------------

// Text is written centered in a given coordinate and can hold any string
// (including other LaTeX commands that can affect the look&feel)
type Text struct {
	Coordinate
	text string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a text box given a coordinate and the text to show
func NewText(coord Coordinate, text string) Text {
	return Text{Coordinate: coord, text: text}
}

// return a valid text and no error if all the keys given in dict are correct
// for defining a text box. Otherwise, return an error. If an error is returned
// the contents of Text are undetermined
//
// A dictionary is correct if and only if it correctly defines a coordinate (see
// VerifyCoordinateDict) and also provides text to be displayed (which can be an
// empty string or might contain LaTeX commands to affect the appearance of the
// text) with the keyword "text"
func VerifyTextDict(dict map[string]interface{}) (Text, error) {

	// first of all, verify that this dictionary correctly provides information
	// for creating a coordinate
	var err error
	var coord Coordinate
	if coord, err = VerifyCoordinateDict(dict); err != nil {
		return Text{}, fmt.Errorf("A coordinate was not properly defined while creating a text box: %v", err)
	}

	// now, beyond the definition of a coordinate, the mandatory keys are given
	// next
	mandatory := []string{"text"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then raise an error and
		// exit
		if _, ok := dict[key]; !ok {
			return Text{}, fmt.Errorf("Mandatory key '%v' for defining a text box not found", key)
		}
	}

	// make also sure that parameters are given with the right type
	if text, ok := dict["text"].(string); ok {
		return Text{Coordinate: coord, text: text}, nil
	}

	// if the type assertion failed then return an error
	return Text{}, errors.New("the text to show in a text box should be given as a string")
}

// methods
// ----------------------------------------------------------------------------

// Return the coordinate of this text box
func (t Text) GetCoord() string {
	return fmt.Sprintf("%v", t.Coordinate)
}

// Return the label of the coordinate of this text box
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

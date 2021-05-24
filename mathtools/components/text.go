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
	"log"
	"text/template"
)

// constants
// ----------------------------------------------------------------------------

// TikZ code to generate text: text is written with a node so that options and a
// label can be specified in addition to the text to write. Note that the tex
// can contain any prefix for the size as well
const tikzText = `\node [{{.GetOptions}}] ({{.GetLabel}}) { {{.GetText}} };`

// types
// ----------------------------------------------------------------------------

// Text is written with a node command so that options and a label can be
// attached to the text to write. Note that the text to write can be preceded of
// other LaTeX commands such as the size of the text ---or other effects, such
// as bold, italic, etc.
type Text struct {
	options string
	label   string
	text    string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a text given the three fields
func NewText(options, label, text string) Text {
	return Text{
		options: options,
		label:   label,
		text:    text,
	}
}

// return a valid text and no error if all the keys given in dict are correct
// for defining a text box. Otherwise, return an error. If an error is returned
// the contents of Text are undefined
//
// No parameter is mandatory ---not even the text itself.
func VerifyTextDict(dict map[string]interface{}) (Text, error) {

	// now, copy the values of the feasible parameters if any are given ---note
	// that none is mandatory
	var ok bool
	var options, label, text string
	for key, value := range dict {

		switch key {
		case "options":
			if options, ok = value.(string); !ok {
				return Text{}, errors.New("The options of a text box should be given as a string")
			}
		case "label":
			if label, ok = value.(string); !ok {
				return Text{}, errors.New("The label of a text box should be given as a string")
			}
		case "text":
			if text, ok = value.(string); !ok {
				return Text{}, errors.New("The text of a text box should be given as a string")
			}
		default:
			log.Printf("The parameter '%v' is not acknowledged for creating a text box and it will be ignored")
		}
	}

	// at this point, the arguments have been verified, so that a new Text is returned
	return Text{
		options: options,
		label:   label,
		text:    text,
	}, nil
}

// methods
// ----------------------------------------------------------------------------

// Return the coordinate of this text box
func (t Text) GetOptions() string {
	return t.options
}

// Return the label of the coordinate of this text box
func (t Text) GetLabel() string {
	return t.label
}

// Return the text to show of this text box
func (t Text) GetText() string {
	return t.text
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

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

// TikZ code to generate text: text can be written with a node so that options
// and a label can be specified in addition to the text to write. Note that the
// tex can contain any prefix for the size as well
const tikzText = `\node [{{.GetOptions}}] ({{.GetLabel}}) { {{.GetText}} };`

// In addition, text can be written specifically located at one label which is
// computed separtely using a draw operation
const tikZLabeledText = `\draw ({{.GetLabel}}) node [{{.GetOptions}}] { {{.GetText}} };`

// Also, text can be positioned specifically at one label computed within the
// text operation using a draw operation
const tikZCoordinatedText = `{{.Coordinate}}
\draw ({{.GetLabel}}) node [{{.GetOptions}}] { {{.GetText}} };`

// types
// ----------------------------------------------------------------------------

// Text can be written with a node command so that options and a label can be
// attached to the text to write. Note that the text to write can be preceded of
// other LaTeX commands such as the size of the text ---or other effects, such
// as bold, italic, etc.
type Text struct {
	options string
	label   string
	text    string
}

// But text can be also written at one specific location (computed separately)
// denoted with a label, using whichever options and text
type LabeledText struct {
	Text
}

// Moreover, text can be located at one specific coordinate which is given to
// the text
type CoordinatedText struct {
	Coordinate
	Text
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

// Create a new instance of text at one specific label computed separately given
// the three fields
func NewLabeledText(options, label, text string) LabeledText {
	return LabeledText{
		Text{
			options: options,
			label:   label,
			text:    text,
		},
	}
}

// Create a new instance of text at one specific label computed within the text
// operation given the three fields. Note that the label has not be given as
// coordinates are labeled necessarily, and the text is written at the label
// computed using the same coordinate
func NewCoordinatedText(coord Coordinate, options, text string) CoordinatedText {
	return CoordinatedText{
		Coordinate: coord,
		Text: Text{
			options: options,
			label:   coord.GetLabel(),
			text:    text,
		},
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

// -- Text

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

// -- LabeledText

// return a TikZ representation of a text box
func (t LabeledText) String() string {

	// create a template with the TikZ code for showing a text box
	tpl, err := template.New("text").Parse(tikZLabeledText)
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

// -- CoordinatedText

// Coordinated text consists of both a coordinate and text and both have labels,
// so there is a conflict which is easily solved by returning any of the labels
// ---as both have to be the same
func (t CoordinatedText) GetLabel() string {
	return t.Coordinate.label
}

// return a TikZ representation of a text box
func (t CoordinatedText) String() string {

	// create a template with the TikZ code for showing a text box
	tpl, err := template.New("text").Parse(tikZCoordinatedText)
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

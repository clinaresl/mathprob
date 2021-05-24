// -*- coding: utf-8 -*-
// line.go
//
// Description: Definition of lines or series of them as reusable components to
//              be used in TikZ drawings
// -----------------------------------------------------------------------------
//
// Started on <lun 24-05-2021 07:04:28.427050044 (1621832668)>
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
	"regexp"
	"text/template"

	"github.com/clinaresl/mathprob/helpers"
)

// constants
// ----------------------------------------------------------------------------

// TikZ code to generate lines as a serie of segments. Note that the number of
// segments to draw is undetermined but at least two references should be given,
// either explicitly or as formulas or using labels
const tikzLine = `\draw [{{.GetOptions}}] {{.GetSegments}};`

// The sequence of end points is completed using the following text templates.
// The first one is used to link the first two end points; the next one is used
// to add new end points
const tikzFirstSegment = `({{.GetReference0}}) -- ({{.GetReference1}})`
const tikzNextSegment = ` -- ({{.GetNextReference}})`

// types
// ----------------------------------------------------------------------------

// Any line has opitons which consist of a comma-separated list of options in a
// string
type BaseLine struct {
	options string
}

// So that a line consists of a list of segments along with some options to draw
// it. Each endpoint is identified with a string which might represent a
// coordinate explicitly given, or a formula, or as the name of a label. A line
// to be correct should contain at least two endpoints
type Line struct {
	refs []string
	BaseLine
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a line given an arbitrary number of segments. Note
// that the options are specified through a dedicated service.
//
// A line to be correct should consist at least two end-points but it is
// possible to provide an arbitrary number of them
func NewLine(ref0, ref1 string, refs ...string) Line {

	// join the first two mandatory references to all the others
	references := append([]string{ref0, ref1}, refs...)

	// and use them to create a line
	return Line{
		refs: references,
	}
}

// return a valid specification of a line with no error if all the keys given in
// dict are correct for defining a line. Otherwise, return an error. If an error
// is returned, the contents of the line are undefined.
//
// A dictionary is correct if and only if it correctly defines a sequence of
// segments, each one identified with the next number after the keyword
// "ref", i.e., "ref0", "ref1", "ref2", etc. These are the
// only mandatory arguments and there can be an arbitrary number of them. In
// addition, it is also possible to specify arbitrary options as a string
func VerifyLineDict(dict map[string]interface{}) (Line, error) {

	// First of all, ensure that all mandatory parameters are given and that
	// they are of the correct type. Create slices for both mandatory and all
	// arguments. Note that, still, there can be more endpoints: "ref2",
	// "ref3", etc.
	all := []string{"ref0", "ref1", "options"}
	mandatory := []string{"ref0", "ref1"}

	// verify that all mandatory arguments are given in the dictionary
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then immediately raise
		// an error
		if _, ok := dict[key]; !ok {
			return Line{}, fmt.Errorf("Mandatory key '%v' for defining a line not found", key)
		}
	}

	// now ensure that the mandatory parameters are of the right type
	var refs []string
	if ref0, ok := dict["ref0"].(string); !ok {
		return Line{}, errors.New("The first end-point of a line should be given as a string")
	} else {
		refs = append(refs, ref0)
	}
	if ref1, ok := dict["ref1"].(string); !ok {
		return Line{}, errors.New("The second end-point of a line should be given as a string")
	} else {
		refs = append(refs, ref1)
	}

	// and also check whether more than two end points were given, ... if so,
	// process them as well. Note that after this process, nextidx gets the
	// integer value of the last reference point plus one
	nextidx := 2
	for {

		// the next reference to look for is of the form "refi" where i is the
		// next integer index
		nextref := fmt.Sprintf("ref%v", nextidx)

		// in case the next end-point has been found, process it
		if _, ok := dict[nextref]; ok {
			if _, ok := dict[nextref].(string); !ok {
				return Line{}, errors.New("Every end-point of a line should be given as a string")
			}
			refs = append(refs, dict[nextref].(string))
		} else {

			// if no more end-points have been found, exit of the loop
			break
		}

		// and increment the index
		nextidx++
	}

	// now, perform the same operation with the optional parameters
	var options string
	if _, ok := dict["options"]; ok {
		if _, ok := dict["options"].(string); !ok {
			return Line{}, errors.New("The options of a line should be given as a string")
		}
		options = dict["options"].(string)
	}

	// in case any other arguments were given, but they are not acknowledged,
	// issue a warning
	for key, _ := range dict {

		// if a key is not found in the list of all arguments then it might be
		// an unnecessary argument unless ...
		if !helpers.Find(key, all) {

			// it is one of the end-points used in the definition of a segment
			// of the line
			r, _ := regexp.Compile(`ref(\d+)$`)
			if r.MatchString(key) {

				// get the integer index. Note that if a match has happened,
				// there can be only one match, right at the end of the string
				idx, _ := helpers.Atoi(r.FindStringSubmatch(key)[1])

				// now, if this index went beyond the last index processed then
				// it is unnecessary, most likely because an intermediate end
				// point has not been specified
				if idx > nextidx {
					log.Printf("A reference to an end-point '%v' has been found, but at least the previous one was not given. The last processed end-point has index %v", key, nextidx-1)
				}
			} else {

				// if this was not a reference to an end-point, then it is
				// clearly an unnecessary argument
				log.Printf("The parameter '%v' is not acknowledged for creating a line and it will be ignored")
			}
		}
	}

	// At this point, the dictionary is correct, return a valid line
	return Line{
		refs:     refs,
		BaseLine: BaseLine{options: options},
	}, nil
}

// The following function is used to return the index to the next reference
// point. In the absence of static variables in Go this is done using closures
func nextEndPoint() (counter func() int) {

	// the last used index is #1
	idx := 1

	// the function just simply increment it
	return func() int {
		idx++
		return idx
	}
}

// methods
// ----------------------------------------------------------------------------

// --BaseLine

// Set the options of a line
func (line *BaseLine) SetOptions(options string) {
	line.options = options
}

// Get the options used
func (line BaseLine) GetOptions() string {
	return line.options
}

// --Line

// Return the reference of the first end-point
func (line Line) GetReference0() string {

	return fmt.Sprintf("%v", line.refs[0])
}

// Return the reference of the second end-point
func (line Line) GetReference1() string {

	return fmt.Sprintf("%v", line.refs[1])
}

// Return the reference of the next end-point
func (line Line) GetNextReference() string {

	nextidx := nextEndPoint()()
	return fmt.Sprintf("%v", line.refs[nextidx])
}

// GetSegments returns a string with the sequence of end points of the line
func (line Line) GetSegments() string {

	// create a template with the TikZ code for showing the segment created by
	// the first two end-points
	tpl, err := template.New("line").Parse(tikzFirstSegment)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, line); err != nil {
		log.Fatal(err)
	}

	// next, in case there are more end-points add them to the output using
	// substitutions with the template used for adding them
	if len(line.refs) > 2 {
		tpl, err = template.New("nextline").Parse(tikzNextSegment)
		if err != nil {
			log.Fatal(err)
		}

		// and now make the appropriate substitution. Note that the execution of the
		// template is written to a string
		if err := tpl.Execute(&tplOutput, line); err != nil {
			log.Fatal(err)
		}
	}

	// and return the resulting string
	return tplOutput.String()

}

// Finally, lines are stringers and these are the means provided for
// automatically reusing this component
func (line Line) String() string {

	// create a template with the TikZ code for showing a line
	tpl, err := template.New("line").Parse(tikzLine)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitution. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, line); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// Local Variables:
// mode:go
// fill-column:80
// End:

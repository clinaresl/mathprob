/*
  sequence.go
  Description: Provides services for automatically creating a sequence
  or a part of it
  -----------------------------------------------------------------------------

  Started on  <Wed Dec 12 20:44:45 2018 >
  Last update <>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by
  Login   <clinares@atlas>
*/

package mathtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"text/template"
	"time"

	"github.com/clinaresl/mathprob/mathtools/components"
)

// constants
// ----------------------------------------------------------------------------

// There are four different types of sequences: "first", "last", "none" or
// "both" if the first, last, none or both numbers of the sequence have to be
// shown
const (
	SEQNONE int = iota
	SEQFIRST
	SEQLAST
	SEQBOTH
)

// the TikZ code for generating arbitrary sequences is shown next. Note that it
// makes use of LaTeX/TikZ components
const latexSequenceCode = `\begin{minipage}{0.25\linewidth}
	\begin{center}
		\begin{tikzpicture}

			% draw the sequence
			{{.GetTikZSequence}}

			% draw an invisible bounding box to properly align all sequences

		\end{tikzpicture}
	\end{center}
\end{minipage}
`

// types
// ----------------------------------------------------------------------------

// A Sequence consists of a type: "first", "last", "none" or "both" if either
// the first number has to be given, the last one, none of them, or both
// respectively. It consists of a number of items, each one greater or equal
// than a given threshold and less or equal than another bound.
type sequence struct {
	seqtype  int
	nbitems  int
	geq, leq int
}

// A Sequence in JSON format consists of a list containing the sequence with
// question marks which should be filled in by the student and another list with
// the answer. The struct also contains another field probtype to recognize this
// structure as a sequence problem in a JSSON file
type sequenceJSON struct {
	Probtype string   `json:"type"`
	Args     []string `json:"args"`
	Solution []string `json:"solution"`
}

// methods
// ----------------------------------------------------------------------------

// -- sequence
// ----------------------------------------------------------------------------

// return the instance of a specific sequence problem in JSON format as a slice
// of bytes. The receiver is assumed to have been fully verified so that it
// should be consistent
func (sequence sequence) GenerateJSONInstance() (data []byte, err error) {

	// from the given instance of a sequence, define the specific problem

	// determine the first number of the sequence ---even if it is not
	// displayed. If the interval [geq, leq] is too narrow to host nbitems,
	// immediately log a fatal error
	if 1+sequence.leq-sequence.geq < sequence.nbitems {
		return data, fmt.Errorf("It is not possible to fit %v different numbers taken from the range [%v, %v]",
			sequence.nbitems, sequence.geq, sequence.leq)
	}

	// The following expression takes into account not only the interval [geq,
	// leq] but also the number of items to display in the sequence
	rand.Seed(time.Now().UTC().UnixNano())
	number1 := sequence.geq + rand.Int()%(2+sequence.leq-sequence.nbitems-sequence.geq)

	// in case this sequence is of type SEQNONE, then randomly choose a position
	// in between to show a number
	pos := 1 + rand.Int()%(sequence.nbitems-2)

	// and now fill in the sequence along with the solution
	args := make([]string, sequence.nbitems)
	solution := make([]string, sequence.nbitems)
	for item := number1; item < number1+sequence.nbitems; item++ {

		// first, write the solution
		idx := item - number1
		solution[idx] = strconv.FormatInt(int64(item), 10)

		// now, depending on the position and type

		switch idx {

		case 0:
			if sequence.seqtype == SEQNONE || sequence.seqtype == SEQLAST {
				args[0] = "?"
			} else {
				args[0] = solution[0]
			}

		case pos:
			if sequence.seqtype == SEQNONE {
				args[pos] = solution[pos]
			} else {
				args[pos] = "?"
			}

		case sequence.nbitems - 1:
			if sequence.seqtype == SEQNONE || sequence.seqtype == SEQFIRST {
				args[idx] = "?"
			} else {
				args[idx] = solution[idx]
			}

		default:
			args[idx] = "?"
		}
	}

	// Marshal the problem in JSON format
	seqprob := &sequenceJSON{
		Probtype: "Sequence",
		Args:     args,
		Solution: solution}
	data, err = json.MarshalIndent(seqprob, "", "\t")

	// and if any error happened then return it immediately
	if err != nil {
		return data, err
	}

	// otherwise, return the data in JSON format
	return data, nil
}

// use the values stored in a sequence to determine the order of the reusable
// components to display the items. The output slice contains items of two
// types, either text (to show numbers) or box (to show answer boxes)
func (sequence sequence) getComponents() []components.ComponentId {

	// in case this sequence is of type SEQNONE, then randomly choose a position
	// in between to show a number
	pos := 1 + rand.Int()%(sequence.nbitems-2)

	// create the output slice
	order := make([]components.ComponentId, sequence.nbitems)

	for idx := 0; idx < sequence.nbitems; idx++ {

		// if it is either the first or last item and this specific element has been
		// requested
		if (idx == 0 && (sequence.seqtype == SEQFIRST || sequence.seqtype == SEQBOTH)) ||
			(idx == sequence.nbitems-1 && (sequence.seqtype == SEQLAST || sequence.seqtype == SEQBOTH)) {

			// then add a text for displaying a number
			order[idx] = components.TEXT
		} else {

			// otherwise, add a text also if this location is the one randomly
			// chosen in a SEQNONE sequence
			if sequence.seqtype == SEQNONE && idx == pos {
				order[idx] = components.TEXT
			} else {

				// in any other case, just add an answer box
				order[idx] = components.BOX
			}
		}
	}

	// and return the order computed so far
	return order
}

// return a valid LaTeX/TikZ representation of this sequence using TikZ
// components
func (sequence sequence) GetTikZSequence() string {

	// determine the first number of the sequence ---even if it is not displayed.
	// If the interval [geq, leq] is too narrow to host nbitems, immediately log a
	// fatal error
	if 1+sequence.leq-sequence.geq < sequence.nbitems {
		log.Fatalf("It is not possible to fit %v different numbers taken from the range [%v, %v]",
			sequence.nbitems, sequence.geq, sequence.leq)
	}

	// and also the number of necessary digits per item. This is computed as the
	// maximum number of digits that might be required ---in spite of the number
	// of digits actually needed. Because it is potentially possible to create
	// sequences with negative numbers then we consider both extrems
	nbdigits := max(float64(nbdigits(sequence.geq)),
		float64(nbdigits(sequence.leq)))

	// the first item to be drawn should be raised by half the height of zero
	// plus 1.5 the baselineskip, while all the other elements should be
	// horizontally aligned
	yshiftheight := 0.5
	yshiftskip := 1.5

	// first, locate a coordinate to mark the origin. This is done using the
	// reusable coordinate
	t := `{{.GetCoordinate (dict "label" "label0" "x" 0.0 "y" 0.0)}}`
	t += "\n"

	// The following expression takes into account not only the interval [geq,
	// leq] but also the number of items to display in the sequence
	rand.Seed(time.Now().UTC().UnixNano())
	number1 := sequence.geq + rand.Int()%(2+sequence.leq-sequence.nbitems-sequence.geq)

	// determine the order of reusable components to draw the sequence. Note
	// that text items are represented with boxes much the same anyway as their
	// usage help to properly locate the items in the whole picture
	for idx, component := range sequence.getComponents() {

		// if this is not the first element of the sequence, restart the y-shift
		if idx > 0 {
			yshiftheight = 0.0
			yshiftskip = 0.0
		}

		// now, depending on the type of reusable component
		switch component {
		case components.TEXT:

			// the number to show in this location is computed as the sum of the
			// first number and the index of this position in the sequence, but
			// other sequences can be created! Note that text is positioned
			// within an invisible box so that it is much easier to be
			// positioned within the whole picture
			t += fmt.Sprintf("{{.GetBox (dict \"label\" \"label%d\" \"formula\" `(label%d) + (%.2f\\zerowidth, %.2f\\zeroheight+%.2f\\baselineskip)` \"minwidth\" `%.2f\\zerowidth` \"minheight\" `\\zeroheight + \\baselineskip` \"draw\" 0 \"text\" `\\huge %d`)}}",
				1+idx, idx, 3+nbdigits, yshiftheight, yshiftskip, 2+nbdigits, idx+number1)

		case components.BOX:
			t += fmt.Sprintf("{{.GetBox (dict \"label\" \"label%d\" \"formula\" `(label%d) + (%.2f\\zerowidth, %.2f\\zeroheight+%.2f\\baselineskip)` \"minwidth\" `%.2f\\zerowidth` \"minheight\" `\\zeroheight + \\baselineskip` \"draw\" 1 \"text\" \"\")}}",
				1+idx, idx, 3+nbdigits, yshiftheight, yshiftskip, 2+nbdigits)

		default:
			log.Fatal("Unexpected type of a reusable component in a sequence")
		}

		// and move to the next line!
		t += "\n"
	}

	// now, execute this template with a masterFile
	var err error
	var result bytes.Buffer
	var masterFile MasterFile
	if result, err = masterFile.MasterToBufferFromTemplate(t); err != nil {
		log.Fatalf("Error when executing the template for creating a sequence: %v", err)
	}

	return result.String()
}

// Return TikZ code that represents a sequence
func (sequence sequence) execute() string {

	// create a template with the TikZ code for showing this sequence
	tpl, err := template.New("sequence").Parse(latexSequenceCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, sequence); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

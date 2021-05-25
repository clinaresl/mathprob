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
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"text/template"
	"time"

	"github.com/clinaresl/mathprob/helpers"
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
const latexSequenceCode = `\begin{minipage}{\linewidth}
    \begin{center}
        \begin{tikzpicture}

            % draw the sequence
            {{.GetTikZPicture}}

        \end{tikzpicture}
    \end{center}
\end{minipage}
`

const tikZSequenceCode = `% --- Coordinates ----------------------------------------------------

        % the lower-left corner is located at (0,0)
{{.GetBottomLabel}}

        % text boxes (either empty or with a hint) have a separation between
        % them equal to epsilon (which here equals 0.5 the width of a digit). To
        % avoid consecutive sequences to collide, twice epsilon is left from the
        % lower-left corner of the bounding box to start the sequence. Since
        % each text box has a width equal to the number of digits to show plus 2
        % (i.e., the additional space of the width of a digit to each side) the
        % first textbox is centered at 2epsilon + (2+nbdigits)/2. Since
        % epsilon=0.5, the previous expression yields: 1.0 + (2+nbdigits)/2
{{.GetFirstLabel}}

        % The distance between the centers of two consecutive textboxes equals
        % the width of any text box plus epsilon (the little space intentionally
        % left between text boxes), resulting in (2+nbdigits+epsilon). Thus, if
        % there are seq.nbitems in the whole sequence, then the distance from
        % the center of the first text box to the last one is equal to
        % (2+nbdigits+epsilon) * (seq.nbitems - 1)
{{.GetLastLabel}}

        % Finally, the upper-right corner is computed from the location of the
        % center of the last text box plus half the width of any text box. Since
        % the width of any text box is (2+nbdigits), the additional space from
        % the center of the last box equals (2+nbdigits)/2
{{.GetRightLabel}}

        % --- Bounding Box ----------------------------------------------------

        % the bounding box is drawn between the lower-left and upper-right
        % coordinates
{{.GetBoundingBox}}

        % ---------------------------------------------------------------------

        % --- Sequence --------------------------------------------------------

        % show all elements of the sequence
{{.GetSequenceItems}}
        % ---------------------------------------------------------------------
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

// A sequence is drawn using TikZ reusable components only. It cconsists of the
// bounding box along with its two coordinates (lower-left and upper-right), and
// additional coordinates for placing each text box which might be empty (it has
// to be filled in by the student) or not ---because it is a hint.
type sequenceTikZ struct {

	// The lower-left coordinate is inserted first to position other coordinates
	// wrt it
	bottom components.Coordinate

	// The first item and last items of the sequence are placed using specific
	// coordinates using only the information from the lower-left coordinate
	first, last components.Coordinate

	// the bounding box is drawn using two coordinates for the lower-left and
	// upper-right. Note that it is implemented as a plain rectangle (instead of
	// a coordinated rectangle), because coordinates are computed separately
	right components.Coordinate
	bBox  components.Rectangle

	// the items of the sequence are stored as text components which might be
	// empty or not, each one located at a different coordinate which is
	// computed from the first cell
	coords []components.Coordinate
	cells  []components.LabeledText
}

// methods
// ----------------------------------------------------------------------------

// --sequenceTikZ

// Generates the TikZ code necessary for positioning the lower-left corner of
// the bounding box
func (tikz sequenceTikZ) GetBottomLabel() string {

	// Coordinates draw themselves
	return fmt.Sprintf("%v", tikz.bottom)
}

// Generates the TikZ code necessary for positioning the coordinate at the
// center of the first item of the sequence
func (tikz sequenceTikZ) GetFirstLabel() string {

	// Coordinates draw themselves
	return fmt.Sprintf("%v", tikz.first)
}

// Generates the TikZ code necessary for positioning the coordinate at the
// center of the last item of the sequence
func (tikz sequenceTikZ) GetLastLabel() string {

	// Coordinates draw themselves
	return fmt.Sprintf("%v", tikz.last)
}

// Generates the TikZ code necessary for positioning the lower-left corner of
// the bounding box
func (tikz sequenceTikZ) GetRightLabel() string {

	// Coordinates draw themselves
	return fmt.Sprintf("%v", tikz.right)
}

// Generates the TikZ code necessary for positioning the bounding box
func (tikz sequenceTikZ) GetBoundingBox() string {

	// Bounding box draw themselves
	return fmt.Sprintf("%v", tikz.bBox)
}

// Generates the TikZ code necessary for positioning all items of the sequence,
// either empty cells or hints
func (tikz sequenceTikZ) GetSequenceItems() string {

	// Use a btyes buffer to append the strings of each cell
	var output bytes.Buffer

	// First, add all coordinates
	for _, coord := range tikz.coords {
		fmt.Fprintf(&output, "%v\n", coord)
	}

	// Draw all text boxes in the cells stored in this pict
	for _, cell := range tikz.cells {
		fmt.Fprintf(&output, "%v\n", cell)
	}

	return output.String()
}

// Return the LaTeX/TikZ commands that show up the picture stored in the
// receiver
func (seq sequenceTikZ) execute() string {

	// create a template with the TikZ code for showing this picture
	tpl, err := template.New("sequenceTikZ").Parse(tikZSequenceCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of
	// the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, seq); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// -- sequence

// return the instance of a specific sequence problem that can be marshalled in
// JSON format. The receiver is assumed to have been fully verified so that it
// should be consistent.
//
// The result is given with a list with as many elements as items in the
// sequence where "?" signals those locations that have to be guessed by the
// student
func (seq sequence) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// determine the first number of the sequence ---even if it is not
	// displayed. If the interval [geq, leq] is too narrow to host nbitems,
	// immediately log a fatal error
	if 1+seq.leq-seq.geq < seq.nbitems {
		return problemJSON{}, fmt.Errorf("It is not possible to fit %v different numbers taken from the range [%v, %v]",
			seq.nbitems, seq.geq, seq.leq)
	}

	// The following expression takes into account not only the interval [geq,
	// leq] but also the number of items to display in the sequence
	number1 := seq.geq + rand.Int()%(2+seq.leq-seq.nbitems-seq.geq)

	// in case this sequence is of type SEQNONE, then randomly choose a position
	// in between to show a number, unless there are only two items in which
	// case randomly chose any
	var pos int
	if seq.nbitems <= 2 {
		pos = rand.Int() % (seq.nbitems)
	} else {
		pos = 1 + rand.Int()%(seq.nbitems-2)
	}

	// and now fill in the sequence along with the solution
	args := make([]string, seq.nbitems)
	solution := make([]string, seq.nbitems)
	for item := number1; item < number1+seq.nbitems; item++ {

		// first, write the solution
		idx := item - number1
		solution[idx] = strconv.FormatInt(int64(item), 10)

		// now, depending on the position and type

		switch idx {

		case 0:
			if seq.seqtype == SEQNONE || seq.seqtype == SEQLAST {
				args[0] = "?"
			} else {
				args[0] = solution[0]
			}

		case pos:
			if seq.seqtype == SEQNONE {
				args[pos] = solution[pos]
			} else {
				args[pos] = "?"
			}

		case seq.nbitems - 1:
			if seq.seqtype == SEQNONE || seq.seqtype == SEQFIRST {
				args[idx] = "?"
			} else {
				args[idx] = solution[idx]
			}

		default:
			args[idx] = "?"
		}
	}

	// and return the problem along with its solution
	return problemJSON{
		Probtype: "Sequence",
		Args:     args,
		Solution: solution}, nil
}

// return a valid LaTeX/TikZ representation of this sequence using TikZ
// components
func (seq sequence) GetTikZPicture() string {

	// -- operands randomly determine the values of the operands. For this, the
	// service that generates problems is the one that can marshal them into
	// JSON format. The numbers of the sequence are given in Args, where a
	// question mark is a number that has to be guessed by the student
	instance, err := seq.generateJSONProblem()
	if err != nil {
		log.Fatalf(" Fatal error while generating a valid sequence: %v", err)
	}

	// in spite of the values geq and leq, it is good to compute the maximum
	// number of digits in each box, so that they look the same (and hence, no
	// additional clues are given to the student ;) )
	nbdigits := 0.0
	for _, item := range instance.Solution {
		if value, err := helpers.Atoi(item); err != nil {
			panic(fmt.Sprintf("Fatal error in the generation of a sequence: %v", err))
		} else {
			if nbd := helpers.NbDigits(value); float64(nbd) > nbdigits {
				nbdigits = float64(nbd)
			}
		}
	}

	// -- Coordinates

	// bottom is the lower-left corner of the bounding box and this coordinate
	// is used to reference others which are computed wrt it
	bottom := components.NewCoordinate(components.Point{
		X: 0.0,
		Y: 0.0,
	}, "bottom")

	// first is the center of the location of the first box
	first := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(bottom) + (%v\zerowidth, 0.5\zeroheight+1.5\baselineskip)$`,
			1.0+(2+nbdigits)/2.0)),
		"first",
	)

	// the last element is placed leaving as much space as required to place
	// intermediate text boxes
	last := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(first) + (%v*\zerowidth, 0.0)$`,
			(2.5+nbdigits)*float64((seq.nbitems-1)))),
		"last",
	)
	right := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(last) + (%v\zerowidth, 0.5\zeroheight + 0.5\baselineskip)$`,
			(2+nbdigits)/2.0)),
		"right",
	)

	// -- Bounding box

	// The bounding box is delimited by bottom and right, as usual
	bBox := components.NewRectangle("bottom", "right")
	bBox.SetOptions("white")

	// -- items

	// the items to draw are given in the Args field of this specific instance.
	// If a question mark is given, then it should be replaced with an empty
	// text box; otherwise, the specified number is shown
	var coords []components.Coordinate
	var cells []components.LabeledText
	for idx, item := range instance.Args {

		// Create two ancilliary variables to store the coordinate and aspect of
		// the next cell
		var box components.LabeledText

		// in spite of the contents, the next cell is located at
		coord := components.NewCoordinate(
			components.Formula(fmt.Sprintf(`$(first) + (%v\zerowidth, 0)$`,
				float64(idx)*(2.5+nbdigits))),
			fmt.Sprintf("cell%v", idx),
		)

		// if this is a question mark
		if item == "?" {

			// then add an empty text box
			box = components.NewLabeledText(
				fmt.Sprintf(`rounded corners, rectangle, minimum width=%v*\zerowidth, minimum height = \zeroheight + \baselineskip, draw`,
					2.0+nbdigits,
				),
				fmt.Sprintf("cell%v", idx),
				"",
			)
		} else {

			// otherwise, add the number itself
			box = components.NewLabeledText(
				"",
				fmt.Sprintf("cell%v", idx),
				`\huge `+item)
		}

		// and add the new box and its coordinates
		coords = append(coords, coord)
		cells = append(cells, box)
	}

	// And put all this elements together to show up the picture of a sequence
	seqPicture := sequenceTikZ{
		bottom: bottom,
		first:  first,
		last:   last,
		right:  right,
		bBox:   bBox,
		coords: coords,
		cells:  cells,
	}

	// and return the TikZ code necessary for drawing the problem
	return seqPicture.execute()
}

// Return TikZ code that represents a sequence
func (seq sequence) execute() string {

	// create a template with the TikZ code for showing this sequence
	tpl, err := template.New("sequence").Parse(latexSequenceCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, seq); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

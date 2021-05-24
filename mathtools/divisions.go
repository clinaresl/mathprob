/*
  divisions.go
  Description: Provides services for automatically creating divisions

  -----------------------------------------------------------------------------

  Started on  <Sat May 25 17:26:25 2019 >
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

// the TikZ code for generating divisions with any parameters is shown
// below
// the TikZ code for generating divisions with any parameters is shown
// below
const latexDivisionCode = `\begin{minipage}{0.25\linewidth}
  \begin{center}
    \begin{tikzpicture}

        % draw the division
        {{.GetTikZPicture}}

    \end{tikzpicture}
  \end{center}
\end{minipage}
`

const tikZDivisionCode = `% --- Coordinates -------------------------------------------------------
{{.GetDivFirstLabel}}
{{.GetDivNextLabels}}
        % -----------------------------------------------------------------------

        % --- Ancilliary reference points
{{.GetDivLine}}
        % -----------------------------------------------------------------------

        % --- Bounding Box ------------------------------------------------------
{{.GetDivBoundingBox}}
        % -----------------------------------------------------------------------
        % show the box enclosing the divisor
{{.GetDivSplitBox}}
        % show the box for writing the quotient
{{.GetDivAnswer}}
        % -----------------------------------------------------------------------
        
        % --- Text ------------------------------------------------------------

        % Dividend
{{.GetDivDividend}}
        % Divisor
{{.GetDivDivisor}}
        % -----------------------------------------------------------------------
`

const latexAnswerCode = `        \node [rounded corners, rectangle, minimum width={{.GetDivWidth}}*\zerowidth, minimum height = \zeroheight+\baselineskip, draw, below=0.15 cm of label3] {};
`

const latexDivOperandCode = `        \node [right=0.0 cm of {{.GetDivRef}}] ({{.GetDivLabel}}) {\huge {{.GetDivValue}}};
`

// types
// ----------------------------------------------------------------------------

// The answer should be written in a box explicitly shown in the
// exercise. The only important parameter for drawing it is its width
// which is computed as a factor of the width of a digit in LaTeX
type latexAnswer struct {
	width float64
}

// Operands (either the dividend or the divisor) are characterized by
// its value, and its location which is computed with respect to a
// reference coordinate which is identified by its name, and a label
type latexDivOperand struct {
	value      int
	ref, label string
}

// A division is characterized by its coordinates, a bounding box
// surrounding all the available area for solving the exercise, the
// box enclosing the divisor, and also the operands
type divisionProblem struct {

	// the first label is computed explicitly whereas the next two labels are
	// computed with respect to the previous ones using formulas. All of them
	// are implemented using the reusable components of TikZ coordinates
	label1         components.Coordinate
	label2, label3 components.Coordinate

	// likewise the coordinate for determining the location of the first line is
	// computed with respect to another coordinate and hence a formula is used.
	// There might be an arbitrary number of points but other points are
	// computed only wrt the location of the first line
	line1 components.Coordinate

	// the bounding box surrounding all the necessary area for solving the
	// exercise is defined next and it is computed with two formulas that
	// specify the lower left and upper right corners
	bBox components.CoordinatedRectangle

	// the box surrounding the dividend consists of a path drawn between three
	// coordinates whose location is determined using formulas
	sBox components.Line

	// the answer should be written within a box explicitly shown
	answer latexAnswer

	// finally, both operands, are created next
	dividend latexDivOperand
	divisor  latexDivOperand
}

// The formal definition of a division problem is given below. It is defined
// with the number of digits of the dividend, divisor and quotient
type division struct {
	nbdvdigits int
	nbdrdigits int
	nbqdigits  int
}

// methods
// ----------------------------------------------------------------------------

// -- division

// return the instance of a specific division problem that can be marshalled in
// JSON format. The receiver is assumed to have been fully verified so that it
// should be consistent
func (div division) generateJSONProblem() (problemJSON, error) {

	// First, verify that parameters are correct. If they are not, take the best
	// action
	if div.nbqdigits < div.nbdvdigits-div.nbdrdigits {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", div.nbqdigits, div.nbdvdigits, div.nbdrdigits, div.nbdvdigits-div.nbdrdigits)
		div.nbqdigits = div.nbdvdigits - div.nbdrdigits
	}

	if div.nbqdigits > div.nbdvdigits-div.nbdrdigits+1 {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", div.nbqdigits, div.nbdvdigits, div.nbdrdigits, div.nbdvdigits-div.nbdrdigits+1)
		div.nbqdigits = div.nbdvdigits - div.nbdrdigits + 1
	}

	// create two slices: one for storing the instance of this problem in the
	// order: dividend, divisor, quotient and remainder where those parts that
	// should be filled in by the student are marked with question marks "?";
	// and another one with the full solution
	args := make([]string, 4)
	solution := make([]string, 4)

	// now, generate numbers in their corresponding range
	var dividend, divisor, quotient int
	for helpers.NbDigits(quotient) != div.nbqdigits || quotient == 0 {
		dividend = helpers.RandN(div.nbdvdigits)
		divisor = helpers.RandN(div.nbdrdigits)
		quotient = dividend / divisor
	}

	// now, copy the arguments and the full solution
	solution[0] = strconv.FormatInt(int64(dividend), 10)
	solution[1] = strconv.FormatInt(int64(divisor), 10)
	solution[2] = strconv.FormatInt(int64(quotient), 10)
	solution[3] = strconv.FormatInt(int64(dividend-divisor*quotient), 10)
	args[0] = solution[0]
	args[1] = solution[1]
	args[2] = "?"
	args[3] = "?"

	// and return the problem along with its solution
	return problemJSON{
		Probtype: "Division",
		Args:     args,
		Solution: solution}, nil
}

// -- answer

// Return the width of the answer box
func (answer latexAnswer) GetDivWidth() string {
	return fmt.Sprintf("%v", answer.width)
}

// Return the Tikz code for drawing answer boxes
func (answer latexAnswer) String() string {

	// use the template defined for drawing answer boxes
	tpl, err := template.New("division").Parse(latexAnswerCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, answer); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// -- operands

// Return the reference point used for drawing this operand
func (operand latexDivOperand) GetDivRef() string {
	return fmt.Sprintf("%v", operand.ref)
}

// Return the label used for drawing this operand
func (operand latexDivOperand) GetDivLabel() string {
	return fmt.Sprintf("%v", operand.label)
}

// Return the of this operand
func (operand latexDivOperand) GetDivValue() string {
	return fmt.Sprintf("%v", operand.value)
}

// Return the Tikz code for drawing an operand
func (operand latexDivOperand) String() string {

	// use the template defined for drawing an operand
	tpl, err := template.New("division").Parse(latexDivOperandCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, operand); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// -- divisionProblem
// ----------------------------------------------------------------------------

// Generates the TikZ code necessary for positioning the first coordinate
func (division divisionProblem) GetDivFirstLabel() string {

	// The label locating the dividend is stored in label1 as a
	// components.Coordinate. It then just suffices to print it
	return fmt.Sprintf("%v", division.label1)
}

// Generates the TikZ code necessary for positioning other coordinates after the
// first one
func (division divisionProblem) GetDivNextLabels() string {

	// The labels used for locating the left margin of the bounding box
	// surrounding the divisor and also its lower margin are stored as
	// components.Coordinate in label2 and label3. It then just suffices to show
	// their location
	return fmt.Sprintf("%v\n%v", division.label2, division.label3)
}

// Generates the TikZ code necessary for positioning the first line of results
func (division divisionProblem) GetDivLine() string {

	// Coordinates draw themselves
	return fmt.Sprintf("%v", division.line1)
}

// Generates the TikZ code necessary for positioning the bounding box
// for solving the whole exercise
func (division divisionProblem) GetDivBoundingBox() string {

	// Bounding box draw themselves
	return fmt.Sprintf("%v", division.bBox)
}

// Generates the TikZ code necessary for drawing the split box
func (division divisionProblem) GetDivSplitBox() string {

	// Lines draw themselves
	return fmt.Sprintf("%v", division.sBox)
}

// Generates the TikZ code necessary for drawing the answer box
func (division divisionProblem) GetDivAnswer() string {

	// answers draw themselves
	return division.answer.String()
}

// Generates the TikZ code necessary for drawing the dividend
func (division divisionProblem) GetDivDividend() string {

	// answers draw themselves
	return division.dividend.String()
}

// Generates the TikZ code necessary for drawing the divisor
func (division divisionProblem) GetDivDivisor() string {

	// answers draw themselves
	return division.divisor.String()
}

// Execute the given division problem and returns legal TikZ code to represent
// it
func (div divisionProblem) execute() string {

	// create a template with the TikZ code for showing this
	// division problem
	tpl, err := template.New("division").Parse(tikZDivisionCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, div); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// return a valid LaTeX/TikZ representation of this sequence using TikZ
// components
func (div division) GetTikZPicture() string {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// --coordinates
	label1 := components.NewCoordinate(components.Point{
		X: 0.0,
		Y: 1 + 2.0*float64(div.nbqdigits) + 0.5,
	}, "label1")

	label2 := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(label1) + %v*(\zerowidth, 0.0)$`,
			2.0+float64(div.nbdvdigits))),
		"label2")

	label3 := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(label2) + (%v*\zerowidth, -\zeroheight)$`,
			0.5*(2+helpers.Max(float64(div.nbdrdigits), float64(div.nbqdigits))))),
		"label3")

	// --lines
	line1 := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(label2) + (-%v\zerowidth, -2*\zeroheight-0.15 cm)$`,
			2.0+float64(div.nbdvdigits))),
		"line1")

	// --bounding box
	bottom := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(line1) + %v*(0.0, -\zeroheight-\baselineskip-0.5/%v*\zeroheight)$`,
			2.0*float64(div.nbqdigits)-1.0,
			2.0*float64(div.nbqdigits)-1.0)),
		"bottom")
	right := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(line1) + %v*(0.0, -\zeroheight-\baselineskip-0.5/%v*\zeroheight)$`,
			2.0*float64(div.nbqdigits)-1.0,
			2.0*float64(div.nbqdigits)-1.0)),
		"right")
	bBox := components.NewCoordinatedRectangle(bottom, right)

	// --split box
	sBox := components.NewLine(`$(label2) + (0.0, \zeroheight)$`,
		`$(label2) + (0.0, -\zeroheight)$`,
		fmt.Sprintf(`$(label2) + %v*(\zerowidth, -\zeroheight/%v)$`,
			2.0+helpers.Max(float64(div.nbdrdigits), float64(div.nbqdigits)),
			2.0+helpers.Max(float64(div.nbdrdigits), float64(div.nbqdigits))))
	sBox.SetOptions("thick, rounded corners")

	// --answer
	answer := latexAnswer{
		width: 2.0 + helpers.Max(float64(div.nbdrdigits), float64(div.nbqdigits)),
	}

	// --operands
	dividend := latexDivOperand{
		ref:   "label1",
		label: "dividend",
	}
	divisor := latexDivOperand{
		ref:   "label2",
		label: "divisor",
	}

	// randomly determine the values of the operands

	// First, verify that parameters are correct. If they are not,
	// take the best action
	if div.nbqdigits < div.nbdvdigits-div.nbdrdigits {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", div.nbqdigits, div.nbdvdigits, div.nbdrdigits, div.nbdvdigits-div.nbdrdigits)
		div.nbqdigits = div.nbdvdigits - div.nbdrdigits
	}

	if div.nbqdigits > div.nbdvdigits-div.nbdrdigits+1 {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", div.nbqdigits, div.nbdvdigits, div.nbdrdigits, div.nbdvdigits-div.nbdrdigits+1)
		div.nbqdigits = div.nbdvdigits - div.nbdrdigits + 1
	}

	// now, generate numbers in their corresponding range
	var qvalue int
	for helpers.NbDigits(qvalue) != div.nbqdigits || qvalue == 0 {
		dividend.value = helpers.RandN(div.nbdvdigits)
		divisor.value = helpers.RandN(div.nbdrdigits)
		qvalue = dividend.value / divisor.value
	}

	// And put all this elements together to bring up the defintion of a division
	divProblem := divisionProblem{
		label1:   label1,
		label2:   label2,
		label3:   label3,
		line1:    line1,
		bBox:     bBox,
		sBox:     sBox,
		answer:   answer,
		dividend: dividend,
		divisor:  divisor,
	}

	// and return the TikZ code necessary for drawing this operation
	return divProblem.execute()
}

// Execute the given division instance and returns legal TikZ code to represent
// it
func (div division) execute() string {

	// create a template with the TikZ code for showing this
	// division problem
	tpl, err := template.New("division").Parse(latexDivisionCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, div); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

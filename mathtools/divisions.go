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
	"strconv"
	"text/template"
)

// constants
// ----------------------------------------------------------------------------

// the TikZ code for generating divisions with any parameters is shown
// below
const latexPreviousDivisionCode = `\begin{minipage}{0.25\linewidth}  
  \begin{center}
    \begin{tikzpicture}

        % --- Coordinates -------------------------------------------------------
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

    \end{tikzpicture}
  \end{center}
\end{minipage}
`
const latexCoordinateExplicitCode = `        \coordinate {{.GetDivLabel}} at {{.GetDivCoords}};`

const latexCoordinateFormulaCode = `        \coordinate {{.GetDivLabel}} at ({{.GetDivFormula}});
`

const latexBoundingBoxCode = `{{.GetDivBottom}}{{.GetDivRight}}        \draw [white] (bottom) rectangle (right);
`

const latexSplitBoxCode = `        \draw [thick, rounded corners] ({{.GetDivFirstCoord}}) -- ({{.GetDivSecondCoord}}) -- ({{.GetDivThirdCoord}});
`

const latexAnswerCode = `        \node [rounded corners, rectangle, minimum width={{.GetDivWidth}}*\zerowidth, minimum height = \zeroheight+\baselineskip, draw, below=0.15 cm of label3] {};
`

const latexDivOperandCode = `        \node [right=0.0 cm of {{.GetDivRef}}] ({{.GetDivLabel}}) {\huge {{.GetDivValue}}};
`

// types
// ----------------------------------------------------------------------------

// A coordinate is defined just with a label. From this basic
// definition it is then possible to specify either coordinates
// explicitly, with a pair (x, y) or with a formula which is used to
// compute the final location of the coordinate
type coordinate struct {
	label string
}

// explicit coordinates are defined with a pair (x,y) defined as
// floating-point numbers
type coordinateExplicit struct {
	x, y float64
	coordinate
}

// coordinates computed with a formula are characterized by a string
// which contains the formula tu use for computing the final location
// of the coordinate
type coordinateFormula struct {
	formula string
	coordinate
}

// A bounding box is defined with a couple of coordinates
// (bottom)--(right) which define the lower-left and upper-right
// corners
type boundingBox struct {

	// both coordinates are computed using formulas
	bottom, right coordinateFormula
}

// the box surrounding the dividend is consists of a path
// drawn between three coordinates whose location is
// determined using formulas
type splitBox struct {
	coord1, coord2, coord3 coordinateFormula
}

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

	// the first label is computed explicitly whereas the next two
	// labels are computed with respect to the previous ones
	label1         coordinateExplicit
	label2, label3 coordinateFormula

	// likewise the coordinate for determining the location of the
	// first line is computed with respect to another coordinate
	// and hence a formula is used
	line1 coordinateFormula

	// the bounding box surrounding all the necessary area for
	// solving the exercise is defined next
	bBox boundingBox

	// the box surrounding the dividend is consists of a path
	// drawn between three coordinates whose location is
	// determined using formulas
	sBox splitBox

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
	for nbdigits(quotient) != div.nbqdigits || quotient == 0 {
		dividend = randN(div.nbdvdigits)
		divisor = randN(div.nbdrdigits)
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

// -- coordinates

// Return the label of this coordinate
func (coord coordinate) GetDivLabel() string {
	return fmt.Sprintf("(%v)", coord.label)
}

func (coord coordinateExplicit) GetDivCoords() string {
	return fmt.Sprintf("(%v, %v)", coord.x, coord.y)
}

func (coord coordinateFormula) GetDivFormula() string {
	return fmt.Sprintf("%v", coord.formula)
}

// Provide TikZ code to represent an explicit coordinate fully
func (coord coordinateExplicit) String() string {

	// create a template with the TikZ code for showing indices
	tpl, err := template.New("coordinate").Parse(latexCoordinateExplicitCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, coord); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Return the pair (x, y) of the current coordinate
func (coord coordinateExplicit) GetCoords() string {

	return fmt.Sprintf("(%v, %v)", coord.x, coord.y)
}

// -- bounding Box

// Get the coordinate at the bottom-left corner of the bounding box
func (bbox boundingBox) GetDivBottom() string {

	// use the template defined for creating formula coords
	tpl, err := template.New("boundingBox").Parse(latexCoordinateFormulaCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, bbox.bottom); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// Get the coordinate at the top-right corner of the bounding box
func (bbox boundingBox) GetDivRight() string {

	// use the template defined for creating formula coords
	tpl, err := template.New("boundingBox").Parse(latexCoordinateFormulaCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, bbox.right); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

func (bbox boundingBox) String() string {

	// use the template defined for creating bounding boxes
	tpl, err := template.New("boundingBox").Parse(latexBoundingBoxCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, bbox); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// -- split box

// Return the TikZ code that represents the first coordinate of the split box
func (sbox splitBox) GetDivFirstCoord() string {
	return fmt.Sprintf("%v", sbox.coord1.formula)
}

// Return the TikZ code that represents the second coordinate of the split box
func (sbox splitBox) GetDivSecondCoord() string {
	return fmt.Sprintf("%v", sbox.coord2.formula)
}

// Return the TikZ code that represents the third coordinate of the split box
func (sbox splitBox) GetDivThirdCoord() string {
	return fmt.Sprintf("%v", sbox.coord3.formula)
}

// Return the TikZ code that draws the split box
func (sbox splitBox) String() string {

	// use the template defined for creating explicit coords
	tpl, err := template.New("division").Parse(latexSplitBoxCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, sbox); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()

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

	// use the template defined for creating split boxes
	tpl, err := template.New("division").Parse(latexCoordinateExplicitCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, division.label1); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// Generates the TikZ code necessary for positioning other coordinates
// after the first one
func (division divisionProblem) GetDivNextLabels() string {

	// use the template defined for creating explicit coords
	tpl, err := template.New("division").Parse(latexCoordinateFormulaCode)
	if err != nil {
		log.Fatal(err)
	}

	// now, use the same template to generate the code of the
	// second and third coordinate. Note that the first execution
	// actually writes to the output writter, which is then
	// buffered when the second execution is invoked, so that
	// reading the resulting string only once works correctly
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, division.label2); err != nil {
		log.Fatal(err)
	}
	if err := tpl.Execute(&tplOutput, division.label3); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// Generates the TikZ code necessary for positioning the first line of results
func (division divisionProblem) GetDivLine() string {

	// use the template defined for creating coords with a formula
	tpl, err := template.New("division").Parse(latexCoordinateFormulaCode)
	if err != nil {
		log.Fatal(err)
	}

	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, division.line1); err != nil {
		log.Fatal(err)
	}
	return tplOutput.String()
}

// Generates the TikZ code necessary for positioning the bounding box
// for solving the whole exercise
func (division divisionProblem) GetDivBoundingBox() string {

	// bounding boxes draw themselves
	return division.bBox.String()
}

// Generates the TikZ code necessary for drawing the split box
func (division divisionProblem) GetDivSplitBox() string {

	// split boxes draw themselves
	return division.sBox.String()
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

// Execute the given division problem and returns legal TikZ code to
// represent it
func (division divisionProblem) execute() string {

	// create a template with the TikZ code for showing this
	// division problem
	tpl, err := template.New("division").Parse(latexPreviousDivisionCode)

	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, division); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

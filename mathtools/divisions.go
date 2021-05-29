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
{{.Label1}}
{{.Label2}}
{{.Label3}}
        % -----------------------------------------------------------------------

        % --- Ancilliary reference points
{{.Line1}}
        % -----------------------------------------------------------------------

        % --- Bounding Box ------------------------------------------------------
{{.BBox}}
        % -----------------------------------------------------------------------
        % show the box enclosing the divisor
{{.SBox}}
        % show the box for writing the quotient
{{.Answer}}
        % -----------------------------------------------------------------------
        
        % --- Text ------------------------------------------------------------

        % Dividend
{{.Dividend}}
        % Divisor
{{.Divisor}}
        % -----------------------------------------------------------------------
`

// types
// ----------------------------------------------------------------------------

// The formal definition of a division problem is given below. It is defined
// with the number of digits of the dividend, divisor and quotient
type division struct {
	nbdvdigits int
	nbdrdigits int
	nbqdigits  int
}

// A division is characterized by its coordinates, a bounding box surrounding
// all the available area for solving the exercise, the box enclosing the
// divisor, and also the operands
type divisionTikZ struct {

	// the first label is computed explicitly whereas the next two labels are
	// computed with respect to the previous ones using formulas. All of them
	// are implemented using the reusable components of TikZ coordinates
	Label1         components.Coordinate
	Label2, Label3 components.Coordinate

	// likewise the coordinate for determining the location of the first line is
	// computed with respect to another coordinate and hence a formula is used.
	// There might be an arbitrary number of points but other points are
	// computed only wrt the location of the first line
	Line1 components.Coordinate

	// the bounding box surrounding all the necessary area for solving the
	// exercise is defined next and it is computed with two formulas that
	// specify the lower left and upper right corners
	BBox components.CoordinatedRectangle

	// the box surrounding the dividend consists of a path drawn between three
	// coordinates whose location is determined using formulas
	SBox components.Line

	// the answer should be written within a box explicitly shown
	Answer components.Text

	// finally, both operands, are created next and implemented as Texts
	Dividend, Divisor components.Text
}

// methods
// ----------------------------------------------------------------------------

// -- divisionTikZ

// Return the LaTeX/TikZ commands that show up the picture stored in the
// receiver
func (tikz divisionTikZ) execute() string {

	// create a template with the TikZ code for showing this picture
	tpl, err := template.New("divisionTikZ").Parse(tikZDivisionCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of
	// the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, tikz); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// -- division

// return the instance of a specific division problem that can be marshalled in
// JSON format. The receiver is assumed to have been fully verified so that it
// should be consistent.
//
// The result is given with four items: dividend, divisor, quotient and
// remainer. The remainder and the quotient are shown as "?" in the arguments as
// they have to be guessed by the student
func (div division) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

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

// return a valid LaTeX/TikZ representation of this sequence using TikZ
// components
func (div division) GetTikZPicture() string {

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

	// note the answer is written withing a text box which necessarily contains
	// nothing. No label is assigned to it as well as no computations are
	// performed from its location
	answer := components.NewText(
		fmt.Sprintf(`rounded corners, rectangle, minimum width=%v*\zerowidth, minimum height = \zeroheight+\baselineskip, draw, below=0.15 cm of label3`,
			2.0+helpers.Max(float64(div.nbdrdigits), float64(div.nbqdigits))),
		"", "",
	)

	// -- operands

	// randomly determine the values of the operands. For this, the service that
	// generates problems is the one that can marshal them into JSON format. The
	// dividend is returned in the first position and the divisor in the second
	instance, err := div.generateJSONProblem()
	if err != nil {
		log.Fatalf(" Fatal error while generating a valid division: %v", err)
	}

	dividend := components.NewText(
		`right=0.0 cm of label1`,
		"dividend",
		`\huge `+instance.Solution[0],
	)
	divisor := components.NewText(
		`right=0.0 cm of label2`,
		"divisor",
		`\huge `+instance.Solution[1],
	)

	// And put all this elements together to show up the picture of a division
	divPicture := divisionTikZ{
		Label1:   label1,
		Label2:   label2,
		Label3:   label3,
		Line1:    line1,
		BBox:     bBox,
		SBox:     sBox,
		Answer:   answer,
		Dividend: dividend,
		Divisor:  divisor,
	}

	// and return the TikZ code necessary for drawing the problem
	return divPicture.execute()
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

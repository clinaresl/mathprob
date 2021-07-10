// -*- coding: utf-8 -*-
// basicoperation.go
//
// Description: Provides services for automatically creating a basic operation
// -----------------------------------------------------------------------------
//
// Started on <mar 25-05-2021 20:47:25.993610806 (1621968445)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package mathtools

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"text/template"
	"time"

	"github.com/clinaresl/mathprob/helpers"
	"github.com/clinaresl/mathprob/mathtools/components"
)

// constants
// ----------------------------------------------------------------------------

// There are two different types of basic operations: "result" or "operand". In
// the first case, all operands are visible and the student has to provide the
// value of the result; in the latter, the result can be seen but one operand is
// missing whose value has to be guessed by the student
const (
	BORESULT int = iota
	BOOPERAND
)

// the TikZ code for generating arbitrary basic operations is shown next. Note that it
// makes use of LaTeX/TikZ components
const latexBasicOperationCode = `\begin{minipage}{0.25\linewidth}
    \begin{center}
        \begin{tikzpicture}

            % draw the basic operation
            {{.GetTikZPicture}}

        \end{tikzpicture}
    \end{center}
\end{minipage}
`

const tikZBasicOperationCode = `% --- Coordinates -----------------------------------------------------

      % Lower-left corner of the bounding box
      {{.Bottom}}

      % the result is located leaving some room to the let so that operations
      % can be drawn next to others withouth colliding. For this, the result
      % is x-shifted 1 plus half the number of digits of the result. It is
      % also always y-shifted 1.5 the baselineskip plus half the height of a
      % digit
      {{.Answer}}

      % --- Split line ------------------------------------------------------

      % Next, a line splitting the operands and result is shown
      {{.Split1}}
      {{.Split2}}
      {{.SplitLine}}

      % --- Operands --------------------------------------------------------

      % next, all operands are shown, as this is a type 0, they are not within
      % a box
      {{.GetOperands}}

      % --- Operator --------------------------------------------------------

      % the operator is shown to the left of the first (lower) operand
      {{.OperatorCoord}}
      {{.Operator}}

      % ---------------------------------------------------------------------

      % --- Bounding Box ----------------------------------------------------

      % the distance between the answer box and the end of the bounding box is
      % half the width of the bounding box. As this bounding box contains one
      % digit its width is 3.0 and hence 1.5 has to be multiplied by the
      % width of zero
      {{.BBox}}

      % ---------------------------------------------------------------------

      % --- Answer box ------------------------------------------------------

      {{.Result}}

      % ---------------------------------------------------------------------
`

// types
// ----------------------------------------------------------------------------

// A basic operation consists of a number of operands related to any of the
// operations: +, -, *, / whose number of digits have to be specified as much as
// the number of desired digits in the result. There are two types of basic
// operations:
//
//    0: all operands are given and the student has to guess the result
//    1: all operands but one are shown but the result can be seen. The student
//    has to provide the value of the missing operand
type basicOperation struct {
	botype       int
	operator     string
	nboperands   int
	nbdigitsop   int
	nbdigitsrslt int
}

// The following struct stores all the information necessary to draw basic
// operations
type basicOperationTikZ struct {

	// the lower left corner of the bounding box is located always at (0, 0)
	Bottom components.Coordinate

	// The answer box is centered at the coordinate answer
	Answer components.Coordinate

	// The split line is drawn between two extremes split1 and split2
	Split1, Split2 components.Coordinate
	SplitLine      components.Line

	// As there can be an arbitrary number of operands, these are stored in a
	// slice along with its coordinates
	coords []components.Coordinate
	ops    []components.LabeledText

	// The operator is located to the left of the first operand ---the one
	// immediately above the split line
	OperatorCoord components.Coordinate
	Operator      components.LabeledText

	// the bounding box surrounding all the necessary area for solving the
	// exercise is defined next and it is computed with two formulas that
	// specify the lower left and upper right corners
	BBox components.CoordinatedRectangle

	// The result is shown in the position stored in the label answer
	Result components.LabeledText
}

// methods
// ----------------------------------------------------------------------------

// -- basicOperationTikZ

// Generates the TikZ code necessary for positioning all operands of the basic
// operation, either empty cells or hints
func (tikz basicOperationTikZ) GetOperands() string {

	// Use a btyes buffer to append the strings of each operand
	var output bytes.Buffer

	// First, add all coordinates
	for _, coord := range tikz.coords {
		fmt.Fprintf(&output, "%v\n", coord)
	}

	// Draw all text boxes in the operands stored in this pict
	for _, op := range tikz.ops {
		fmt.Fprintf(&output, "%v\n", op)
	}

	// and return the concatenation of the LaTeX/TikZ code used for drawing all
	// operands
	return output.String()
}

// Return the LaTeX/TikZ commands that show up the picture stored in the
// receiver
func (tikz basicOperationTikZ) execute() string {

	// create a template with the TikZ code for showing this picture
	tpl, err := template.New("basicOperationTikZ").Parse(tikZBasicOperationCode)
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

// -- basicOperation

// return the instance of a specific basic operation problem that can be
// marshalled in JSON format. The receiver is assumed to have been fully
// verified so that it should be consistent.
//
// The result is given as an array of numbers:
//    1. The first string is the operation to perform: "+", "-", "*" or "/"
//    2. First, all operands are given
//    3. The last string is the result
func (bo basicOperation) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// first, ensure that the number of digits both for the operands and the
	// result are compatible
	switch bo.operator {
	case "+":

		// no math expression! I just compute the upper and lower bound on the
		// number of digits in the result and compare it to the value given
		if helpers.NbDigits(bo.nboperands*int(math.Pow(10, float64(bo.nbdigitsop))-1)) < bo.nbdigitsrslt ||
			helpers.NbDigits(bo.nboperands*int(math.Pow(10, float64(bo.nbdigitsop-1)))) > bo.nbdigitsrslt {
			return problemJSON{}, fmt.Errorf("It is not possible to generate summations with %v digits using %v operands with %v digits each",
				bo.nbdigitsrslt, bo.nboperands, bo.nbdigitsop)
		}

	case "-":

		// watch out! the possibility of generating negative numbers is also
		// considered. Thus, the resulting guard is pretty close to the previous
		// one when computing the maximum, though the number of operands minus
		// one is used instead because the largest number (in magnitude) can be
		// generated if and only if the first one is zero, so that the maximum
		// number of digits in the result is the same as if we are summing up
		// all operands but the first one. As for the lower bound in the number
		// of digits it is clearly one
		if helpers.NbDigits((bo.nboperands-1)*int(math.Pow(10, float64(1+bo.nbdigitsop)))) < bo.nbdigitsrslt ||
			1 > bo.nbdigitsrslt {
			return problemJSON{}, fmt.Errorf("It is not possible to generate subtractions with %v digits using %v operands with %v digits each",
				bo.nbdigitsrslt, bo.nboperands, bo.nbdigitsop)
		}

	case "*":

		// this is easy ...
		if bo.nboperands*bo.nbdigitsop < bo.nbdigitsrslt ||
			1+bo.nboperands*(bo.nbdigitsop-1) > bo.nbdigitsrslt {
			return problemJSON{}, fmt.Errorf("It is not possible to generate multiplications with %v digits using %v operands with %v digits each",
				bo.nbdigitsrslt, bo.nboperands, bo.nbdigitsop)
		}

	case "/":

		// Divisions can consist only of two arguments
		if bo.nboperands > 2 {
			return problemJSON{}, errors.New("Divisions can consist only of two items!")
		}

		// and considering that both operands have the same number of digits,
		// the result necessarily consists of one single digit
		if bo.nbdigitsrslt != 1 {
			return problemJSON{}, errors.New("Divisions can only generate results with 1 digit")
		}
	}

	// in case type 1 was selected, randomly choose any location among all
	// operands
	pos := 1 + rand.Int()%bo.nboperands

	// next, create the instance.
	var result int
	solution := make([]string, 2+bo.nboperands)

	// The first position of the solution slice is the operation to perform
	solution[0] = bo.operator

	// Create first a solution which will be masked later on. The result is
	// stored in the last location of the slice. Note that the vector directly
	// contains strings. This is done because it will then be inserted in the
	// struct to return even if it forces casting types all the time

	// and now randomly generate operands of the given width until a result of
	// the desired width is generated. Also, basic operations are intended for
	// very beginners and thus, negative values are intentionally removed
	for helpers.NbDigits(result) != bo.nbdigitsrslt ||
		result <= 0 {

		// generate all operands first and write them tentatively in the
		// solution slice
		for i := 0; i < bo.nboperands; i++ {
			solution[1+i] = fmt.Sprintf("%v", helpers.RandN(bo.nbdigitsop))
		}

		// compute the specified operation over these items. First initialize
		// the result to the value of the first operand
		result, _ = helpers.Atoi(solution[1])
		for i := 1; i < bo.nboperands; i++ {
			value, _ := helpers.Atoi(solution[1+i])
			switch bo.operator {
			case "+":
				result += value
			case "-":
				result -= value
			case "*":
				result *= value
			case "/":
				result /= value
			}
		}
	}
	solution[1+bo.nboperands] = fmt.Sprintf("%v", result)

	// now, copy the solution to the args but ...
	args := make([]string, 2+bo.nboperands)
	for i := 0; i < 2+bo.nboperands; i++ {
		args[i] = fmt.Sprintf("%v", solution[i])
	}

	// ... replace the location pos in case this is a type 1 problem
	if bo.botype == BOOPERAND {
		args[pos] = "?"
	} else {

		// otherwise, mask the result of the basic operation
		args[1+bo.nboperands] = "?"
	}

	return problemJSON{
		Probtype: "BasicOperation",
		Args:     args,
		Solution: solution,
	}, nil
}

// return a valid LaTeX/TikZ representation of this basic operation using TikZ
// components
func (bo basicOperation) GetTikZPicture() string {

	// -- operands: randomly determine the values of the operands. For this, the
	//              service that generates problems is the one that can marshal
	//              them into JSON format. The operands and the result are given
	//              in Args, where a question mark is a number that has to be
	//              guessed by the student
	instance, err := bo.generateJSONProblem()
	if err != nil {
		log.Fatalf(" Fatal error while generating a valid basic operation: %v", err)
	}

	// compute the number of digits required to draw all operands and the result
	nbdigits := helpers.Max(float64(bo.nbdigitsop), float64(bo.nbdigitsrslt))

	// -- Coordinates

	// Bottom is the lower-left corner of the bounding box
	bottom := components.NewCoordinate(components.Point{
		X: 0.0,
		Y: 0.0,
	}, "bottom")

	// The answer box is located in the last row of the figure
	answer := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(bottom) + (%v\zerowidth, 0.5\zeroheight+1.0\baselineskip)$`,
			1.5+(2.0+nbdigits)/2.0)),
		"answer",
	)

	// The split line is drawn between two endpoints whose coordinates are
	// computed separately
	split1 := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(answer) + (%v\zerowidth, 1.5\baselineskip)$`,
			-(0.75+(2.0+nbdigits)/2.0))),
		"split1",
	)
	split2 := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(answer) + (%v\zerowidth, 1.5\baselineskip)$`,
			(2.0+nbdigits)/2.0)),
		"split2",
	)
	splitLine := components.NewLine("split1", "split2")
	splitLine.SetOptions("thick")

	// -- operands

	// the operands to draw are given in the Args field of this specific
	// instance. If a question mark is given, then it should be replaced with an
	// empty text box; otherwise, the specified number is shown. Note that the
	// operands are the numbers in the arguments of this instance but the first
	// (which is the operand) and the last ---which is the result
	var coords []components.Coordinate
	var ops []components.LabeledText
	for idx, item := range instance.Args[1 : len(instance.Args)-1] {

		// Create two ancilliary variables to store the coordinate and aspect of
		// the next cell
		var box components.LabeledText

		// in spite of the contents, the next cell is located at the following
		// coordinate. Note that the name of the coordinate is a little bit
		// weird. This is because in the LaTeX manual for basic operations, op1
		// is the one right immediately above the split line
		ith := float64(len(instance.Args)-idx) - 2.0
		coord := components.NewCoordinate(
			components.Formula(fmt.Sprintf(`$(answer) + (0, %v\zeroheight + %v\baselineskip)$`,
				ith-1.0,
				2.0+ith)),
			fmt.Sprintf("op%v", ith),
		)

		// if this is a question mark
		if item == "?" {

			// then add an empty text box
			box = components.NewLabeledText(
				fmt.Sprintf(`rounded corners, rectangle, minimum width=%v*\zerowidth, minimum height = \zeroheight + \baselineskip, draw`,
					2.0+nbdigits,
				),
				fmt.Sprintf("op%v", ith),
				"",
			)
		} else {

			// otherwise, add the number itself
			box = components.NewLabeledText(
				"",
				fmt.Sprintf("op%v", ith),
				`\huge `+item)
		}

		// and add the new box and its coordinates
		coords = append(coords, coord)
		ops = append(ops, box)
	}

	// -- operator
	operatorCoord := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(op1) + (%v\zerowidth, 0.0)$`,
			-0.75-(2+nbdigits)/2.0)),
		"operator",
	)

	// the text to show for the operator depends upon the operator requested
	var opLaTeX string
	switch instance.Args[0] {
	case "+", "-":
		opLaTeX = instance.Args[0]
	case "*":
		opLaTeX = `$\times$`
	case "/":
		opLaTeX = `$\div$`
	}
	operator := components.NewLabeledText("", "operator", `\huge `+opLaTeX)

	// -- bounding box
	right := components.NewCoordinate(
		components.Formula(fmt.Sprintf(`$(split2) + (0.75\zerowidth, %v\baselineskip)$`,
			1+2.0*(len(instance.Args)-2))),
		"right")
	bBox := components.NewCoordinatedRectangle(bottom, right)
	bBox.SetOptions("white")

	// -- result

	// Note that the result might be either unknown (basic operations type 0) or
	// known (basic operations type 1)
	var result components.LabeledText
	if instance.Args[len(instance.Args)-1] == "?" {

		// in case it is unknown, draw an empty box
		result = components.NewLabeledText(
			fmt.Sprintf(`rounded corners, rectangle, minimum width=%v*\zerowidth, minimum height = \zeroheight + \baselineskip, draw`,
				2.0+nbdigits,
			),
			fmt.Sprintf("answer"),
			"",
		)
	} else {

		// otherwise, show the number in the arguments of this instance
		result = components.NewLabeledText(
			"",
			fmt.Sprintf("answer"),
			`\huge `+instance.Args[len(instance.Args)-1])
	}

	// And put all these elements together to show up the picture of a basic
	// operation
	boPicture := basicOperationTikZ{
		Bottom:        bottom,
		Answer:        answer,
		Split1:        split1,
		Split2:        split2,
		SplitLine:     splitLine,
		coords:        coords,
		ops:           ops,
		OperatorCoord: operatorCoord,
		Operator:      operator,
		BBox:          bBox,
		Result:        result,
	}

	// and return the TikZ code necessary for drawing the problem
	return boPicture.execute()
}

// Return TikZ code that represents a basic operation
func (bo basicOperation) execute() string {

	// create a template with the TikZ code for showing this basic operation
	tpl, err := template.New("basicOperation").Parse(latexBasicOperationCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, bo); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// Local Variables:
// mode:go
// fill-column:80
// End:

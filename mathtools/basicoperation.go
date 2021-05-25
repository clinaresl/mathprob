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

// the TikZ code for generating arbitrary sequences is shown next. Note that it
// makes use of LaTeX/TikZ components
const latexBasicOperationCode = `\begin{minipage}{0.25\linewidth}
    \begin{center}
        \begin{tikzpicture}

            % draw the sequence
            {{.GetTikZPicture}}

        \end{tikzpicture}
    \end{center}
\end{minipage}
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

// methods
// ----------------------------------------------------------------------------

// -- basicOperation

// return the instance of a specific basic operation problem that can be
// marshalled in JSON format. The receiver is assumed to have been fully
// verified so that it should be consistent.
//
// The result is given
func (bo basicOperation) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// first, ensure that the number of digits both for the operands and the
	// result are compatible
	switch bo.operator {
	case "+":

		// no math expression! I just compute the upper and lower bound on the
		// number of digits in the result and compare it to the value given
		if helpers.NbDigits(bo.nboperands*int(math.Pow(10, float64(1+bo.nbdigitsop)))) < bo.nbdigitsrslt ||
			helpers.NbDigits(bo.nboperands*int(math.Pow(10, float64(bo.nbdigitsop)))) > bo.nbdigitsrslt {
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
	pos := rand.Int() % bo.nboperands

	// next, create the instance. Create first a solution which will be masked
	// later on. The result is stored in the last location of the slice. Note
	// that the vector directly contains strings. This is done because it will
	// then be inserted in the struct to return even if it forces casting types
	// all the time
	var result int
	solution := make([]string, 1+bo.nboperands)

	// and now randomly generate operands of the given width until a result of
	// the desired width is generated. Also, basic operations are intended for
	// very beginners and thus, negative values are intentionally removed
	for helpers.NbDigits(result) != bo.nbdigitsrslt ||
		result <= 0 {

		// generate all operands first and write them tentatively in the
		// solution slice
		for i := 0; i < bo.nboperands; i++ {
			solution[i] = fmt.Sprintf("%v", helpers.RandN(bo.nbdigitsop))
		}

		// compute the specified operation over these items. First initialize
		// the result to the value of the first operand
		result, _ = helpers.Atoi(solution[0])
		for i := 1; i < bo.nboperands; i++ {
			value, _ := helpers.Atoi(solution[i])
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
	solution[bo.nboperands] = fmt.Sprintf("%v", result)

	// now, copy the solution to the args but ...
	args := make([]string, 1+bo.nboperands)
	for i := 0; i < 1+bo.nboperands; i++ {
		args[i] = fmt.Sprintf("%v", solution[i])
	}

	// ... replace the location pos in case this is a type 1 problem
	if bo.botype == BOOPERAND {
		args[pos] = "?"
	} else {

		// otherwise, mask the result of the basic operation
		args[bo.nboperands] = "?"
	}

	return problemJSON{
		Probtype: "BasicOperation",
		Args:     args,
		Solution: solution,
	}, nil
}

// return a valid LaTeX/TikZ representation of this sequence using TikZ
// components
func (bo basicOperation) GetTikZPicture() string {
	return "Here I come!!"
}

// Return TikZ code that represents a sequence
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

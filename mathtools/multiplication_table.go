// -*- coding: utf-8 -*-
// multiplication_table.go
//
// Description: Provides services for automatically generating multiplication
// tables
// -----------------------------------------------------------------------------
//
// Started on <lun 31-05-2021 20:00:08.961920795 (1622484008)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package mathtools

import (
	"bytes"
	"fmt"
	"log"
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
	MTRESULT int = iota
	MTOPERAND
)

// the TikZ code for generating arbitrary multiplication tables is shown next.
// Note that it makes use of LaTeX/TikZ components
const latexMultiplicationTableCode = `\begin{minipage}{\linewidth}
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

// A multiplication table consists of 10 different rows where the same number is
// multiplied by geq, geq + 1, ... leq. If not given, the values of geq and leq
// should be 1 and 10 by default. The definition of a multiplication table
// includes the number of digits used in the factor repeated in all rows. Other
// options are whether the factors are presented in the usual order: 5x1
// (inv=false) or 1x5 (inv=true) and also whether the rows are sorted or not
// ---field sort. If inv is enabled it is randomly chosen whether each row is
// shown in the regular or inversed order

// In addition, there are two different types of multiplication tables:
//
//    0: both operands are given and the student has to guess the result
//    1: only one operand is given, and the student has to guess the value of
//    the other operand so that the equality holds
type multiplicationTable struct {
	mttype   int
	nbdigits int
	geq, leq int
	inv      bool
	sorted   bool
}

// methods
// ----------------------------------------------------------------------------

// -- multiplicationTable

// return the instance of a specific basic multiplication table that can be
// marshalled in JSON format. The receiver is assumed to have been fully
// verified so that it should be consistent.
//
// The result is given as an array of numbers:
//    1. The first string is the number used in the multiplication table
//    2. Next, all items of each row are given in sorted order, e.g., "5", "1",
//    "5" which stands for "5x1=5". If one item has to be guessed it is shown as
//    a question mark "?"
func (mt multiplicationTable) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// first, determine the factor to use in all rows of the multiplication
	// table
	factor := helpers.RandN(mt.nbdigits)

	// now, make room to store the full solution of the multiplication table. In
	// total (1+leq-geq) rows have to be generated, each with three digits and
	// write down the number used in the multiplication table
	solution := make([]string, 1+(1+mt.leq-mt.geq)*3)
	solution[0] = fmt.Sprintf("%v", factor)

	// fill in the table
	for i := mt.geq; i <= mt.leq; i++ {

		// compute the relative position of this number
		idx := i - mt.geq

		// store the values in the solution with the usual order
		solution[1+idx*3] = fmt.Sprintf("%v", factor)
		solution[2+idx*3] = fmt.Sprintf("%v", i)

		// if the inverted presentation of factors was requested then randomly
		// determine whether to write the factor first or later
		if mt.inv {

			// if a number randomly generated in the interval [0, 100) falls in
			// the first half, then reverse the operands
			if rand.Int()%100 < 50 {
				solution[1+idx*3], solution[2+idx*3] = solution[2+idx*3], solution[1+idx*3]
			}
		}

		// in both cases, write down the result next
		solution[3+idx*3] = fmt.Sprintf("%v", factor*i)
	}

	// In case sorted is false, then shuffle all items in the multiplication
	// table
	if !mt.sorted {

		// For this, shuffle a slice of ints with the indexes of each row
		identity := make([]int, 10)
		for i := 0; i < 10; i++ {
			identity[i] = i
		}

		// and now shuffle them
		rand.Shuffle(len(identity),
			func(i, j int) {
				identity[i], identity[j] = identity[j], identity[i]
			})

		// now, affect the order of the solution as specified in the shuffled
		// slice. Note that as this is a destructive operation over solution, a
		// copy is necessary
		isolution := make([]string, len(solution))
		copy(isolution, solution)
		for i := 0; i < 10; i++ {
			solution[1+i*3], solution[2+i*3], solution[3+i*3] =
				isolution[1+identity[i]*3], isolution[2+identity[i]*3], isolution[3+identity[i]*3]
		}
	}

	// once the multiplication table has been fully generated, it is now the
	// turn to create the specific instance determining what numbers are hidden.
	// Note that the arguments preserve the first value, the factor used in the
	// multiplication table
	args := make([]string, 1+(1+mt.leq-mt.geq)*3)
	args[0] = solution[0]
	for i := 0; i < 1+mt.leq-mt.geq; i++ {

		// in case this is an ordinary multiplication table, just create the
		// instance as usual
		if mt.mttype == MTRESULT {
			args[1+i*3], args[2+i*3] = solution[1+i*3], solution[2+i*3]
			args[3+i*3] = "?"
		} else {

			// randomly determine what operand to mask. Tentaively, preserve the
			// first operand and the result
			args[1+i*3], args[3+i*3] = solution[1+i*3], solution[3+i*3]
			args[2+i*3] = "?"

			// if a number randomly generated in the interval [0, 100) falls in
			// the first half then mask the first operand instead
			if rand.Int()%100 < 50 {

				args[2+i*3], args[3+i*3] = solution[2+i*3], solution[3+i*3]
				args[1+i*3] = "?"
			}
		}
	}

	// Now, generate the multiplication table
	return problemJSON{
		Probtype: "MultiplicationTable",
		Args:     args,
		Solution: solution,
	}, nil
}

// return a valid LaTeX/TikZ representation of this multiplication table using
// TikZ components
func (mt multiplicationTable) GetTikZPicture() string {
	return "Hi there!"
}

// Return TikZ code that represents a sequence
func (mt multiplicationTable) execute() string {

	// create a template with the TikZ code for showing this basic operation
	tpl, err := template.New("basicOperation").Parse(latexMultiplicationTableCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the execution of the
	// template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, mt); err != nil {
		log.Fatal(err)
	}

	// and return the resulting string
	return tplOutput.String()
}

// Local Variables:
// mode:go
// fill-column:80
// End:

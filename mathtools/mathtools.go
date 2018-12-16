/*
  mathtools.go
  Description: Different tools for handling mathematical problems
  -----------------------------------------------------------------------------

  Started on  <Mon Jul 10 19:29:17 2017 >
  Last update <>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by
  Login   <clinares@atlas>
*/

// Provides various services for randomly generating math problems
package mathtools

import (
	"fmt"
	"log" // logging services
	"math"
	"math/rand"
	"os" // access to file mgmt functions
	"time"

	"text/template" // go facility for processing templates

	"bitbucket.org/mathprob/fstools"
)

// types
// ----------------------------------------------------------------------------

// Positions are used by virtually any operation to create. Thus, it
// is provided here. Every position is qualified by its x- and y-
// coordinates
type position struct {
	x, y float64
}

// A master file consists of an input filename that stores the
// tempalte to fill in to generate the final sheet of exercises, and
// an output tex filename. It also comes with other fields that can be
// used for customizing the resulting file such as the student's name
type MasterFile struct {
	Infile  string
	Name    string
	Outfile string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a master file with the given name
func NewMasterFile(filename, name string) MasterFile {

	return MasterFile{Infile: filename, Name: name}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// methods
// ----------------------------------------------------------------------------

// -- position
// ----------------------------------------------------------------------------

// String provides a human readable/TikZ digestible forma of the
// contents of any position
func (pos position) String() string {
	return fmt.Sprintf("(%v, %v)", pos.x, pos.y)
}

// -- MasterFile
// ----------------------------------------------------------------------------

// Return the input filename that shall store the template file to
// generate the exercises
func (masterFile MasterFile) GetInfile() string {
	return masterFile.Infile
}

// Return the student's name of this master file
func (masterFile MasterFile) GetName() string {
	return masterFile.Name
}

// Return the output tex filename that shall contain the exercises in tex
func (masterFile MasterFile) GetOutfile() string {
	return masterFile.Outfile
}

// the following function is provided just to allow the text/template to repeat
// the same statement an arbitrary number of times. It just return a slice of
// MasterFiles of a given length. Each element can then be used to invoke the
// various services provided for text/templates
func (masterFile MasterFile) Slice(n int) []MasterFile {
	return make([]MasterFile, n)
}

// Get the LaTeX code for a simple operation where both operands must
// be lower than 10
func (masterFile MasterFile) GetSimpleOperation11(operator int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// create the two operands
	latexOperandA := latexOperand{
		label: "label1",
		id:    "num1",
		pos:   position{x: 2.50, y: 3.50},
		value: 1 + rand.Intn(9),
	}
	latexOperandB := latexOperand{
		label: "label2",
		id:    "num2",
		pos:   position{x: 2.50, y: 2.50},
		value: 1 + rand.Intn(9),
	}

	// create the operator and, in passing, verify specific conditions for
	// the operands to make sense with the requested operator
	latexOperator := latexOperator{
		label: "op",
		pos:   position{x: 1.00, y: 2.50},
	}
	switch operator {
	case ADD:
		latexOperator.symbol = `$+$`
	case SUB:
		latexOperator.symbol = `$-$`

		// make sure that operandB is strictly less or equal than
		// operandA
		if latexOperandA.value < latexOperandB.value {
			latexOperandA.value, latexOperandB.value = latexOperandB.value, latexOperandA.value
		}

	case PROD:
		latexOperator.symbol = `$\times$`
	case DIV:
		latexOperator.symbol = `$\div$`

		// ensure the second operand is not 0
		if latexOperandB.value == 0 {
			latexOperandB.value = 1
		}

		// ensure the result is an integer value
		if latexOperandA.value%latexOperandB.value != 0 {

			// just divide by the minimum of both operands and
			// randomly select the larger number to be below 10
			latexOperandB.value = min(latexOperandA.value, latexOperandB.value)
			latexOperandA.value = int(float64(latexOperandB.value) * math.Floor((math.Log(float64(10)))/(math.Log(float64(latexOperandB.value)))))
		}
	default:
		log.Fatalf("Unknown operator '%v'", operator)
	}

	// create the result
	latexResult := latexResult{
		label:         "label3",
		pos:           position{x: 2.75, y: 1.00},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the operation
	latexOperation := singleOperation{
		operandA: latexOperandA,
		operandB: latexOperandB,
		operator: latexOperator,
		result:   latexResult,
	}

	return latexOperation.Execute()
}

// Get the LaTeX code for a simple operation where the first operand
// has to be less than 20 and the second one has to be less than 10,
// and the sum of both operands strictly less or equal than 20
func (masterFile MasterFile) GetSimpleOperation21Bounded20(operator int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// create the two operands
	latexOperandB := latexOperand{
		label: "label2",
		id:    "num2",
		pos:   position{x: 2.50, y: 2.50},
		value: 1 + rand.Intn(9),
	}
	latexOperandA := latexOperand{
		label: "label1",
		id:    "num1",
		pos:   position{x: 2.50, y: 3.50},
		value: 10 + rand.Intn(10-latexOperandB.value),
	}

	// create the operator and, in passing, verify specific conditions for
	// the operands to make sense with the requested operator
	latexOperator := latexOperator{
		label: "op",
		pos:   position{x: 1.00, y: 2.50},
	}
	switch operator {
	case ADD:
		latexOperator.symbol = `$+$`
	case SUB:
		latexOperator.symbol = `$-$`

		// make sure that operandB is strictly less or equal than
		// operandA
		if latexOperandA.value < latexOperandB.value {
			latexOperandA.value, latexOperandB.value = latexOperandB.value, latexOperandA.value
		}

	case PROD:
		latexOperator.symbol = `$\times$`
	case DIV:
		latexOperator.symbol = `$\div$`

		// ensure the second operand is not 0
		if latexOperandB.value == 0 {
			latexOperandB.value = 1
		}

		// ensure the result is an integer value
		if latexOperandA.value%latexOperandB.value != 0 {

			// just divide by the minimum of both operands and
			// randomly select the larger number to be below 10
			latexOperandB.value = min(latexOperandA.value, latexOperandB.value)
			latexOperandA.value = int(float64(latexOperandB.value) * math.Floor((math.Log(float64(10)))/(math.Log(float64(latexOperandB.value)))))
		}
	default:
		log.Fatalf("Unknown operator '%v'", operator)
	}

	// create the result
	latexResult := latexResult{
		label:         "label3",
		pos:           position{x: 2.75, y: 1.00},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the operation
	latexOperation := singleOperation{
		operandA: latexOperandA,
		operandB: latexOperandB,
		operator: latexOperator,
		result:   latexResult,
	}

	return latexOperation.Execute()
}

// Get the LaTeX code for a simple operation where the first operand
// has to be less than 100 and the second one has to be less than 10
func (masterFile MasterFile) GetSimpleOperation21(operator int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// create the two operands
	latexOperandA := latexOperand{
		label: "label1",
		id:    "num1",
		pos:   position{x: 2.50, y: 3.50},
		value: 1 + rand.Intn(99),
	}
	latexOperandB := latexOperand{
		label: "label2",
		id:    "num2",
		pos:   position{x: 2.50, y: 2.50},
		value: 1 + rand.Intn(9),
	}

	// create the operator and, in passing, verify specific conditions for
	// the operands to make sense with the requested operator
	latexOperator := latexOperator{
		label: "op",
		pos:   position{x: 1.00, y: 2.50},
	}
	switch operator {
	case ADD:
		latexOperator.symbol = `$+$`
	case SUB:
		latexOperator.symbol = `$-$`

		// make sure that operandB is strictly less or equal than
		// operandA
		if latexOperandA.value < latexOperandB.value {
			latexOperandA.value, latexOperandB.value = latexOperandB.value, latexOperandA.value
		}

	case PROD:
		latexOperator.symbol = `$\times$`
	case DIV:
		latexOperator.symbol = `$\div$`

		// ensure the second operand is not 0
		if latexOperandB.value == 0 {
			latexOperandB.value = 1
		}

		// ensure the result is an integer value
		if latexOperandA.value%latexOperandB.value != 0 {

			// just divide by the minimum of both operands and
			// randomly select the larger number to be below 10
			latexOperandB.value = min(latexOperandA.value, latexOperandB.value)
			latexOperandA.value = int(float64(latexOperandB.value) * math.Floor((math.Log(float64(10)))/(math.Log(float64(latexOperandB.value)))))
		}
	default:
		log.Fatalf("Unknown operator '%v'", operator)
	}

	// create the result
	latexResult := latexResult{
		label:         "label3",
		pos:           position{x: 2.75, y: 1.00},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the operation
	latexOperation := singleOperation{
		operandA: latexOperandA,
		operandB: latexOperandB,
		operator: latexOperator,
		result:   latexResult,
	}

	return latexOperation.Execute()
}

// Get the LaTeX code for a simple operation where both operands must
// be lower than 100
func (masterFile MasterFile) GetSimpleOperation22(operator int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// create the two operands
	latexOperandA := latexOperand{
		label: "label1",
		id:    "num1",
		pos:   position{x: 2.50, y: 3.50},
		value: 1 + rand.Intn(99),
	}
	latexOperandB := latexOperand{
		label: "label2",
		id:    "num2",
		pos:   position{x: 2.50, y: 2.50},
		value: 1 + rand.Intn(99),
	}

	// create the operator and, in passing, verify specific conditions for
	// the operands to make sense with the requested operator
	latexOperator := latexOperator{
		label: "op",
		pos:   position{x: 1.00, y: 2.50},
	}
	switch operator {
	case ADD:
		latexOperator.symbol = `$+$`
	case SUB:
		latexOperator.symbol = `$-$`

		// make sure that operandB is strictly less or equal than
		// operandA
		if latexOperandA.value < latexOperandB.value {
			latexOperandA.value, latexOperandB.value = latexOperandB.value, latexOperandA.value
		}

	case PROD:
		latexOperator.symbol = `$\times$`
	case DIV:
		latexOperator.symbol = `$\div$`

		// ensure the second operand is not 0
		if latexOperandB.value == 0 {
			latexOperandB.value = 1
		}

		// ensure the result is an integer value
		if latexOperandA.value%latexOperandB.value != 0 {

			// just divide by the minimum of both operands and
			// randomly select the larger number to be below 10
			latexOperandB.value = min(latexOperandA.value, latexOperandB.value)
			latexOperandA.value = int(float64(latexOperandB.value) * math.Floor((math.Log(float64(100)))/(math.Log(float64(latexOperandB.value)))))
		}
	default:
		log.Fatalf("Unknown operator '%v'", operator)
	}

	// create the result
	latexResult := latexResult{
		label:         "label3",
		pos:           position{x: 2.75, y: 1.00},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the operation
	latexOperation := singleOperation{
		operandA: latexOperandA,
		operandB: latexOperandB,
		operator: latexOperator,
		result:   latexResult,
	}

	return latexOperation.Execute()
}

// Get the LaTeX code for generating a problem on sequences of the
// given type: to either guess the previous (seqtype 1), subsequent
// (2), or both (3). Indices are strictly less than 10
func (masterFile MasterFile) GetSequence10(seqtype int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// decide first the x-position of the index and the box for
	// drawing the answer
	var xIndex, xAnswer float64
	switch seqtype {
	case PREVIOUS:
		xIndex, xAnswer = 2.00, 0.50
	case SUBSEQUENT:
		xIndex, xAnswer = 0.50, 2.00
	case FULL: // not implemented yet!
		xIndex, xAnswer = 0.00, 0.00
	}

	// create the index. The value is chosen randomly
	latexIndex := latexIndex{
		label: "label1",
		id:    "num1",
		pos:   position{x: xIndex, y: 0.75},
		value: 1 + rand.Intn(9),
	}

	// create the box to write the answer
	latexAnswer := latexSequenceAnswer{
		label:         "label2",
		pos:           position{x: xAnswer, y: 0.75},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the sequence problem
	latexSequenceProblem := sequenceProblem{
		index:   latexIndex,
		answer:  latexAnswer,
		seqtype: seqtype,
	}

	return latexSequenceProblem.Execute(seqtype)
}

// Get the LaTeX code for generating a problem on sequences of the
// given type: to either guess the previous (seqtype 1), subsequent
// (2), or both (3). Indices are strictly less than 100
func (masterFile MasterFile) GetSequence100(seqtype int) (latexCode string) {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// decide first the x-position of the index and the box for
	// drawing the answer
	var xIndex, xAnswer float64
	switch seqtype {
	case PREVIOUS:
		xIndex, xAnswer = 2.00, 0.50
	case SUBSEQUENT:
		xIndex, xAnswer = 0.50, 2.00
	case FULL: // not implemented yet!
		xIndex, xAnswer = 0.00, 0.00
	}

	// create the index. The value is chosen randomly
	latexIndex := latexIndex{
		label: "label1",
		id:    "num1",
		pos:   position{x: xIndex, y: 0.75},
		value: 1 + rand.Intn(99),
	}

	// create the box to write the answer
	latexAnswer := latexSequenceAnswer{
		label:         "label2",
		pos:           position{x: xAnswer, y: 0.75},
		minimumWidth:  1.5,
		minimumHeight: 1.25,
	}

	// and use all of these to create the sequence problem
	latexSequenceProblem := sequenceProblem{
		index:   latexIndex,
		answer:  latexAnswer,
		seqtype: seqtype,
	}

	return latexSequenceProblem.Execute(seqtype)
}

// templates
// ----------------------------------------------------------------------------

// Writes into the specified dst file the result of instantiating the
// given master file. For a full description, see the manual.
func (masterFile MasterFile) MasterToFileFromTemplate(dst string) {

	// verify that the given master file exists and is accessible
	masterisregular, _ := fstools.IsRegular(masterFile.Infile)
	if !masterisregular {
		log.Fatalf("the master file '%s' does not exist or is not accessible",
			masterFile.Infile)
	}

	// access a template and parse its contents
	master, err := template.ParseFiles(masterFile.Infile)
	if err != nil {
		log.Fatal(err)
	}

	// if the given filename already exists, then number it and so
	// on until the resulting filename does not exist. If
	// re-numbering is required, start with index 2
	index := 2
	current := dst
	for _, err := os.Stat(dst); err == nil; {
		log.Printf("The file '%v' already exists", dst)

		// renumber this filename
		dst = fstools.NumberFilename(current, index)

		// move forward to the next index and perform the test
		// again
		index += 1
		_, err = os.Stat(dst)
	}

	// now, open the file in read/write mode
	file, err := os.Create(dst)
	if err != nil {
		log.Fatalf("It was not possible to create the file '%v'", dst)
	}

	// make sure the file is closed before leaving
	defer file.Close()

	// and now execute the template
	err = master.Execute(file, masterFile)
	if err != nil {
		log.Fatal(err)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

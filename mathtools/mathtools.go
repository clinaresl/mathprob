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
	"errors"
	"fmt"
	"log" // logging services
	"math"
	"math/rand"
	"os" // access to file mgmt functions
	"path"
	"text/template"
	"time"

	// go facility for processing templates
	"github.com/clinaresl/mathprob/fstools"
	"github.com/clinaresl/mathprob/mathtools/components"
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
	Class   string
	Outfile string
}

// functions
// ----------------------------------------------------------------------------

// Create a new instance of a master file with the given name and clas
func NewMasterFile(filename, name, class string) MasterFile {

	return MasterFile{Infile: filename, Name: name, Class: class}
}

// helpers

// compute the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// compute the maximum of two floats
func max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

// return the number of digits of number n
func nbdigits(n int) int {
	return 1 + int(math.Ceil(math.Log10(float64(n))))
}

// return a random number with exactly n digits
func randN(n int) int {
	lower := int(math.Pow(float64(10), float64(n)-1))
	upper := int(math.Pow(float64(10), float64(n)))
	return lower + rand.Int()%(upper-lower)
}

// return true if and only if the given value has been found in the
// specified slice
func find(item string, container []string) bool {

	// for all items in the container
	for _, value := range container {

		// in case it has been found, then exit immediately
		if value == item {
			return true
		}
	}

	// if it has not been found after traversing the container,
	// then return false
	return false
}

// verify that the keys given in dict are correct for defining
// divisions. A dictionary is correct if and only if all the mandatory
// arguments have been given. If not, an error is raised and execution
// is aborted. Unnecessary keys are reported
func verifyDivisionDict(dict map[string]interface{}) {

	// the mandatory keys are given next
	mandatory := []string{"nbdvdigits", "nbdrdigits", "nbqdigits"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then
		// raise an error and exit
		if _, ok := dict[key]; !ok {
			log.Fatalf(" Fatal Error: Mandatory key '%v' for defining a division not found", key)
		}
	}

	// next, verify if there are some parameters
	for key := range dict {

		// if this key was not requested then report a message
		if !find(key, mandatory) {
			log.Printf(" Warning: The key '%v' is not necessary for creating a division and it will be ignored", key)
		}
	}
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

// Return the student's class of this master file
func (masterFile MasterFile) GetClass() string {
	return masterFile.Class
}

// Return the output tex filename that shall contain the exercises in tex
func (masterFile MasterFile) GetOutfile() string {
	return masterFile.Outfile
}

// the following function is provided just to allow the text/template to repeat
// the same statement an arbitrary number of times. It just returns a slice of
// MasterFiles of a given length. Each element can then be used to invoke the
// various services provided for text/templates
func (masterFile MasterFile) Slice(n int) []MasterFile {
	return make([]MasterFile, n)
}

// Components
// ----------------------------------------------------------------------------

func (masterFile MasterFile) GetCoordinate(dict map[string]interface{}) string {

	// first things first, verify that the given dictionary is correct
	if _, err := components.VerifyCoordinateDict(dict); err != nil {
		log.Fatal(err)
	}

	// now, get the positionable item, either a point or a formula
	var pos components.Position

	// if a formula was given, then set the position to a formula
	if value, ok := dict["formula"]; ok {
		svalue := value.(string)
		pos = components.Formula(svalue)
	} else {

		// otherwise, the dictionary contains a point which is accessible via
		// the keys "x" and "y"
		x := dict["x"].(float64)
		y := dict["y"].(float64)
		pos = components.Point{X: x, Y: y}
	}

	// so that at this point a valid coordinate can be returned
	label := dict["label"].(string)
	coord := components.NewCoordinate(pos, label)

	// and return the string that represents this coordinate
	return coord.String()
}

// Simple Operations
// ----------------------------------------------------------------------------

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

// Sequences
// ----------------------------------------------------------------------------

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

// divisions
// ----------------------------------------------------------------------------

func (masterFile MasterFile) GetDivision(dict map[string]interface{}) string {

	// seed the random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// Verify the given keys in the dictionary are correct. Note
	// that the types are not verified, only the presence of the
	// keys
	verifyDivisionDict(dict)

	// Now, build the components of the division according to the given parameters

	// --coordinates
	label1 := coordinateExplicit{
		x: 0.0,
		y: 1 + 2.0*dict["nbqdigits"].(float64) + 0.5,
	}
	label1.label = "label1"

	label2 := coordinateFormula{
		formula: fmt.Sprintf(`$(label1) + %v*(\zerowidth, 0.0)$`,
			2.0+dict["nbdvdigits"].(float64)),
	}
	label2.label = "label2"

	label3 := coordinateFormula{
		formula: fmt.Sprintf(`$(label2) + (%v*\zerowidth, -\zeroheight)$`,
			0.5*(2+max(dict["nbdrdigits"].(float64), dict["nbqdigits"].(float64)))),
	}
	label3.label = "label3"

	line1 := coordinateFormula{
		formula: fmt.Sprintf(`$(label2) + (%v\zerowidth, -2*\zeroheight-0.15 cm)$`,
			2.0+dict["nbdvdigits"].(float64)),
	}
	line1.label = "line1"

	// --bounding box
	bottom := coordinateFormula{
		formula: fmt.Sprintf(`$(line1) + %v*(0.0, -\zeroheight-\baselineskip-0.5/%v*\zeroheight)$`,
			2.0*dict["nbqdigits"].(float64)-1.0,
			2.0*dict["nbqdigits"].(float64)-1.0),
	}
	bottom.label = "bottom"
	right := coordinateFormula{
		formula: fmt.Sprintf(`$(label2) + (%v*\zerowidth, \zeroheight)$`,
			2.0+max(dict["nbdrdigits"].(float64), dict["nbqdigits"].(float64))),
	}
	right.label = "right"
	bBox := boundingBox{
		bottom: bottom,
		right:  right,
	}

	// --split box
	coord1 := coordinateFormula{
		formula: `$(label2) + (0.0, \zeroheight)$`,
	}
	coord2 := coordinateFormula{
		formula: `$(label2) + (0.0, -\zeroheight)$`,
	}
	coord3 := coordinateFormula{
		formula: fmt.Sprintf(`$(label2) + %v*(\zerowidth, -\zeroheight/%v)$`,
			2.0+max(dict["nbdrdigits"].(float64), dict["nbqdigits"].(float64)),
			2.0+max(dict["nbdrdigits"].(float64), dict["nbqdigits"].(float64))),
	}
	sBox := splitBox{
		coord1: coord1,
		coord2: coord2,
		coord3: coord3,
	}

	// --answer
	answer := latexAnswer{
		width: 2.0 + max(dict["nbdrdigits"].(float64), dict["nbqdigits"].(float64)),
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
	if dict["nbqdigits"].(float64) < dict["nbdvdigits"].(float64)-dict["nbdrdigits"].(float64) {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", dict["nbqdigits"], dict["nbdvdigits"], dict["nbdrdigits"], dict["nbdvdigits"].(float64)-dict["nbdrdigits"].(float64))
		dict["nbqdigits"] = dict["nbdvdigits"].(float64) - dict["nbdrdigits"].(float64)
	}

	if dict["nbqdigits"].(float64) > dict["nbdvdigits"].(float64)-dict["nbdrdigits"].(float64)+1 {
		log.Printf(" It is not possible to generate quotients with %v digits if the dividend has %v digits and the divisor has %v digits. Thus, %v digits in the quotient are generated instead", dict["nbqdigits"], dict["nbdvdigits"], dict["nbdrdigits"], dict["nbdvdigits"].(float64)-dict["nbdrdigits"].(float64)+1)
		dict["nbqdigits"] = dict["nbdvdigits"].(float64) - dict["nbdrdigits"].(float64) + 1
	}

	// now, generate numbers in their corresponding range
	log.Print(dict)
	var qvalue int
	nbdivdigits := int(dict["nbdivdigits"].(float64))
	nbdrdigits := int(dict["nbdrdigits"].(float64))
	for nbdigits(qvalue) < nbdivdigits || qvalue == 0 {
		dividend.value = randN(nbdivdigits)
		divisor.value = randN(nbdrdigits)
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
	return divProblem.Execute()

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

	// access a template and parse its contents. In addition it
	// registers a function "dict" which allows the user to
	// introduce in the text template any arguments
	master := template.Must(template.New(path.Base(masterFile.Infile)).Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {

			// if the number of items is not even (as many
			// pairs of the form "Key" "Value" should be
			// given) then an error is raised
			if len(values)%2 != 0 {
				return nil, errors.New("Invalid dict call. There should be an even number of arguments of the form 'Key' 'Value'")
			}

			// Create a map with as many elements as keys
			// have been specified
			dict := make(map[string]interface{}, len(values)/2)

			// and process them
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("Dict keys must be strings")
				}
				dict[key] = values[i+1]
			}

			// at this point no error has been reported, move therefore back
			return dict, nil
		}}).ParseFiles(masterFile.Infile))

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

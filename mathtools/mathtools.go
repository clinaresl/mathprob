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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log" // logging services
	"math"
	"math/rand"
	"os" // access to file mgmt functions
	"text/template"
	"time"

	// go facility for processing templates
	"github.com/clinaresl/mathprob/fstools"
	"github.com/clinaresl/mathprob/helpers"
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

// 	// if the type was not recognized, then return an error
// 	return 0, fmt.Errorf("It was not possible to cast '%v' into an integer")
// }

// return a valid specification of a sequence with no error if all the keys
// given in dict are correct for defining a sequence. If not, an error is
// returned. If an error is returned, the contents of the sequence are
// undetermined
//
// A dictionary is correct if and only if it correctly provides a type of
// sequence with the keyword "type", a number of items with the keyword
// "nbitems", and a lower and upper bound with "geq" and "leq"
func verifySequenceDict(dict map[string]interface{}) (sequence, error) {

	// the mandatory keys are given next
	mandatory := []string{"type", "nbitems", "geq", "leq"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then
		// raise an error and exit
		if _, ok := dict[key]; !ok {
			return sequence{}, fmt.Errorf("Mandatory key '%v' for defining a sequence not found", key)
		}
	}

	// make also sure that parameters are given with the right type
	var err error
	var seqtype, nbitems, geq, leq int
	if seqtype, err = helpers.Atoi(dict["type"]); err != nil {
		return sequence{}, errors.New("the type of a sequence should be given as a integer")
	}
	if nbitems, err = helpers.Atoi(dict["nbitems"]); err != nil {
		return sequence{}, errors.New("the number of items in a sequence should be given as an integer")
	}
	if geq, err = helpers.Atoi(dict["geq"]); err != nil {
		return sequence{}, errors.New("the lower bound of a sequence should be given as an integer")
	}
	if leq, err = helpers.Atoi(dict["leq"]); err != nil {
		return sequence{}, errors.New("the upper bound of a sequence should be given as a string")
	}

	// finally, ensure the type is correct
	if seqtype < SEQNONE || seqtype > SEQBOTH {
		return sequence{}, fmt.Errorf("the type of a sequence given '%v' is incorrect", seqtype)
	}

	// next, verify if there are some unnecessary parameters
	for key := range dict {

		// if this key was not requested then report a message
		if !helpers.Find(key, mandatory) {
			log.Printf(" Warning: The key '%v' is not necessary for creating a sequence and it will be ignored", key)
		}
	}

	// otherwise, the dictionary is correct
	return sequence{seqtype: seqtype,
		nbitems: nbitems,
		geq:     geq,
		leq:     leq}, nil
}

// verify that the keys given in dict are correct for defining
// divisions. A dictionary is correct if and only if all the mandatory
// arguments have been given. If not, an error is raised and execution
// is aborted. Unnecessary keys are reported
func verifyDivisionDict(dict map[string]interface{}) (division, error) {

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

	// make also sure that parameters are given with the right type
	var err error
	var nbdvdigits, nbdrdigits, nbqdigits int
	if nbdvdigits, err = helpers.Atoi(dict["nbdvdigits"]); err != nil {
		return division{}, errors.New("the number of digits of the dividend should be given as a integer")
	}
	if nbdrdigits, err = helpers.Atoi(dict["nbdrdigits"]); err != nil {
		return division{}, errors.New("the number of digits of the divisor should be given as an integer")
	}
	if nbqdigits, err = helpers.Atoi(dict["nbqdigits"]); err != nil {
		return division{}, errors.New("the number of digits of the quotient should be given as an integer")
	}

	// next, verify if there are some unnecessary parameters
	for key := range dict {

		// if this key was not requested then report a message
		if !helpers.Find(key, mandatory) {
			log.Printf(" Warning: The key '%v' is not necessary for creating a division and it will be ignored", key)
		}
	}

	// now, return the proper definition of a division problem
	return division{
		nbdvdigits: nbdvdigits,
		nbdrdigits: nbdrdigits,
		nbqdigits:  nbqdigits}, nil
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

// This method is intended to be used in master files. It is substituted by TikZ
// contents that create a coordinate with a label (identified with the key
// "label") and located at a given position which can be identified either with
// a position (using both keys "x" and "y") or a formula, with the key
// "formula". The coordinates x and y must be given as floating-point numbers
// whereas formulas should be given as strings.
func (masterFile MasterFile) Coordinate(dict map[string]interface{}) string {

	// first things first, verify that the given dictionary is correct
	var err error
	var coord components.Coordinate
	if coord, err = components.VerifyCoordinateDict(dict); err != nil {
		log.Fatal(err)
	}

	// otherwise return the string that represents this coordinate
	return coord.String()
}

// This method is intended to be used in master files. It is substituted by TikZ
// contents that create a text box located at a coordinate (either by providing
// the coordinates of a Point or giving a Formula) with the contents
// specified in the key "text"
func (masterFile MasterFile) Text(dict map[string]interface{}) string {

	// first things first, verify that the given dictionary is correct
	var err error
	var text components.Text
	if text, err = components.VerifyTextDict(dict); err != nil {
		log.Fatal(err)
	}

	// and return the string that shows up the contents of this text box
	return text.String()
}

// This method is intended to be used in master files. It is substituted by TikZ
// contents that create a box located at a coordinate (either by providing the
// coordinates of a Point or giving a Formula) and with the contents specified
// in the key "text" which has the minimum width and height given in "minwidth"
// and "minheight"
func (masterFile MasterFile) Box(dict map[string]interface{}) string {

	// first things first, verify that the given dictionary is correct
	var err error
	var box components.Box
	if box, err = components.VerifyBoxDict(dict); err != nil {
		log.Fatal(err)
	}

	// and return the string that shows up the contents of this box
	return box.String()
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
			latexOperandB.value = helpers.Min(latexOperandA.value, latexOperandB.value)
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
			latexOperandB.value = helpers.Min(latexOperandA.value, latexOperandB.value)
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
			latexOperandB.value = helpers.Min(latexOperandA.value, latexOperandB.value)
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
			latexOperandB.value = helpers.Min(latexOperandA.value, latexOperandB.value)
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

// Return the LaTeX code in TikZ format that generates a sequence with the
// keywords given in the dictionary. The type of a sequence is: "first", if the
// first number has to be given; "last", if the last number has to be given;
// "none" or "both" if either none of them or both have to displayed. In
// addition, a sequence is made up of a number of items, each one greater or
// equal than a given threshold and lower or equal than another bound using the
// keywords "geq" and "leq" respectively
func (masterFile MasterFile) Sequence(dict map[string]interface{}) string {

	// verify the given dictionary is correct and get an instance of a valid
	// sequence
	sequence, err := verifySequenceDict(dict)
	if err != nil {
		log.Fatalf("The dictionary given for creating a sequence is incorrect: %v", err)
	}

	// and return the LaTeX/TikZ code for representing this sequence
	return sequence.execute()
}

// divisions
// ----------------------------------------------------------------------------

// Return the LaTeX code in TikZ format that generates a division with the
// keywords given in the dictionary:
//
// nbdvdigits: number of digits of the dividend
// nbdrdigits: number of digits of the divisor
// nbqdigits: number of digits of the quotient
func (masterFile MasterFile) Division(dict map[string]interface{}) string {

	// Verify the given keys in the dictionary are correct. Note
	// that the types are not verified, only the presence of the
	// keys. In case of an error, just generate a fatal error
	div, err := verifyDivisionDict(dict)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return div.execute()
}

// templates
// ----------------------------------------------------------------------------

// Parse the template given in contents to a masterfile and returns the result
// in a buffer, and nil if no error was found
func (masterFile MasterFile) masterToBufferFromTemplate(contents string) (bytes.Buffer, error) {

	// create the buffer to return the result of the execution
	var result bytes.Buffer

	// access a template and parse its contents. In addition it registers a
	// function "dict" which allows the user to introduce in the text template
	// any arguments
	t := template.Must(template.New(contents).Funcs(template.FuncMap{
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
		}}).Parse(contents))

	// execute the template with the information in this instance
	err := t.Execute(&result, masterFile)
	if err != nil {

		// note that the result might contain some partial results
		return result, err
	}

	// at this point return the result with no error
	return result, nil
}

// Writes into the specified dst file the result of instantiating the
// given master file
func (masterFile MasterFile) MasterToFileFromTemplate(dst string) {

	// verify that the given master file exists and is accessible
	masterisregular, _ := fstools.IsRegular(masterFile.Infile)
	if !masterisregular {
		log.Fatalf("the master file '%s' does not exist or is not accessible",
			masterFile.Infile)
	}

	// these files are expected to be not too long, actually, so read the entire
	// contents of the file into main memory
	contents, err := ioutil.ReadFile(masterFile.Infile)
	if err != nil {
		log.Fatalf("It was not possible to read the input file '%v'", masterFile.Infile)
	}

	// if the given filename already exists, then number it and so on until the
	// resulting filename does not exist. If re-numbering is required, start
	// with index 2
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

	// execute the template
	result, err := masterFile.masterToBufferFromTemplate(string(contents))
	if err != nil {
		log.Fatalf("Error when executing the template over the master file", result)
	}

	// and write the result in the output file
	if _, err := file.WriteString(result.String()); err != nil {
		log.Fatalf("Error while writing the result of a template in '%v'", dst)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

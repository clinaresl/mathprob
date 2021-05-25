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
	"os"  // access to file mgmt functions
	"text/template"

	// go facility for processing templates
	"github.com/clinaresl/mathprob/fstools"
	"github.com/clinaresl/mathprob/helpers"
	"github.com/clinaresl/mathprob/mathtools/components"
)

// types
// ----------------------------------------------------------------------------

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

// return a valid specification of a basic operation with no error ir all the
// keys given in dict are correct for defining a basic sequence. If not, an
// error is returned. If an error is returned, the contents of the basic
// operation are undefined
//
// A dictionary is correct if and only if it correctly provides a type of basic
// operation with the keyword "type", a number of digits of the operands, and
// the result, and the number of operands to show.
func verifyBasicOperationDict(dict map[string]interface{}) (basicOperation, error) {

	// the mandatory keys are given next
	mandatory := []string{"type", "operator", "nboperands", "nbdigitsop", "nbdigitsrslt"}

	// now, verify that all mandatory parameters are present in the dict
	for _, key := range mandatory {

		// if a mandatory parameter has not been given, then
		// raise an error and exit
		if _, ok := dict[key]; !ok {
			return basicOperation{}, fmt.Errorf("Mandatory key '%v' for defining a basic operation not found", key)
		}
	}

	// make also sure that parameters are given with the right type
	var ok bool
	var err error
	var operator string
	var botype, nboperands, nbdigitsop, nbdigitsrslt int
	if operator, ok = dict["operator"].(string); !ok {
		return basicOperation{}, errors.New("The operator of a basic operation should be given as a stirng")
	} else {
		operators := []string{"+", "-", "*", "/"}
		if !helpers.Find(operator, operators) {
			return basicOperation{}, errors.New("The operator of a basic operation has to be one and only one among the following: '+', '-', '*' or '/'")
		}
	}
	if botype, err = helpers.Atoi(dict["type"]); err != nil {
		return basicOperation{}, errors.New("the type of a basic operation should be given as an integer")
	}
	if nboperands, err = helpers.Atoi(dict["nboperands"]); err != nil {
		return basicOperation{}, errors.New("the number of operands in a basic operation should be given as an integer")
	}
	if nbdigitsop, err = helpers.Atoi(dict["nbdigitsop"]); err != nil {
		return basicOperation{}, errors.New("the number of digits of all operands should be given as an integer")
	}
	if nbdigitsrslt, err = helpers.Atoi(dict["nbdigitsrslt"]); err != nil {
		return basicOperation{}, errors.New("the number of digits of the result of a basic operation should be given as a string")
	}

	// finally, ensure the type is correct
	if botype < BORESULT || botype > BOOPERAND {
		return basicOperation{}, fmt.Errorf("the type of a basic operation given '%v' is incorrect", botype)
	}

	// next, verify if there are some unnecessary parameters
	for key := range dict {

		// if this key was not requested then report a message
		if !helpers.Find(key, mandatory) {
			log.Printf(" Warning: The key '%v' is not necessary for creating a basic operation and it will be ignored", key)
		}
	}

	// otherwise, the dictionary is correct
	return basicOperation{
		botype:       botype,
		operator:     operator,
		nboperands:   nboperands,
		nbdigitsop:   nbdigitsop,
		nbdigitsrslt: nbdigitsrslt,
	}, nil
}

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
		return sequence{}, errors.New("the type of a sequence should be given as an integer")
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
	return sequence{
		seqtype: seqtype,
		nbitems: nbitems,
		geq:     geq,
		leq:     leq,
	}, nil
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
		nbqdigits:  nbqdigits,
	}, nil
}

// methods
// ----------------------------------------------------------------------------

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

// TikZ reusable components
//
// The following meethods provide direct access to the TikZ reusable components
// to be used in a master file directly
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

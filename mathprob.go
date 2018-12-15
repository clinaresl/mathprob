/*
  mathprob.go
  Description: Automatic generator of Math Problems
  -----------------------------------------------------------------------------

  Started on  <Mon Jul 10 09:31:10 2017 >
  Last update <>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by
  Login   <carlos.linares@uc3m.es>
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/mathprob/fstools"
	"bitbucket.org/mathprob/mathtools"
)

// imports
// ----------------------------------------------------------------------------

// global variables
// ----------------------------------------------------------------------------
var VERSION string = "0.1.0" // current version
var EXIT_SUCCESS int = 0     // exit with success

// Options
var masterFilename string // master file
var texFilename string    // output tex filename
var jsonFilename string   // JSON filename with info of all records to process
var studentName string    // student's name
var helpMaster bool       // is help on master files requested?
var helpJSON bool         // is help about JSON files requested?
var verbose bool          // has verbose output been requested?
var version bool          // has version info been requested?

// functions
// ----------------------------------------------------------------------------

// initializes the command-line parser
func init() {

	// Flag to store the master file to process
	flag.StringVar(&masterFilename, "infile", "", "master file to use for generating the sheets of exercises. If a JSON file is given also, this parameter is automatically discarded. Use '-help-master' to obtain additional information")
	flag.StringVar(&texFilename, "outfile", "", "output filename with the TeX code of the exercises generated from the template file. If not given, then the student's name provided with -student-name is used instead. If none is provided, then 'main.tex' is used by default. In case the resulting TeX file already exists, then it is re-numbered to avoid overwritting existing contents")
	flag.StringVar(&jsonFilename, "json-file", "", "file with information of all records to process in JSON format. If a JSON file is given, the input file given with -infile is automatically discarded. It is not allowed to provide more than 1024 records in the JSON file. Use 'help-json' to obtain additional information")
	flag.StringVar(&studentName, "name", "", "Student's name")

	flag.BoolVar(&helpMaster, "help-master", false, "provides information about the format and usage of master files")
	flag.BoolVar(&helpMaster, "help-json", false, "provides information about the JSON format used to specify multiple records")

	// other optional parameters are verbose and version
	flag.BoolVar(&verbose, "verbose", false, "provides verbose output")
	flag.BoolVar(&version, "version", false, "shows version info and exists")
}

// shows version info and exists with the specified signal
func showVersion(signal int) {

	fmt.Printf("\n %v", os.Args[0])
	fmt.Printf("\n Version: %v\n\n", VERSION)
	os.Exit(signal)
}

// shows informmation on master files
func showHelpMaster(signal int) {

	fmt.Println(` 
 TBW
`)
	os.Exit(signal)
}

// shows informmation on master files
func showHelpJSON(signal int) {

	fmt.Println(` 
 TBW
`)
	os.Exit(signal)
}

// parse the flags and verifies that proper values were given. If not, a fatal
// error is raised
func verify() {

	// if version information was requested show it now and exit
	if version {
		showVersion(EXIT_SUCCESS)
	}

	// in case further assistance on a particular subject is requested, then
	// show it here and exit
	if helpMaster {
		showHelpMaster(EXIT_SUCCESS)
	}
	if helpJSON {
		showHelpJSON(EXIT_SUCCESS)
	}

	// verify that a master file has been given
	if masterFilename == "" && jsonFilename == "" {
		log.Fatalf("Use either -master-file or -json-file to provide a master file. See -help for more details")
	}

	// if a student's name has not been provided, issue a warning
	// as it might be used in the master file
	if studentName == "" && jsonFilename == "" {
		log.Println("No student's name has been provided!")
	}
}

// the following function applies the following rules to derive the TeX filename:
//
//    1. If -tex-file has been used, then return immediately the user's choice
//    2. Otherwise, if a student's name was given, then use it instead
//    3. If none has been provided, then use 'main.tex' by default
func getTexName() string {

	//    1. If -tex-file has been used, then return immediately the user's choice
	if texFilename != "" {
		return fstools.AddSuffix(texFilename, ".tex")
	} else {
		//    2. Otherwise, if a student's name was given,
		//    then use it instead
		if studentName != "" {
			return fstools.AddSuffix(studentName, ".tex")
		} else {
			//    3. If none has been provided, then use
			//    'main.tex' by default
			return "main.tex"
		}
	}

}

// Main body
func main() {

	// first, parse the flags
	flag.Parse()

	// verify the values parsed
	verify()

	// in case a JSON file was provided
	if jsonFilename != "" {

		// Unmarshal the JSON file to get all records to process
		jsonData, _ := ioutil.ReadFile(jsonFilename)
		var records []mathtools.MasterFile
		records = make([]mathtools.MasterFile, 5)
		_ = json.Unmarshal([]byte(jsonData), &records)

		fmt.Println()
		for _, field := range records {

			// show info
			fmt.Println(" * Processing ...")
			fmt.Printf("\t Master file    : %s\n", field.GetInfile())
			fmt.Printf("\t Student's name : %v\n", field.GetName())
			fmt.Printf("\t TeX file       : %v\n\n", field.GetOutfile())

			// process this specific record
			masterFile := mathtools.NewMasterFile(field.GetInfile(), field.GetName())
			masterFile.MasterToFileFromTemplate(fstools.AddSuffix(field.GetOutfile(),
				".tex"))
		}
	} else {

		// Otherwise, use the parameters given by the user to
		// generate a unique TeX file

		// get the tex filename and show it on the standard output
		texFilename = getTexName()
		log.Printf("TeX filename: %s\n", texFilename)

		// now, instantiate the master file with the data generated
		masterFile := mathtools.NewMasterFile(masterFilename, studentName)
		masterFile.MasterToFileFromTemplate(texFilename)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

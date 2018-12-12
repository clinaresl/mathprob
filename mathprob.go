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
	"flag"
	"fmt"
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
var helpMaster bool       // is help on master files requested?
var verbose bool          // has verbose output been requested?
var version bool          // has version info been requested?

// functions
// ----------------------------------------------------------------------------

// initializes the command-line parser
func init() {

	// Flag to store the master file to process
	flag.StringVar(&masterFilename, "master", "", "master file to use for generating the sheets of exercises")
	flag.BoolVar(&helpMaster, "help-master", false, "provides information about the format and usage of master files")

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

// parse the flags and verifies that proper values were given. If not, a fatal
// error is raised
func verify() {

	// first, parse the flags ---in case help was given, it is automatically
	// handled by the flags package
	flag.Parse()

	// if version information was requested show it now and exit
	if version {
		showVersion(EXIT_SUCCESS)
	}

	// in case further assistance on a particular subject is requested, then
	// show it here and exit
	if helpMaster {
		showHelpMaster(EXIT_SUCCESS)
	}

	// verify that a master file has been given
	if masterFilename == "" {
		log.Fatalf("Use --master to provide a master file. See --help for more details")
	}

	// verify that the given master file exists and is accessible
	masterisregular, _ := fstools.IsRegular(masterFilename)
	if !masterisregular {
		log.Fatalf("the master file '%s' does not exist or is not accessible",
			masterFilename)
	}
}

// Main body
func main() {

	// verify the values parsed
	verify()

	// now, instantiate the master file with the data generated
	masterFile := mathtools.NewMasterFile(masterFilename)
	masterFile.MasterToFileFromTemplate("main.tex")
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

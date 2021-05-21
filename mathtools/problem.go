// -*- coding: utf-8 -*-
// problem.go
// -----------------------------------------------------------------------------
//
// Started on <lun 17-05-2021 22:03:12.129169010 (1621281792)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

//
// Description
//
package mathtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// types
// ----------------------------------------------------------------------------

// A master file consists of an input filename that stores the
// tempalte to fill in to generate the final sheet of exercises, and
// an output tex filename. It also comes with other fields that can be
// used for customizing the resulting file such as the student's name

// A master problem consists of a number of arbitrary arguments of any type
// indexed by a string, a specific type and a number of problems to generate
type MasterProblem struct {
	probtype string
	args     map[string]interface{}
	nbprobs  int
}

// functions
// ----------------------------------------------------------------------------

// verify the type of problem requested is acknowledged by mathprob. In case it
// is acknowledged, an error is returned; otherwise nil
func AckProblemType(probtype string) error {

	// if this problem is not among those supporting the automated generation of
	// json problems then immediately return an error
	if strings.ToUpper(probtype) != "SEQUENCE" {
		return fmt.Errorf("Unknown problem type '%v'", probtype)
	}

	// otherwise acknwoledge by saying no error
	return nil
}

// return an array of instances of MasterProblem from the contents of a json
// file. In case it is not possible to unmarshall the contents of the json file,
// then an error is returned and the contents of the slice are undefined
func Unmarshall(data []byte) (output []MasterProblem, err error) {

	// first things first, decode the data in the JSON file, which is expected
	// to be a slice of entries, each specifying a different problem type
	var jsondata []map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(&jsondata); err != nil {
		return output, errors.New("Error while decoding JSON data to generate instances of master problems")
	}

	// ensure that data has been properly processed
	for _, entry := range jsondata {

		var ok bool

		// Right now just unmarshall this specific entry to an instance of
		// MasterProblem by decoding each entry separately

		// First, process the type. Verify first that it exists as a key in the
		// JSON input file; if so, extract its value
		var probtype string
		if _, ok = entry["type"]; !ok {
			return output, errors.New("The type of problem to generate has not been found!")
		} else {
			if probtype, ok = entry["type"].(string); !ok {
				return output, errors.New("The problem type could not be casted into a string")
			}
		}
		// Second, process the number of problems to generate. Verify first that
		// it exists as a key in the JSON input file; if so, extract its value
		var nbprobs int
		if _, ok = entry["nbprobs"]; !ok {
			return output, errors.New("The number of problems has not been found!")
		} else {
			if _, ok = entry["nbprobs"].(float64); !ok {
				return output, errors.New("The number of problems to generate could not be casted into an integer")
			} else {
				nbprobs = int(entry["nbprobs"].(float64))
			}
		}

		// Thirdly, process the arguments. Verify first that they exist as a map
		// of keys to values of any type; if so, extract them
		var args map[string]interface{} = make(map[string]interface{})
		if _, ok = entry["args"]; !ok {
			return output, errors.New("The arguments for generating a type of problem has not been found!")
		} else {
			if _, ok = entry["args"].(map[string]interface{}); !ok {
				return output, errors.New("The arguments are not given as a dict of strings to values of any type")
			} else {

				// then process all entries in this dictionary, looking for
				// specific keys
				jsonargs := entry["args"].(map[string]interface{})
				for _, ikey := range []string{"type", "nbitems", "geq", "leq"} {

					// if this key has not been given, then immediately issue an
					// error
					if _, ok = jsonargs[ikey]; !ok {
						return output, fmt.Errorf("The arg '%v' has not been found!", ikey)
					} else {
						if ivalue, ok := jsonargs[ikey].(interface{}); !ok {
							return output, fmt.Errorf("It was not possible to decode the value of the arg '%v'", ikey)
						} else {

							// store this key
							args[ikey] = ivalue
						}
					}
				}
			}
		}

		// and generate a master problem
		masterProblem := MasterProblem{
			probtype: probtype,
			args:     args,
			nbprobs:  nbprobs,
		}
		output = append(output, masterProblem)
	}

	// and finally return with the data computed so far
	return
}

// given an array of master problems (of any type) return a slice of bytes in
// JSON format with the requested problems. If a problem could not be generated,
// the contents of the returned data are undefined and an error is raised
func GenerateJSON(problems []MasterProblem) (data []byte, err error) {

	// for all problems
	for _, problem := range problems {

		// each master problem requests a specific number of instances to
		// generate
		for i := 0; i < problem.nbprobs; i++ {

			// get a specific representation of this problem as a JSON stream
			if jsondata, err := problem.getJSONProblem(); err != nil {
				return data, err
			} else {

				// and if everything went fine then add it to the overall JSON
				// data stream to return
				data = append(data, jsondata...)
			}
		}
	}

	// Before leaving, enclose all problems in a json slice
	data = append([]byte("[\n"), append(data, []byte("\n]")...)...)

	return
}

// Methods
// ----------------------------------------------------------------------------
// Return a JSON string with instances of all problems specified in the master
// problem. Note an error is returned if anything goes wrong
func (m MasterProblem) getJSONProblem() (data []byte, err error) {

	// first things first, verify that this problem has a correct type
	if err := AckProblemType(m.probtype); err != nil {
		return data, err
	}

	// given that this problem is correct, verify its arguments
	switch strings.ToUpper(m.probtype) {

	case "SEQUENCE":

		// First, verify that all items in the dictionary of args are correct
		if instance, err := verifySequenceDict(m.args); err != nil {
			return data, err
		} else {

			// if so, generate a JSON stream with the representation of this
			// specific problem
			if data, err := instance.GenerateJSONInstance(); err != nil {
				return data, err
			} else {

				// if everything went on correctly, then return the JSON data
				// stream
				return data, nil
			}
		}

	default:
		return data, fmt.Errorf("Unsupported generation of JSON problems for problem type '%v'", m.probtype)
	}

	// Data should have been returned before, this is just to avoid compiler errors
	return
}

// Local Variables:
// mode:go
// fill-column:80
// End:

// -*- coding: utf-8 -*-
// problem.go
// -----------------------------------------------------------------------------
//
// Started on <lun 17-05-2021 22:03:12.129169010 (1621281792)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package mathtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// This file contains general functions for handling requests to automatically
// generate problems of different types in JSON format. It is intended to serve
// as a bridge to connect mathprob to a front-end. It works from a list of
// requests given also in JSON format

// types
// ----------------------------------------------------------------------------

// A master problem consists of a number of arbitrary arguments of any type
// indexed by a string, a specific type and a number of problems to generate
type MasterProblem struct {
	probtype string
	args     map[string]interface{}
	nbprobs  int
}

// A problem in JSON format consists mainly of two fields: the arguments of the
// problem and its solution. Those records in the arguments of the problem that
// have to be filled in by the student are marked with a question mark "?". In
// addition, different problems might have different types and thus, a probtype
// field is given also
type problemJSON struct {
	Probtype string   `json:"type"`
	Id       int      `json:"id"`
	Args     []string `json:"args"`
	Solution []string `json:"solution"`
}

// functions
// ----------------------------------------------------------------------------

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

		// Thirdly, ensure the arguments are given and extract'em
		var args map[string]interface{} = make(map[string]interface{})
		if _, ok = entry["args"]; !ok {
			return output, errors.New("The arguments for generating a type of problem has not been found!")
		} else {
			if args, ok = entry["args"].(map[string]interface{}); !ok {
				return output, errors.New("The arguments are not given as a dict of strings to values of any type")
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

	// -- initialization: create a slice of JSON problems where each request is
	//                    filled in. These is the slice to marshal
	var jsonprobs []problemJSON

	// for all problems
	for _, problem := range problems {

		// each master problem requests a specific number of instances to
		// generate
		for i := 0; i < problem.nbprobs; i++ {

			// depending upon the type of problem to generate
			switch strings.ToUpper(problem.probtype) {

			case "BASICOPERATION":

				// First, verify that all items in the dictionary of args are correct
				if instance, err := verifyBasicOperationDict(problem.args); err != nil {
					return data, err
				} else {

					// if so, generate a JSON stream with the representation of this
					// specific problem
					if iprob, err := instance.generateJSONProblem(); err != nil {
						return data, err
					} else {

						// if everything went on correctly, then correctly
						// number this problem and add this problem to the slice
						// of problems to marshal
						iprob.Id = i
						jsonprobs = append(jsonprobs, iprob)
					}
				}

			case "DIVISION":

				// First, verify that all items in the dictionary of args are correct
				if instance, err := verifyDivisionDict(problem.args); err != nil {
					return data, err
				} else {

					// if so, generate a JSON stream with the representation of this
					// specific problem
					if iprob, err := instance.generateJSONProblem(); err != nil {
						return data, err
					} else {

						// if everything went on correctly, then correctly
						// number this problem and add this problem to the slice
						// of problems to marshal
						iprob.Id = i
						jsonprobs = append(jsonprobs, iprob)
					}
				}

			case "MULTIPLICATIONTABLE":

				// First, verify that all items in the dictionary of args are correct
				if instance, err := verifyMultiplicationTableDict(problem.args); err != nil {
					return data, err
				} else {

					// if so, generate a JSON stream with the representation of this
					// specific problem
					if iprob, err := instance.generateJSONProblem(); err != nil {
						return data, err
					} else {

						// if everything went on correctly, then correctly
						// number this problem and add this problem to the slice
						// of problems to marshal
						iprob.Id = i
						jsonprobs = append(jsonprobs, iprob)
					}
				}

			case "SEQUENCE":

				// First, verify that all items in the dictionary of args are correct
				if instance, err := verifySequenceDict(problem.args); err != nil {
					return data, err
				} else {

					// if so, generate a JSON stream with the representation of this
					// specific problem
					if iprob, err := instance.generateJSONProblem(); err != nil {
						return data, err
					} else {

						// if everything went on correctly, then correctly
						// number this problem and add this problem to the slice
						// of problems to marshal
						iprob.Id = i
						jsonprobs = append(jsonprobs, iprob)
					}
				}

			default:
				return data, fmt.Errorf("Unsupported generation of JSON problems for problem type '%v'", problem.probtype)
			}
		}
	}

	// Now, marshal data and return the json bytes stream. Note that this
	// function returns straight away the same error returned by the Marshal
	// function
	data, err = json.MarshalIndent(jsonprobs, "", "\t")
	return data, err
}

// Local Variables:
// mode:go
// fill-column:80
// End:

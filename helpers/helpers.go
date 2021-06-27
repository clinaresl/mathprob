// -*- coding: utf-8 -*-
// helpers.go
// -----------------------------------------------------------------------------
//
// Started on <dom 23-05-2021 19:43:54.087493637 (1621791834)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Definition of general purpose functions used in the implementation of the
// math tools
package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// transform the input into a bool by making sure the input is either an int or
// a string. In case an integer is given, 0 is false and any other value is 1;
// if a string is given, "" and "false" (with any mixture of upper/lower case
// letter) is false and any other string is true. In case it is not possible,
// the value returned is undefined and an error is signaled
func Atob(n interface{}) (bool, error) {

	switch value := n.(type) {
	case int:
		return value != 0, nil
	case string:
		return value != "" && strings.ToLower(value) != "false", nil
	}

	// if the type has not been recognized, then return an error
	return false, fmt.Errorf("It was not possible to cast '%v' into a bool")
}

// transform the input into an integer by making sure that the input is either
// an int, a float or a string. In case it is not possible, the value returned
// is undefined and an error is signaled
func Atoi(n interface{}) (int, error) {

	switch value := n.(type) {
	case int:
		return value, nil
	case float32:
		return int(value), nil
	case float64:
		return int(value), nil
	case string:
		if result, err := strconv.Atoi(value); err != nil {
			return 0, err
		} else {
			return result, nil
		}
	}

	// if the type was not recognized, then return an error
	return 0, fmt.Errorf("It was not possible to cast '%v' into an integer")
}

// return true if and only if the given value has been found in the
// specified slice
func Find(item string, container []string) bool {

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

// compute the maximum of two floats
func Max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

// compute the minimum of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// return the number of digits of number n. In case the number is negative, then
// 1 is added to display the unary -
func NbDigits(n int) int {

	// because we use log10 to compute the number of digits of any number, we
	// have to consider separately the case of 0
	if n == 0 {
		return 1
	} else if n < 0 {

		// also, if a number is negative, we should use its magnitude and add 1
		// accounting for the unary -
		return 2 + int(math.Log10(math.Abs(float64(n))))
	}

	// if the number is strictly positive, then
	return 1 + int(math.Log10(float64(n)))
}

// return a random number with exactly n digits
func RandN(n int) int {
	lower := int(math.Pow(float64(10), float64(n)-1))
	upper := int(math.Pow(float64(10), float64(n)))
	return lower + rand.Int()%(upper-lower)
}

// In case any of the arguments given in args does not appear in the specified
// dictionary then return an error explicitly mentioning the missing key.
// Otherwise, return no error
func VerifyArgs(dict map[string]interface{}, args []string) error {

	// for all entries in the slice
	for _, key := range args {

		// if this key is not found, then immediately return false
		if _, ok := dict[key]; !ok {
			return fmt.Errorf("Missing key '%v'", key)
		}
	}

	// at this point everything went fine!
	return nil
}

// Return whether any of the keys in the specified dict does not appear in the
// slice of args. If so, it explicitly mentions the missing key; otherwise,
// return false and the value of the second return value is undefined
func VerifyKeys(dict map[string]interface{}, args []string) (bool, string) {

	// for all keys in the dictionary
	for key := range dict {

		// if this key does not exist in the slice of arguments
		if !Find(key, args) {

			// then return false and the offending key
			return false, key
		}
	}

	// at this point return true with any string
	return true, ""
}

// Local Variables:
// mode:go
// fill-column:80
// End:

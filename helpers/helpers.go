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
)

// compute the minimum of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// compute the maximum of two floats
func Max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
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

// transform the input into an integer by making sure that the input is either
// an int, a float or a string. In case it is not possible, the value return is
// undefined and an error is signaled
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

// Local Variables:
// mode:go
// fill-column:80
// End:

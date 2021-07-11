// -*- coding: utf-8 -*-
// mystery_operation.go
// -----------------------------------------------------------------------------
//
// Started on <dom 11-07-2021 03:35:49.505779062 (1625967349)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

//
// Description
//
package mathtools

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/clinaresl/mathprob/helpers"
)

// constants
// ----------------------------------------------------------------------------

// types
// ----------------------------------------------------------------------------

// A mystery operation consists of an incomplete operation of any type (+, -, /,
// *) where some digits have been masked in the first operand, second, or the
// result or any combination of these.
type mysteryOperation struct {

	// number of digits of the first and second operand
	nbdigits1, nbdigits2 int

	// number of masked digits of the first and second operand
	nbmasked1, nbmasked2 int

	// number of digits and number of masked digits in the answer
	nbdigitsanswer, nbmaskedanswer int

	// operator
	operator string
}

// methods
// ----------------------------------------------------------------------------

// -- mysteryOperation

// return the instance of a specific mystery operation that can be marshalled in
// JSON format. The receiver is assumed to have been fully verified so that it
// should be consistent
//
// The result is given as an array of numbers:
//    1. The first string is the operator
//    2. The 2nd-3th strings are the number of digits of the first and second
//    operand
//    3. The 4th string is the number of digits of the answer
//    4. Next, all digits of both operands and the digits of the answer are
//    given consecutively. If one item has to be guessed it is masked with a
//    question mark "?"
func (mo mysteryOperation) generateJSONProblem() (problemJSON, error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// create a slice with all digits to choose from
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	// first, verify that the values given to the arguments of this mystery
	// operation make sense

	// first of all, ensure there are no more masked digits in each item than
	// digits in it
	if mo.nbmasked1 > mo.nbdigits1 {
		return problemJSON{}, fmt.Errorf("There are more masked digits (%v) in the first operand than digits in it (%v)",
			mo.nbmasked1, mo.nbdigits1)
	}
	if mo.nbmasked2 > mo.nbdigits2 {
		return problemJSON{}, fmt.Errorf("There are more masked digits (%v) in the second operand than digits in it (%v)",
			mo.nbmasked2, mo.nbdigits2)
	}
	if mo.nbmaskedanswer > mo.nbdigitsanswer {
		return problemJSON{}, fmt.Errorf("There are more masked digits (%v) in the answer than digits in it (%v)",
			mo.nbmaskedanswer, mo.nbdigitsanswer)
	}

	// also, make sure that the number of digits makes sense
	switch mo.operator {

	case "+":
		if mo.nbdigitsanswer < int(helpers.Max(float64(mo.nbdigits1), float64(mo.nbdigits2))) ||
			mo.nbdigitsanswer > 1+int(helpers.Max(float64(mo.nbdigits1), float64(mo.nbdigits2))) {
			return problemJSON{}, fmt.Errorf("It is not possible to generate a sum with %v digits with %v and %v digits in the first and second operands",
				mo.nbdigitsanswer, mo.nbdigits1, mo.nbdigits2)
		}

	case "-":
		if mo.nbdigitsanswer < 1 ||
			mo.nbdigitsanswer > int(helpers.Max(float64(mo.nbdigits1), float64(mo.nbdigits2))) {
			return problemJSON{}, fmt.Errorf("It is not possible to generate a subtraction with %v digits with %v and %v digits in the first and second operands",
				mo.nbdigitsanswer, mo.nbdigits1, mo.nbdigits2)
		}

	case "*":
		if mo.nbdigitsanswer < mo.nbdigits1+mo.nbdigits2-1 ||
			mo.nbdigitsanswer > mo.nbdigits1+mo.nbdigits2 {
			return problemJSON{}, fmt.Errorf("It is not possible to generate a multiplication with %v digits with %v and %v digits in the first and second operands",
				mo.nbdigitsanswer, mo.nbdigits1, mo.nbdigits2)
		}

	case "/":
		if mo.nbdigitsanswer < 1 ||
			mo.nbdigitsanswer > int(helpers.Max(float64(mo.nbdigits1), float64(mo.nbdigits2)))-
				int(helpers.Min(mo.nbdigits1, mo.nbdigits2)) {
			return problemJSON{}, fmt.Errorf("It is not possible to generate a division with %v digits with %v and %v digits in the first and second operands",
				mo.nbdigitsanswer, mo.nbdigits1, mo.nbdigits2)
		}
	}

	// randomly pick up operands for this instance. Retry as many times as
	// necessary as getting one instance which is compliant with the given
	// parameters
	var operand1, operand2, answer string
	for {

		// create the first operand
		operand1 = ""
		for i := 0; i < mo.nbdigits1; i++ {
			operand1 = operand1 + digits[rand.Intn(len(digits))]
		}

		// create the second operand
		operand2 = ""
		for i := 0; i < mo.nbdigits2; i++ {
			operand2 = operand2 + digits[rand.Intn(len(digits))]
		}

		// compute the answer
		op1, _ := helpers.Atoi(operand1)
		op2, _ := helpers.Atoi(operand2)

		// Intentionally remove combinations of subtractions/divisions where
		// the second operand is greater than the first operand
		if (mo.operator == "-" || mo.operator == "/") &&
			(op2 > op1) {
			continue
		}

		// compute the answer
		switch mo.operator {
		case "+":
			answer = fmt.Sprintf("%v", op1+op2)
		case "-":
			answer = fmt.Sprintf("%v", op1-op2)
		case "*":
			answer = fmt.Sprintf("%v", op1*op2)
		case "/":
			answer = fmt.Sprintf("%v", op1/op2)
		}

		// and verify that an answer with the given number of digits has been
		// generated. If so, exit
		if len(answer) == mo.nbdigitsanswer {
			break
		}
	}

	// -- solution

	// the solution is given as a concatenation of the digits of both operands
	// and the result, but the first four strings have to provide information
	// about the size of the different items of this operation
	solution := make([]string, 4+mo.nbdigits1+mo.nbdigits2+mo.nbdigitsanswer)
	solution[0] = fmt.Sprintf("%v", mo.operator)
	solution[1] = fmt.Sprintf("%v", mo.nbdigits1)
	solution[2] = fmt.Sprintf("%v", mo.nbdigits2)
	solution[3] = fmt.Sprintf("%v", mo.nbdigitsanswer)

	// next, copy all digits of all operands and the result to the solution
	for i := 0; i < mo.nbdigits1; i++ {
		digit, _ := helpers.Atoi(operand1[i])
		solution[4+i] = fmt.Sprintf("%v", digit)
	}
	for i := 0; i < mo.nbdigits2; i++ {
		digit, _ := helpers.Atoi(operand2[i])
		solution[4+len(operand1)+i] = fmt.Sprintf("%v", digit)
	}
	for i := 0; i < mo.nbdigitsanswer; i++ {
		digit, _ := helpers.Atoi(answer[i])
		solution[4+len(operand1)+len(operand2)+i] = fmt.Sprintf("%v", digit)
	}

	// -- args

	// for creating the specific arguments for this problem, randomly mask the
	// specified number of digits in each item. The following vectors contain
	// the positions that have to be masked in each item
	var masked1, masked2, maskedanswer []int
	for {
		idx := rand.Intn(mo.nbdigits1)
		if !helpers.FindInt(idx, masked1) {
			masked1 = append(masked1, idx)
		}
		if len(masked1) == mo.nbmasked1 {
			break
		}
	}
	for {
		idx := rand.Intn(mo.nbdigits2)
		if !helpers.FindInt(idx, masked2) {
			masked2 = append(masked2, idx)
		}
		if len(masked2) == mo.nbmasked2 {
			break
		}
	}
	for {
		idx := rand.Intn(mo.nbdigitsanswer)
		if !helpers.FindInt(idx, maskedanswer) {
			maskedanswer = append(maskedanswer, idx)
		}
		if len(maskedanswer) == mo.nbmaskedanswer {
			break
		}
	}

	// next, copy the solution to the args
	var args []string
	for i := 0; i < len(solution); i++ {

		// if this item has been chosen to be masked, then do show
		if i >= 4 && i < 4+len(operand1) {
			if helpers.FindInt(i-4, masked1) {
				args = append(args, "?")
				continue
			}
		}

		if i >= 4+len(operand1) && i < 4+len(operand1)+len(operand2) {
			if helpers.FindInt(i-4-len(operand1), masked2) {
				args = append(args, "?")
				continue
			}
		}

		if i >= 4+len(operand1)+len(operand2) {
			if helpers.FindInt(i-4-len(operand1)-len(operand2), maskedanswer) {
				args = append(args, "?")
				continue
			}
		}

		// otherwise, copy it from the solution
		args = append(args, solution[i])
	}

	// Now, generate the mystery operation
	return problemJSON{
		Probtype: "MysteryOperation",
		Args:     args,
		Solution: solution,
	}, nil
}

// Local Variables:
// mode:go
// fill-column:80
// End:

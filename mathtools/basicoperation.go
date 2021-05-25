// -*- coding: utf-8 -*-
// basicoperation.go
//
// Description: Provides services for automatically creating a basic operation
// -----------------------------------------------------------------------------
//
// Started on <mar 25-05-2021 20:47:25.993610806 (1621968445)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package mathtools

// constants
// ----------------------------------------------------------------------------

// There are two different types of basic operations: "result" or "operand". In
// the first case, all operands are visible and the student has to provide the
// value of the result; in the latter, the result can be seen but one operand is
// missing whose value has to be guessed by the student
const (
	BORESULT int = iota
	BOOPERAND
)

// types
// ----------------------------------------------------------------------------

// A basic operation consists of a number of operands related to any of the
// operations: +, -, *, / whose number of digits have to be specified as much as
// the number of desired digits in the result. There are two types of basic
// operations:
//
//    0: all operands are given and the student has to guess the result
//    1: all operands but one are shown but the result can be seen. The student
//    has to provide the value of the missing operand
type basicOperation struct {
	botype       int
	operator     string
	nboperands   int
	nbdigitsop   int
	nbdigitsrslt int
}

// Local Variables:
// mode:go
// fill-column:80
// End:

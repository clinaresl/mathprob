/*
  simpleoperation.go
  Description: Provides services for automatically creating a simple algebraic operation
  -----------------------------------------------------------------------------

  Started on  <Wed Jul 12 13:19:25 2017 >
  Last update <>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by
  Login   <clinares@atlas>
*/

package mathtools

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

// constants
// ----------------------------------------------------------------------------

// Different basic operations (represented as integers) are allowed,
// including the following
const (
	ADD int = iota
	SUB
	PROD
	DIV
)

// the TikZ code for generating algebraic operations is shown below
const latexOperationCode = `\begin{minipage}{0.25\linewidth}  
  \begin{center}
    \begin{tikzpicture}

      % show the operands
      {{.GetOperandA}}
      {{.GetOperandB}}

      % show a straight line
      \node [below right=0.35 cm and 0.25 cm of {{.GetLabelB}}] (right) {};
      \node [below left=0.35 cm and 2.05 cm of {{.GetLabelB}}] (left) {};
      \draw [thick] (left) -- (right);

      % show the operator
      {{.GetOperator}}

      % show the result
      {{.GetResult}}

    \end{tikzpicture}
  \end{center}
\end{minipage}
`

// the TikZ code for generating the operands is the following
const latexOperandCode = `\node ({{.GetLabel}}) at {{.GetPosition}} {};
      \node [left = 0.0 cm of {{.GetLabel}}] ({{.GetId}}) {\huge {{.GetValue}}};`

// the TikZ code for generating operators is shown below
const latexOperatorCode = `\node [left = 0.0 cm of num2] ({{.GetLabel}}) {\huge{{.GetSymbol}}};`

// the TikZ code for generating the result is the following
const latexResultCode = `\node ({{.GetLabel}}) at {{.GetPosition}} {};
      \node [rectangle, minimum width={{.GetMinimumWidth}}, minimum height = {{.GetMinimumHeight}}, draw, left = 0.0 cm of {{.GetLabel}}] {};`

// types
// ----------------------------------------------------------------------------

// A LaTeX operand consists of a label, an identifier, a position and an integer
type latexOperand struct {
	label, id string
	pos       position
	value     int
}

// A LaTeX operator consists of a label and a mathematical symbol
type latexOperator struct {
	label, symbol string
}

// A LaTeX result consists of a label, a position, and a minimum width
// and height
type latexResult struct {
	label                       string
	pos                         position
	minimumWidth, minimumHeight float64
}

// A simple operation consists of two operands, an operator and a result
type singleOperation struct {
	operandA, operandB latexOperand
	operator           latexOperator
	result             latexResult
}

// methods
// ----------------------------------------------------------------------------

// -- latexOperand
// ----------------------------------------------------------------------------

// Provide TikZ code to represent operands
func (operand latexOperand) String() string {

	// create a template with the TikZ code for showing operands
	tpl, err := template.New("operand").Parse(latexOperandCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, operand); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Return the label of an operand
func (operand latexOperand) GetLabel() string {
	return fmt.Sprintf("%v", operand.label)
}

// Return the identifier of an operand
func (operand latexOperand) GetId() string {
	return fmt.Sprintf("%v", operand.id)
}

// Return the position of an operand
func (operand latexOperand) GetPosition() string {
	return fmt.Sprintf("%v", operand.pos)
}

// Return the value of an operand
func (operand latexOperand) GetValue() string {
	return fmt.Sprintf("%v", operand.value)
}

// -- latexOperator
// ----------------------------------------------------------------------------

// Provide TikZ code to represent operators
func (operator latexOperator) String() string {

	// create a template with the TikZ code for showing operands
	tpl, err := template.New("operator").Parse(latexOperatorCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, operator); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Return the label of an operator
func (operator latexOperator) GetLabel() string {
	return fmt.Sprintf("%v", operator.label)
}

// Return the symbol of an operator
func (operator latexOperator) GetSymbol() string {
	return fmt.Sprintf("%v", operator.symbol)
}

// -- latexResult
// ----------------------------------------------------------------------------

// Provide TikZ code to represent operands
func (result latexResult) String() string {

	// create a template with the TikZ code for showing the result
	tpl, err := template.New("result").Parse(latexResultCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, result); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// return the label of a result
func (result latexResult) GetLabel() string {
	return fmt.Sprintf("%v", result.label)
}

// return the position of a result
func (result latexResult) GetPosition() string {
	return fmt.Sprintf("%v", result.pos)
}

// return the minimum width of a result in centimeters
func (result latexResult) GetMinimumWidth() string {
	return fmt.Sprintf("%v cm", result.minimumWidth)
}

// return the minimum height of a result in centimeters
func (result latexResult) GetMinimumHeight() string {
	return fmt.Sprintf("%v cm", result.minimumHeight)
}

// -- singleOperation
// ----------------------------------------------------------------------------

// Execute the given operation and returns legal TikZ code to represent it
func (operation singleOperation) Execute() string {

	// create a template with the TikZ code for showing single
	// operations
	tpl, err := template.New("operation").Parse(latexOperationCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, operation); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Return the label of the first operand
func (operation singleOperation) GetLabelA() string {
	return operation.operandA.label
}

// Return the label of the second operand
func (operation singleOperation) GetLabelB() string {
	return operation.operandB.label
}

// Generates TikZ code to draw the first operand
func (operation singleOperation) GetOperandA() string {
	return fmt.Sprintf("%v", operation.operandA)
}

// Generates TikZ code to draw the second operand
func (operation singleOperation) GetOperandB() string {
	return fmt.Sprintf("%v", operation.operandB)
}

// Generates TikZ code to draw the operator
func (operation singleOperation) GetOperator() string {
	return fmt.Sprintf("%v", operation.operator)
}

// Generates TikZ code to draw the result
func (operation singleOperation) GetResult() string {
	return fmt.Sprintf("%v", operation.result)
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

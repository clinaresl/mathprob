/*
  sequence.go
  Description: Provides services for automatically creating a sequence or a part of it
  -----------------------------------------------------------------------------

  Started on  <Wed Dec 12 20:44:45 2018 >
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

// Three different types of sequences are acknowledged: either
// computing the previous, the subsequent, or both
const (
	PREVIOUS   int = 1
	SUBSEQUENT     = 2
	FULL           = 3
)

// the TikZ code for generating sequences type "previous" is shown below
const latexPreviousSequenceCode string = `\begin{minipage}{0.25\linewidth}  
  \begin{center}
    \begin{tikzpicture}

      % show the area to draw the answer
      {{.GetSequenceAnswer}}

      % show the index
      {{.GetIndex}}

    \end{tikzpicture}
  \end{center}
\end{minipage}
`

// the TikZ code for generating sequences type "subsequent" is shown below
const latexSubsequentSequenceCode string = `\begin{minipage}{0.25\linewidth}  
  \begin{center}
    \begin{tikzpicture}

      % show the index
      {{.GetIndex}}

      % show the area to draw the answer
      {{.GetSequenceAnswer}}

    \end{tikzpicture}
  \end{center}
\end{minipage}
`

// the TikZ code for generating the indices is the following
const latexIndexCode string = `\node ({{.GetLabel}}) at {{.GetPosition}} {};
      \node [left = 0.0 cm of {{.GetLabel}}] ({{.GetId}}) {\huge {{.GetValue}}};`

// the TikZ code for generating the result is the following
const latexSequenceAnswerCode string = `\node ({{.GetLabel}}) at {{.GetPosition}} {};
      \node [rectangle, minimum width={{.GetMinimumWidth}} cm, minimum height = {{.GetMinimumHeight}} cm, draw, left = 0.0 cm of {{.GetLabel}}] {};`

// An Index is the number shown whose precendent and/or subsequent has
// to be obtained. It consists of a label, an identifier, a position
// and an integer
type latexIndex struct {
	label, id string
	pos       position
	value     int
}

// A LaTeX Sequence Answer is the area where the student has to write
// her answer. It consists of a label, a position, and a minimum width
// and height
type latexSequenceAnswer struct {
	label                       string
	pos                         position
	minimumWidth, minimumHeight float64
}

// Finally, a problem on sequences is defined as shown next. Each one
// consists of just an index and an answer that might precede or
// follow the index according to the given sequence type.
//
// WARNING - Full sequences not implemented yet!
type sequenceProblem struct {
	index   latexIndex
	answer  latexSequenceAnswer
	seqtype int
}

// methods
// ----------------------------------------------------------------------------

// -- latexIndex
// ----------------------------------------------------------------------------

// Provide TikZ code to represent an index
func (index latexIndex) String() string {

	// create a template with the TikZ code for showing indices
	tpl, err := template.New("index").Parse(latexIndexCode)
	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, index); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Return the label of an index
func (index latexIndex) GetLabel() string {
	return fmt.Sprintf("%v", index.label)
}

// Return the identifier of an index
func (index latexIndex) GetId() string {
	return fmt.Sprintf("%v", index.id)
}

// Return the position of an index
func (index latexIndex) GetPosition() string {
	return fmt.Sprintf("%v", index.pos)
}

// Return the value of an index
func (index latexIndex) GetValue() string {
	return fmt.Sprintf("%v", index.value)
}

// -- latexSequenceAnswer
// ----------------------------------------------------------------------------

// Provide TikZ code to represent operands
func (result latexSequenceAnswer) String() string {

	// create a template with the TikZ code for showing the result
	tpl, err := template.New("result").Parse(latexSequenceAnswerCode)
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
func (result latexSequenceAnswer) GetLabel() string {
	return fmt.Sprintf("%v", result.label)
}

// return the position of a result
func (result latexSequenceAnswer) GetPosition() string {
	return fmt.Sprintf("%v", result.pos)
}

// return the minimum width of a result in centimeters
func (result latexSequenceAnswer) GetMinimumWidth() string {
	return fmt.Sprintf("%v", result.minimumWidth)
}

// return the minimum height of a result in centimeters
func (result latexSequenceAnswer) GetMinimumHeight() string {
	return fmt.Sprintf("%v", result.minimumHeight)
}

// -- sequenceProblem
// ----------------------------------------------------------------------------

// Execute the given sequence and returns legal TikZ code to represent it
func (sequence sequenceProblem) Execute(seqtype int) string {

	// create a template with the TikZ code for showing this
	// sequence according to the given type
	var tpl *template.Template
	var err error
	switch seqtype {
	case PREVIOUS:
		tpl, err = template.New("sequence").Parse(latexPreviousSequenceCode)
	case SUBSEQUENT:
		tpl, err = template.New("sequence").Parse(latexSubsequentSequenceCode)
	case FULL:
		tpl, err = template.New("sequence").Parse(latexPreviousSequenceCode)
	}

	if err != nil {
		log.Fatal(err)
	}

	// and now make the appropriate substitutions. Note that the
	// execution of the template is written to a string
	var tplOutput bytes.Buffer
	if err := tpl.Execute(&tplOutput, sequence); err != nil {
		log.Fatal(err)
	}

	return tplOutput.String() // and return the resulting string
}

// Generates TikZ code to draw the index
func (sequence sequenceProblem) GetIndex() string {
	return fmt.Sprintf("%v", sequence.index)
}

// Generates TikZ code to draw the area to write the answer
func (sequence sequenceProblem) GetSequenceAnswer() string {
	return fmt.Sprintf("%v", sequence.answer)
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */

/*
  components.go
  Description: Definition of different reusable components to be used in TikZ
               drawings
  -----------------------------------------------------------------------------

  Started on  <Mon Jul 10 09:31:10 2017 >
  Last update <>
  -----------------------------------------------------------------------------
  Made by Carlos Linares López
  Login <carlos.linares@uc3m.es>
*/

// This package provides a number of reusable components that can be used for
// creating exercises automatically.
package components

// the different components are distinguished by an integer index
type ComponentId int

// which is used to number all the different reusable components implemented in
// this package
const (
	POSITION ComponentId = iota
	FORMULA
	COORDINATE
	TEXT
	BOX
)

//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

// Implement this interface to make your object connectable by this library
type Connectable interface {
	// GetId should return unique identifier value for each connectable instance
	GetId() string
}

// Implement this interface to make your object connectable objects to connect
type Operator interface {
	// Evaluate should connect connectable a to connectable: a -> b
	Evaluate(a Connectable, b Connectable) error
	// GetId should return unique identifier value for this connection type
	GetId() string
}

// These are the operations you can do with the terms of the category equation
type Category interface {
	GetSources() ConnectableSet
	GetSinks() ConnectableSet
	GetOperator() Operator
	// Category is plainly a container for planned connection operations
	GetOperations() OperationSet
	// returns true, if the both categories involved have the same planned connection operations
	Equals(another Category) bool
	IsZero() bool
	IsIdentity() bool
	// Calls the operator to connect the all the planned connections
	Evaluate() error
	// Calls the operator to connect the all the planned connections in alphabetical order
	EvaluateSorted() error

	String() string
}

type EquationTerm interface {
	Category
	GetProcessedTerm() ProcessedTerm
	// the '+' of the equations
	Add(category Category) EquationTerm
	// the '-' of the equations
	Discard(category Category) EquationTerm
	// the '->' of the equations
	Connect(category Category) EquationTerm
}

// And finally the category equations can be made in
//
// Please do not mix equations done with two different operators, because they might not work well together
type EquationFactory interface {
	// returns the Idetitity term
	I() EquationTerm
	// returns the terminatOr term
	O() EquationTerm
	// wraps your connectable into a equation term
	W(c Connectable) EquationTerm
}

// Equation factory constructor.
//
// Please do not mix equations done with two different operators, because they might not work well together
func NewEquationFactory(operator Operator) EquationFactory {
	return &equationFactory{Operator: operator}
}

// Implementation details

type equationFactory struct {
	Operator Operator
}

func (p *equationFactory) I() EquationTerm {
	return NewIdentityTerm(p.Operator)
}

func (p *equationFactory) O() EquationTerm {
	return NewZeroTerm(p.Operator)
}

func (p *equationFactory) W(c Connectable) EquationTerm {
	return NewWrapperTerm(p.Operator, c)
}

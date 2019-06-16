//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
// @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

// Connectable has to be implemented in order to use this library.
// Connectables are the base items connected by the library.
type Connectable interface {
	// GetId should return unique identifier value for each connectable instance
	GetId() string
}

// Operator has to be implemented in order to use this library. Operator connects the Connactable objects.
type Operator interface {
	// Evaluate should connect connectable a to connectable: a -> b
	Evaluate(a Connectable, b Connectable) error
	// GetId should return unique identifier value for this connection type
	GetId() string
}

// Category is plainly a container for already planned connection operations, sources and sinks
type Category interface {
	// GetSources returns the sources which are outputs from this category to the next one
	GetSources() ConnectableSet
	// GetSinks returns the sinks which are inputs to this category from the previous one
	GetSinks() ConnectableSet
	// GetOperator returns the connecting operator
	GetOperator() Operator
	// GetOperations returns the  already planned connection operations
	GetOperations() OperationSet
	// Equals returns true, if the both categories involved have the same planned connection operations, sources and sinks
	Equals(another Category) bool
	// IsZero indicates if this category bahaves like a terminator term on the equations
	IsZero() bool
	// IsIdentity indicates if this category bahaves like an identity term on the equations
	IsIdentity() bool
	// Evaluate calls the operator to connect the all the planned connections
	Evaluate() error
	// EvaluateSorted calls the operator to connect the all the planned connections in alphabetical order
	EvaluateSorted() error
	// String prints the category in human readable form. Do not use for serialization.
	String() string
}

// EquationTerm describes the arithmetic operations for the category equations
type EquationTerm interface {
	Category
	// GetProcessedTerm returns a description of the processed arithmetic operation which has lead to this term
	GetProcessedTerm() ProcessedTerm
	// Add is the '+' of the equation operations
	Add(category Category) EquationTerm
	// Discard is the '-' of the equation operations
	Discard(category Category) EquationTerm
	// Connect is the '->' of the equation operations
	Connect(category Category) EquationTerm
}

// EquationFactory is finally the place where category equations can be made from
// Please do not mix equations done with two different operators, because they might not work well together
type EquationFactory interface {
	// returns the Idetitity term
	I() EquationTerm
	// returns the terminatOr term
	O() EquationTerm
	// wraps your connectable into a equation term
	W(c Connectable) EquationTerm
}

// NewEquationFactory creates a new Equation factory
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

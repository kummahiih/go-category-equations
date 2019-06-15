//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
	"fmt"
)

type Connectable interface {
	GetId() string
}

type ConnectableSet interface {
	Union(another ConnectableSet) ConnectableSet
	DiscardAll(another ConnectableSet) ConnectableSet
	Clone() ConnectableSet
	Add(f Connectable)
	Remove(f Connectable)
	Equals(another ConnectableSet) bool
	AsArray() []Connectable
	String() string
}

type Operator interface {
	Evaluate(a Connectable, b Connectable) error
	GetId() string
}

type FreezedOperation interface {
	GetSource() Connectable
	GetSink() Connectable
	Evaluate() error
	Equals(another FreezedOperation) bool
	GetOperator() Operator
}

type OperationSet interface {
	Union(another OperationSet) OperationSet
	DiscardAll(another OperationSet) OperationSet
	Clone() OperationSet
	//inplce add
	Add(f FreezedOperation)
	//inplace remove
	Remove(f FreezedOperation)

	Equals(another OperationSet) bool

	AsArray() []FreezedOperation
	GetOperator() Operator
}

type Category interface {
	GetSources() ConnectableSet
	GetSinks() ConnectableSet
	GetOperator() Operator
	GetOperations() OperationSet
	Equals(another Category) bool
	IsZero() bool
	IsIdentity() bool
	Evaluate() error
	String() string
}

type Operation int

const (
	ADD Operation = iota
	DISCARD
	ARROW
)

type ProcessedTerm interface {
	GetSink() Category
	GetOperation() Operation
	GetSource() Category
	Equals(another ProcessedTerm) bool
	String() string
}

type EquationTerm interface {
	Category
	GetProcessedTerm() ProcessedTerm
	Add(category Category) EquationTerm
	Discard(category Category) EquationTerm
	Connect(category Category) EquationTerm
}

// Helper functions

func O2S(op Operation) string {
	switch op {
	case ADD:
		return "+"
	case DISCARD:
		return "-"
	case ARROW:
		return "*"
	}
	panic("invalid operation")
}

func EqualOperators(f Operator, another Operator) bool {
	return f.GetId() == another.GetId()
}

func CompatibleOperators(f Operator, another Operator) error {
	if !EqualOperators(f, another) {
		return fmt.Errorf("Expected operator %s, got %s", f.GetId(), another.GetId())
	}
	return nil
}

// Factories

func NewConnectableSet() ConnectableSet {
	return newConnectableSetFromArray([]Connectable{})
}

func NewFreezedOperation(operator Operator, source Connectable, sink Connectable) FreezedOperation {
	return &freezedOperation{Source: source, Sink: sink, Operator: operator}
}

func NewOperationSet(operator Operator) OperationSet {
	return &operationSet{
		Operator:          operator,
		FreezedOperations: make(map[FreezedOperation]bool),
	}
}

func NewProcessedTerm(source Category, operation Operation, sink Category) ProcessedTerm {
	return &processedTerm{
		Source:    source,
		Sink:      sink,
		Operation: operation}
}

func NewIdentityTerm(operator Operator) EquationTerm {
	return &equationTerm{
		categoryImpl: categoryImpl{
			Sources:    NewConnectableSet(),
			Sinks:      NewConnectableSet(),
			Operator:   operator,
			Operations: NewOperationSet(operator),
			isZero:     false,
			isIdentity: true,
			stringImpl: func(c *categoryImpl) string { return "I" }},
		processedTerm: nil}
}

func NewZeroTerm(operator Operator) EquationTerm {
	return &equationTerm{
		categoryImpl: categoryImpl{
			Sources:    NewConnectableSet(),
			Sinks:      NewConnectableSet(),
			Operator:   operator,
			Operations: NewOperationSet(operator),
			isZero:     true,
			isIdentity: false,
			stringImpl: func(c *categoryImpl) string { return "O" }},
		processedTerm: nil}
}

func NewWrapperTerm(operator Operator, connectable Connectable) EquationTerm {
	sources := NewConnectableSet()
	sources.Add(connectable)

	sinks := NewConnectableSet()
	sinks.Add(connectable)

	return &equationTerm{
		categoryImpl: categoryImpl{
			Sources:    sources,
			Sinks:      sinks,
			Operator:   operator,
			Operations: NewOperationSet(operator),
			isZero:     false,
			isIdentity: false,
			stringImpl: func(c *categoryImpl) string { return connectable.GetId() }},
		processedTerm: nil}
}

func NewIntermediateTerm(
	operator Operator,
	sources ConnectableSet,
	sinks ConnectableSet,
	operations OperationSet,
	processedTerm ProcessedTerm) EquationTerm {
	return &equationTerm{
		categoryImpl: categoryImpl{
			Sources:    sources,
			Sinks:      sinks,
			Operator:   operator,
			Operations: operations,
			isZero:     false,
			isIdentity: false,
			stringImpl: func(c *categoryImpl) string { return processedTerm.String() }},
		processedTerm: processedTerm}
}

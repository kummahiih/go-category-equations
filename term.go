//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
// @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
	"fmt"
)

// ProcessedTerm contains a description of the processed arithmetic operation which has lead to some term
type ProcessedTerm interface {
	// GetSink returns the sink of the done operation
	GetSink() Category
	// GetOperation returns the done operation
	GetOperation() Operation
	// GetSource returns the source of the done operation
	GetSource() Category
	// Equals returns true, if the operation had the same parameters
	Equals(another ProcessedTerm) bool
	// String returns a human readable description of the done operation
	String() string
}

// Operation is used to identify the different operations one can make to the category equation terms
type Operation int

const (
	// ADD denotes the Union -operation
	ADD Operation = iota
	// DISCARD denotes the DiscardAll  -operation
	DISCARD
	// ARROW denotes the Connect -operation
	ARROW
)

// Helper functions for Operation

// O2S returns a string representetion of the given operation
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

// EqualOperators returns true, if the Ids are same
func EqualOperators(f Operator, another Operator) bool {
	return f.GetId() == another.GetId()
}

// CompatibleOperators returns an error, if the connection operators are not same
func CompatibleOperators(f Operator, another Operator) error {
	if !EqualOperators(f, another) {
		return fmt.Errorf("Expected operator %s, got %s", f.GetId(), another.GetId())
	}
	return nil
}

// NewProcessedTerm returns a new ProcessedTerm instance
func NewProcessedTerm(source Category, operation Operation, sink Category) ProcessedTerm {
	return &processedTerm{
		Source:    source,
		Sink:      sink,
		Operation: operation}
}

// NewIdentityTerm returns a new identity term instance
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

// NewZeroTerm returns a new zero term instance
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

// NewWrapperTerm returns a new Connectable wrapping term instance
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

// NewIntermediateTerm returns a new intermediate term. Used on the arithmetic operation implementations
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
			isIdentity: processedTerm != nil && processedTerm.GetOperation() == ADD && (processedTerm.GetSink().IsIdentity() || processedTerm.GetSource().IsIdentity()),
			stringImpl: func(c *categoryImpl) string { return processedTerm.String() }},
		processedTerm: processedTerm}
}

// implementation details

type processedTerm struct {
	Sink      Category
	Source    Category
	Operation Operation
}

func (p *processedTerm) GetSink() Category {
	return p.Sink
}

func (p *processedTerm) GetOperation() Operation {
	return p.Operation
}

func (p *processedTerm) GetSource() Category {
	return p.Source

}

func (p *processedTerm) Equals(another ProcessedTerm) bool {
	return p.Sink.Equals(another.GetSink()) && p.Source.Equals(another.GetSource()) && p.Operation == another.GetOperation()
}

func (p *processedTerm) String() string {
	return fmt.Sprintf("(%s) %s (%s)", p.Source.String(), O2S(p.Operation), p.Sink.String())
}

type equationTerm struct {
	categoryImpl
	processedTerm ProcessedTerm
}

func (e *equationTerm) GetProcessedTerm() ProcessedTerm {
	return e.processedTerm
}

func (e *equationTerm) Add(category Category) EquationTerm {
	return NewIntermediateTerm(
		e.Operator,
		e.Sources.Union(category.GetSources()),
		e.Sinks.Union(category.GetSinks()),
		e.Operations.Union(category.GetOperations()),
		NewProcessedTerm(e, ADD, category))
}

func (e *equationTerm) Discard(category Category) EquationTerm {
	return NewIntermediateTerm(
		e.Operator,
		e.Sources.DiscardAll(category.GetSources()),
		e.Sinks.DiscardAll(category.GetSinks()),
		e.Operations.DiscardAll(category.GetOperations()),
		NewProcessedTerm(e, DISCARD, category))
}

func (e *equationTerm) Connect(anext Category) EquationTerm {
	if e.IsZero() {
		return NewIntermediateTerm(
			anext.GetOperator(),
			anext.GetSources().Clone(),
			NewConnectableSet(),
			anext.GetOperations().Clone(),
			NewProcessedTerm(e, ARROW, anext))
	}
	if anext.IsZero() {
		return NewIntermediateTerm(
			anext.GetOperator(),
			NewConnectableSet(),
			e.GetSinks().Clone(),
			e.Operations.Clone(),
			NewProcessedTerm(e, ARROW, anext))
	}

	newOperations := NewOperationSet(e.Operator)

	for _, source := range e.Sources.AsArray() {
		for _, sink := range anext.GetSinks().AsArray() {
			newOperations.Add(NewFreezedOperation(e.Operator, source, sink))
		}
	}

	newSources := NewConnectableSet()
	for _, source := range anext.GetSources().AsArray() {
		newSources.Add(source)
	}
	if anext.IsIdentity() { // a -> (I+b)
		for _, source := range e.GetSources().AsArray() {
			newSources.Add(source)
		}
	}

	newSinks := NewConnectableSet()
	for _, sink := range e.GetSinks().AsArray() {
		newSinks.Add(sink)
	}

	if e.IsIdentity() { // (a+I) -> b
		for _, sink := range anext.GetSinks().AsArray() {
			newSinks.Add(sink)
		}
	}

	operations := e.Operations.Union(anext.GetOperations()).Union(newOperations)

	return NewIntermediateTerm(
		anext.GetOperator(),
		newSources,
		newSinks,
		operations,
		NewProcessedTerm(e, ARROW, anext))
}

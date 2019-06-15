package category

import (
	"fmt"
)

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
	if e.IsIdentity() {
		return NewIntermediateTerm(
			anext.GetOperator(),
			anext.GetSources().Clone(),
			anext.GetSinks().Clone(),
			anext.GetOperations().Clone(),
			NewProcessedTerm(e, ARROW, anext))
	}
	if e.IsZero() {
		return NewIntermediateTerm(
			anext.GetOperator(),
			anext.GetSources().Clone(),
			NewConnectableSet(),
			anext.GetOperations().Clone(),
			NewProcessedTerm(e, ARROW, anext))
	}
	if anext.IsIdentity() {
		return NewIntermediateTerm(
			anext.GetOperator(),
			e.GetSources().Clone(),
			e.GetSinks().Clone(),
			e.Operations.Clone(),
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

	newSinks := NewConnectableSet()
	for _, sink := range e.GetSinks().AsArray() {
		newSinks.Add(sink)
	}

	operations := e.Operations.Union(anext.GetOperations()).Union(newOperations)

	return NewIntermediateTerm(
		anext.GetOperator(),
		newSources,
		newSinks,
		operations,
		NewProcessedTerm(e, ARROW, anext))
}

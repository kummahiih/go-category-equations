//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
//"fmt"
)

type freezedOperation struct {
	Source   Connectable
	Sink     Connectable
	Operator Operator
}

func (f *freezedOperation) GetSink() Connectable   { return f.Sink }
func (f *freezedOperation) GetSource() Connectable { return f.Source }

func (f *freezedOperation) Evaluate() error {
	return f.Operator.Evaluate(f.Source, f.Sink)
}

func (f *freezedOperation) Equals(another FreezedOperation) bool {
	return another != nil && f.Source == another.GetSource() && f.Sink == another.GetSink() && f.Operator.GetId() == another.GetOperator().GetId()
}

func (f *freezedOperation) GetOperator() Operator { return f.Operator }

type operationSet struct {
	FreezedOperations map[FreezedOperation]bool
	Operator          Operator
}

func newOperationSetFromArray(operator Operator, operations []FreezedOperation) *operationSet {
	aSet := &operationSet{
		Operator:          operator,
		FreezedOperations: make(map[FreezedOperation]bool, len(operations)),
	}
	for _, v := range operations {
		aSet.FreezedOperations[v] = true
	}
	return aSet
}

func (fs *operationSet) Union(another OperationSet) OperationSet {
	operations := another.AsArray()

	for k, _ := range fs.FreezedOperations {
		operations = append(operations, k)
	}

	unionSet := newOperationSetFromArray(fs.Operator, operations)

	return unionSet
}

func (fs *operationSet) DiscardAll(another OperationSet) OperationSet {

	discardSet := newOperationSetFromArray(fs.Operator, another.AsArray())

	for k, _ := range fs.FreezedOperations {
		discardSet.Remove(k)
	}
	return discardSet
}

func (fs *operationSet) Clone() OperationSet {
	freezeds := make(map[FreezedOperation]bool, len(fs.FreezedOperations))

	for k, v := range fs.FreezedOperations {
		freezeds[k] = v
	}

	return &operationSet{Operator: fs.Operator, FreezedOperations: freezeds}
}

func (fs *operationSet) Add(f FreezedOperation) {
	fs.FreezedOperations[f] = true
}

func (fs *operationSet) Remove(f FreezedOperation) {
	delete(fs.FreezedOperations, f)
}

func (fs *operationSet) Equals(another OperationSet) bool {
	if !EqualOperators(fs.Operator, another.GetOperator()) {
		return false
	}

	operations := another.AsArray()

	if len(operations) != len(fs.FreezedOperations) {
		return false
	}

	for _, f := range operations {
		_, found := fs.FreezedOperations[f]
		if !found {
			return false
		}
	}
	return true
}

func (fs *operationSet) AsArray() []FreezedOperation {
	values := make([]FreezedOperation, len(fs.FreezedOperations))
	i := 0
	for k, _ := range fs.FreezedOperations {
		values[i] = k
		i++
	}
	return values
}

func (fs *operationSet) GetOperator() Operator {
	return fs.Operator
}

func (fs *operationSet) CompatibleWithSet(another OperationSet) error {
	return CompatibleOperators(fs.GetOperator(), another.GetOperator())
}

func (fs *operationSet) CompatibleWithList(another []FreezedOperation) error {
	for f := range another {
		err := CompatibleOperators(fs.GetOperator(), another[f].GetOperator())
		if err != nil {
			return err
		}
	}
	return nil
}

//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
	"fmt"
	"sort"
	"strings"
)

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
	AsSortedArray() []FreezedOperation
	GetOperator() Operator
}

func NewFreezedOperation(operator Operator, source Connectable, sink Connectable) FreezedOperation {
	return &freezedOperation{Source: source, Sink: sink, Operator: operator}
}

func NewOperationSet(operator Operator) OperationSet {
	return &operationSet{
		Operator:          operator,
		FreezedOperations: make(map[freezedOperationKey]FreezedOperation),
	}
}

// implementation details

type freezedOperationKey struct {
	Source   string
	Sink     string
	Operator string
}

func getKey(op FreezedOperation) freezedOperationKey {
	return freezedOperationKey{
		Source:   op.GetSource().GetId(),
		Sink:     op.GetSink().GetId(),
		Operator: op.GetOperator().GetId()}
}

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
	FreezedOperations map[freezedOperationKey]FreezedOperation
	Operator          Operator
}

func newOperationSetFromArray(operator Operator, operations []FreezedOperation) *operationSet {
	aSet := &operationSet{
		Operator:          operator,
		FreezedOperations: make(map[freezedOperationKey]FreezedOperation, len(operations)),
	}
	for _, v := range operations {
		aSet.FreezedOperations[getKey(v)] = v
	}
	return aSet
}

func (fs *operationSet) Union(another OperationSet) OperationSet {
	operations := another.AsArray()

	for _, v := range fs.FreezedOperations {
		operations = append(operations, v)
	}

	unionSet := newOperationSetFromArray(fs.Operator, operations)

	return unionSet
}

func (fs *operationSet) DiscardAll(another OperationSet) OperationSet {

	discardSet := fs.Clone()

	for _, v := range another.AsArray() {
		discardSet.Remove(v)
	}
	return discardSet
}

func (fs *operationSet) Clone() OperationSet {
	freezeds := make(map[freezedOperationKey]FreezedOperation, len(fs.FreezedOperations))

	for k, v := range fs.FreezedOperations {
		freezeds[k] = v
	}

	return &operationSet{Operator: fs.Operator, FreezedOperations: freezeds}
}

func (fs *operationSet) Add(f FreezedOperation) {
	fs.FreezedOperations[getKey(f)] = f
}

func (fs *operationSet) Remove(f FreezedOperation) {
	fmt.Printf("%+v %+v\n", fs.FreezedOperations, f)
	delete(fs.FreezedOperations, getKey(f))
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
		_, found := fs.FreezedOperations[getKey(f)]
		if !found {
			return false
		}
	}
	return true
}

func (fs *operationSet) AsArray() []FreezedOperation {
	values := make([]FreezedOperation, len(fs.FreezedOperations))
	i := 0
	for _, v := range fs.FreezedOperations {
		values[i] = v
		i++
	}
	return values
}

func (fs *operationSet) AsSortedArray() []FreezedOperation {
	arr := fs.AsArray()
	sort.Slice(arr, func(i, j int) bool {
		cmpOp := strings.Compare(arr[i].GetOperator().GetId(), arr[j].GetOperator().GetId())
		cmpSource := strings.Compare(arr[i].GetSource().GetId(), arr[j].GetSource().GetId())
		cmpSink := strings.Compare(arr[i].GetSink().GetId(), arr[j].GetSink().GetId())
		return (cmpOp < 0) || (cmpOp == 0 && cmpSource < 0) || (cmpOp == 0 && cmpSource == 0 && cmpSink < 0)
	})

	return arr
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

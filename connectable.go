//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
// @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
	"sort"
	"strings"
)

// ConnectableSet contains a set of operations and a set of set operations for the connectable set
type ConnectableSet interface {
	// Union is a set union. Returns a new set
	Union(another ConnectableSet) ConnectableSet
	// DiscardAll removes all instances of the another set from this one and returns it as a new
	DiscardAll(another ConnectableSet) ConnectableSet
	// Clone clones the set
	Clone() ConnectableSet
	// Add adds an item to this set
	Add(f Connectable)
	// Remove removes and item from this set
	Remove(f Connectable)
	// Equals is a set equality check
	Equals(another ConnectableSet) bool
	// AsArray returns the operations as an array
	AsArray() []Connectable
	// AsSortedArray returns the operations as a sorted array
	AsSortedArray() []Connectable
	// String prints the content of the ConnectableSet in an human understable form. Don't use for serialization
	String() string
}

// NewConnectableSet creates a new ConnectableSet instance
func NewConnectableSet() ConnectableSet {
	return newConnectableSetFromArray([]Connectable{})
}

// implementation details

type connectableSet struct {
	Connectables map[string]Connectable
}

func newConnectableSetFromArray(operations []Connectable) *connectableSet {
	aSet := &connectableSet{
		Connectables: make(map[string]Connectable, len(operations)),
	}
	for _, v := range operations {
		aSet.Connectables[v.GetId()] = v
	}
	return aSet
}

func (fs *connectableSet) Union(another ConnectableSet) ConnectableSet {
	operations := another.AsArray()

	for _, v := range fs.Connectables {
		operations = append(operations, v)
	}

	unionSet := newConnectableSetFromArray(operations)

	return unionSet
}

func (fs *connectableSet) DiscardAll(another ConnectableSet) ConnectableSet {

	discardSet := fs.Clone()

	for _, v := range another.AsArray() {
		discardSet.Remove(v)
	}
	return discardSet
}

func (fs *connectableSet) Clone() ConnectableSet {
	freezeds := make(map[string]Connectable, len(fs.Connectables))

	for k, v := range fs.Connectables {
		freezeds[k] = v
	}

	return &connectableSet{Connectables: freezeds}
}

func (fs *connectableSet) Add(f Connectable) {
	fs.Connectables[f.GetId()] = f
}

func (fs *connectableSet) Remove(f Connectable) {
	delete(fs.Connectables, f.GetId())
}

func (fs *connectableSet) Equals(another ConnectableSet) bool {
	operations := another.AsArray()

	if len(operations) != len(fs.Connectables) {
		return false
	}

	for _, f := range operations {
		_, found := fs.Connectables[f.GetId()]
		if !found {
			return false
		}
	}
	return true
}

func (fs *connectableSet) AsArray() []Connectable {
	values := make([]Connectable, len(fs.Connectables))
	i := 0
	for _, v := range fs.Connectables {
		values[i] = v
		i++
	}
	return values
}

func (fs *connectableSet) AsSortedArray() []Connectable {
	arr := fs.AsArray()
	sort.Slice(arr, func(i, j int) bool {
		cmpOp := strings.Compare(arr[i].GetId(), arr[j].GetId())
		return cmpOp < 0
	})

	return arr
}

func (fs *connectableSet) String() string {
	ret := ""
	for _, v := range fs.AsSortedArray() {
		ret = ret + v.GetId() + ","
	}
	return ret
}

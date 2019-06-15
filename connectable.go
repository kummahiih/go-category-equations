//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

import (
//"fmt"
)

type connectableSet struct {
	Connectables map[Connectable]bool
}

func newConnectableSetFromArray(operations []Connectable) *connectableSet {
	aSet := &connectableSet{
		Connectables: make(map[Connectable]bool, len(operations)),
	}
	for _, v := range operations {
		aSet.Connectables[v] = true
	}
	return aSet
}

func (fs *connectableSet) Union(another ConnectableSet) ConnectableSet {
	operations := another.AsArray()

	for k, _ := range fs.Connectables {
		operations = append(operations, k)
	}

	unionSet := newConnectableSetFromArray(operations)

	return unionSet
}

func (fs *connectableSet) DiscardAll(another ConnectableSet) ConnectableSet {

	discardSet := newConnectableSetFromArray(another.AsArray())

	for k, _ := range fs.Connectables {
		discardSet.Remove(k)
	}
	return discardSet
}

func (fs *connectableSet) Clone() ConnectableSet {
	freezeds := make(map[Connectable]bool, len(fs.Connectables))

	for k, v := range fs.Connectables {
		freezeds[k] = v
	}

	return &connectableSet{Connectables: freezeds}
}

func (fs *connectableSet) Add(f Connectable) {
	fs.Connectables[f] = true
}

func (fs *connectableSet) Remove(f Connectable) {
	delete(fs.Connectables, f)
}

func (fs *connectableSet) Equals(another ConnectableSet) bool {
	operations := another.AsArray()

	if len(operations) != len(fs.Connectables) {
		return false
	}

	for _, f := range operations {
		_, found := fs.Connectables[f]
		if !found {
			return false
		}
	}
	return true
}

func (fs *connectableSet) AsArray() []Connectable {
	values := make([]Connectable, len(fs.Connectables))
	i := 0
	for k, _ := range fs.Connectables {
		values[i] = k
		i++
	}
	return values
}

func (fs *connectableSet) String() string {
	ret := ""
	for k, _ := range fs.Connectables {
		ret = ret + k.GetId() + ","
	}
	return ret
}

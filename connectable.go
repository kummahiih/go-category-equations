//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

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

func (fs *connectableSet) String() string {
	ret := ""
	for k, _ := range fs.Connectables {
		ret = ret + k + ","
	}
	return ret
}

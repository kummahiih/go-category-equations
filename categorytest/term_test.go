//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//
package categorytest

import (
	"category"
	"testing"
)

func TestConnect(t *testing.T) {
	connect := PrintConnection
	a := NewConnectable("a")
	b := NewConnectable("b")

	res := connect(a, b)
	if res != nil {
		t.Fatalf("wtf")
	}
}

func TestIdentity(t *testing.T) {
	connect := NewConnectionPrinter()
	a := category.NewWrapperTerm(connect, NewConnectable("a"))
	t.Log(a.String())

	I := category.NewIdentityTerm(connect)
	t.Log(I.String())

	k := a.Connect(I)
	t.Log(k.String())
	res := k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
	t.Log("sources")
	t.Log(k.GetSources().String())
	if k.GetSources().String() != "a," {
		t.Fatalf("!")
	}
	t.Log("sinks")
	t.Log(k.GetSinks().String())
	if k.GetSinks().String() != "a," {
		t.Fatalf("!")
	}

	k = I.Connect(a)
	t.Log(k.String())
	res = k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
	t.Log("sources")
	t.Log(k.GetSources().String())
	if k.GetSources().String() != "a," {
		t.Fatalf("!")
	}
	t.Log("sinks")
	t.Log(k.GetSinks().String())
	if k.GetSinks().String() != "a," {
		t.Fatalf("!")
	}
}

func TestZero(t *testing.T) {
	connect := NewConnectionPrinter()
	a := category.NewWrapperTerm(connect, NewConnectable("a"))
	t.Log(a.String())

	O := category.NewZeroTerm(connect)
	t.Log(O.String())

	k := a.Connect(O)
	t.Log(k.String())
	res := k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
	t.Log("sources")
	t.Log(k.GetSources().String())
	if k.GetSources().String() != "" {
		t.Fatalf("!")
	}
	t.Log("sinks")
	t.Log(k.GetSinks().String())
	if k.GetSinks().String() != "a," {
		t.Fatalf("!")
	}

	k = O.Connect(a)
	t.Log(k.String())
	res = k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
	t.Log("sources")
	t.Log(k.GetSources().String())
	if k.GetSources().String() != "a," {
		t.Fatalf("!")
	}
	t.Log("sinks")
	t.Log(k.GetSinks().String())
	if k.GetSinks().String() != "" {
		t.Fatalf("!")
	}

}

func TestZeroOne(t *testing.T) {
	connect := NewConnectionPrinter()
	O := category.NewZeroTerm(connect)
	I := category.NewIdentityTerm(connect)

	k := I.Connect(O)
	t.Log(k.String())
	res := k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
	k = O.Connect(I)
	t.Log(k.String())
	res = k.Evaluate()
	if res != nil {
		t.Fatalf("wtf")
	}
}

func TestZeroOneWrapperEqual(t *testing.T) {
	connect := NewConnectionPrinter()
	O := category.NewZeroTerm(connect)
	I := category.NewIdentityTerm(connect)
	a := category.NewWrapperTerm(connect, NewConnectable("a"))
	t.Log(a.Equals(a))
	if !a.Equals(a) {
		t.Fatalf("!")
	}
	t.Log(a.Equals(I))
	if a.Equals(I) {
		t.Fatalf("!")
	}

	t.Log(a.Equals(O))
	if a.Equals(O) {
		t.Fatalf("!")
	}

	t.Log(I.Equals(I))
	if !I.Equals(I) {
		t.Fatalf("!")
	}

	t.Log(I.Equals(O))
	if !I.Equals(O) {
		t.Fatalf("!")
	}

	t.Log(O.Equals(O))
	if !O.Equals(O) {
		t.Fatalf("!")
	}
}

func TestTermMultiply(t *testing.T) {
	connect := NewConnectionPrinter()
	a := category.NewWrapperTerm(connect, NewConnectable("a"))
	b := category.NewWrapperTerm(connect, NewConnectable("b"))
	c := category.NewWrapperTerm(connect, NewConnectable("c"))
	d := category.NewWrapperTerm(connect, NewConnectable("d"))

	k := a.Add(b).Connect(c.Add(d))
	t.Log(k.String())
	err := k.Evaluate()
	if err != nil {
		t.Fatalf("wtf")
	}

	operations := k.GetOperations()
	if len(operations.AsArray()) != 4 {
		t.Fatalf("!!")
	}

	kk := b.Add(a).Connect(d.Add(c))
	_ = kk.Evaluate()

	if !k.Equals(kk) {
		t.Fatalf("!!")
	}

	aa := a.Add(a).Connect(d.Add(c))
	_ = aa.Evaluate()

	if k.Equals(aa) {
		t.Fatalf("!!")
	}
}

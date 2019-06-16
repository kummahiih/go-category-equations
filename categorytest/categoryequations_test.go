//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//
package categorytest

import (
	"category"
	"fmt"
	"testing"
)

func TestEquationAddOrder(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))

	first := a.Add(b).Connect(c.Add(d))
	second := b.Add(a).Connect(d.Add(c))

	t.Log(first.String())
	t.Log(second.String())

	if !first.Equals(second) {
		t.Fatalf("adding order problem")
	}
	t.Log("are equal")

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestEquationMultiplicationOrder(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))

	first := a.Add(b).Connect(c.Add(d))
	second := c.Add(d).Connect(a.Add(b))

	t.Log(first.String())
	t.Log(second.String())

	if first.Equals(second) {
		t.Fatalf("multiplication order problem")
	}
	t.Log("are not equal")

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestEquationDiscard(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))
	O := G.O()

	first := c.Add(d).Connect(a.Add(b))
	second := c.Add(d).Connect(a.Add(b)).Discard(d.Connect(a))
	third := c.Connect(a.Add(b)).Add(d.Connect(b))
	fourth := third.Discard(
		d.Connect(O).Add(O.Connect(a)))

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")
	err = third.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")
	err = fourth.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(first.String())
	t.Log(second.String())

	if first.Equals(second) {
		t.Fatalf("discard problem")
	}
	t.Log("are not equal")

	t.Log(second.String())
	t.Log(third.String())

	if second.Equals(third) {
		t.Fatalf("discard problem")
	}
	t.Log("are not equal")

	t.Log(second.String())
	t.Log(fourth.String())

	if !second.Equals(fourth) {
		t.Fatalf("discard problem")
	}
	t.Log("are equal")

}

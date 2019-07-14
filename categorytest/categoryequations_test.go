//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
// @license: MIT <http://www.opensource.org/licenses/mit-license.php>
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

func TestEquationSinksAndSources(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))
	//(C(3,4) * C(1,2) - C(4) * C(1)) == C(3) * C(1,2) + C(4) * C(2)
	first := c.Add(d).Connect(a.Add(b)).Discard(d.Connect(a))
	second := c.Connect(a.Add(b)).Add(d.Connect(b))

	t.Log(first.String())
	t.Log(second.String())

	if first.Equals(second) {
		t.Fatalf("sinks problem")
	}

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(first.GetSinks())
	t.Log(second.GetSinks())

	if first.GetSinks().String() != "c," {
		t.Fatalf("sinks problem")
	}

	if second.GetSinks().String() != "c,d," {
		t.Fatalf("sinks problem")
	}
}

func TestEquationIdentityMultiplicationWith3Terms(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	I := G.I()

	first := a.Connect(I).Connect(b)
	second := a.Connect(b)

	t.Log(first.String())
	t.Log(second.String())

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !first.Equals(second) {
		t.Fatalf("identity problem")
	}
}

func TestEquationIdentityMultiplicationWithAddition(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	//c := G.W(NewConnectable("c"))
	I := G.I()
	//O := G.O()

	first := a.Connect(b.Add(I)) //.Connect(c.Add(I))
	t.Log(first.String())
	t.Log(first.GetSinks())
	t.Log(first.GetSources())

	second := a.Connect(b).Add(a)
	t.Log(second.String())
	t.Log(second.GetSinks())
	t.Log(second.GetSources())

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !first.Equals(second) {
		t.Fatalf("identity problem")
	}
}

func TestEquationIdentityMultiplication(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))
	e := G.W(NewConnectable("e"))
	I := G.I()

	//       C(1,2)     *     C(3,4)        * C(5) + C(1,2) *          C(5)
	first := a.Add(b).Connect(c.Add(d)).Connect(e).Add(a.Add(b).Connect(e))
	t.Log(first.String())
	t.Log(first.GetSinks())
	t.Log(first.GetSources())

	// C(1,2) *            ( C(3,4) + I ) * C(5)
	second := a.Add(b).Connect(c.Add(d).Add(I)).Connect(e)
	t.Log(second.String())
	t.Log(second.GetSinks())
	t.Log(second.GetSources())

	err := first.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("vs\n")

	err = second.EvaluateSorted()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !first.Equals(second) {
		t.Fatalf("identity problem")
	}
}

func TestExampleForDocumentation(t *testing.T) {
	connect := NewConnectionPrinter()
	G := category.NewEquationFactory(connect)

	a := G.W(NewConnectable("a"))
	b := G.W(NewConnectable("b"))
	c := G.W(NewConnectable("c"))
	d := G.W(NewConnectable("d"))
	e := G.W(NewConnectable("e"))
	I := G.I()
	O := G.O()

	first := O.Connect(
		a.Add(b).Connect(c.Add(d)).Connect(e).Add(a.Add(b).Connect(e))).Connect(O)
	simplified := O.Connect(
		a.Add(b).Connect(c.Add(d).Add(I)).Connect(e)).Connect(O)

	if first.Equals(simplified) {
		fmt.Printf("Equations %s and %s are equal.\n", first, simplified)
	}

	fmt.Printf("Connections:\n")

	err := first.EvaluateSorted()
	if err != nil {
		fmt.Printf("something bad happened")
	}

	fmt.Printf("sinks: %s sources: %s\n", first.GetSinks(), first.GetSources())

	fmt.Printf("vs:\n")

	err = simplified.EvaluateSorted()
	if err != nil {
		fmt.Printf("something bad happened")
	}

	fmt.Printf("sinks: %s sources: %s\n", simplified.GetSinks(), simplified.GetSources())

}

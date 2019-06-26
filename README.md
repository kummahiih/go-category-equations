# go-category-equations

Rewrite of the https://github.com/kummahiih/python-category-equations -library in Golang. 
The library should be fluent to use and have the same tests or more. 
The full API documentation can be found from https://godoc.org/github.com/kummahiih/go-category-equations .

With the tools provided here you can create category like equations for the given operator. On the equations the underlaying 'Add' ('+' on prints) and 'Discard' ('-' on prints) operations are basic set operations called union and discard and the 'Connect' ('*' on prints) operator connects sources to sinks. The equation system also has a Identity 'I' term and zerO -like termination term 'O'. For futher details go https://en.wikipedia.org/wiki/Category_(mathematics)#Definition

## Usage Example

Implement 'category.Connectable' (here 'NewConnectable') and 'category.Operator':

	// Connectable has to be implemented in order to use this library.
	// Connectables are the base items connected by the library.
	type Connectable interface {
		// GetId should return unique identifier value for each connectable instance
		GetId() string
	}

	// Operator has to be implemented in order to use this library. Operator connects the Connactable objects.
	type Operator interface {
		// Evaluate should connect connectable a to connectable: a -> b
		Evaluate(a Connectable, b Connectable) error
		// GetId should return unique identifier value for this connection type
		GetId() string
	}


This has been done for tests in 'category/categorytest/testprinter.go' ('NewConnectable' and 'NewConnectionPrinter') and their usage can be seen for example from a test 'TestExampleForDocumentation' which lies in the file 'category/categorytest/categoryequations_test.go':

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
		a.Add(b).Connect( c.Add(d).Add(I)).Connect(e)).Connect(O)

	if (first.Equals(simplified)) {
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

Which prints (when ran )

    go/src/category/categorytest$ go test -v -run TestExampleForDocumentation
    === RUN   TestExampleForDocumentation
    Equations ((O) * (((((a) + (b)) * ((c) + (d))) * (e)) + (((a) + (b)) * (e)))) * (O) and ((O) * ((((a) + (b)) * (((c) + (d)) + (I))) * (e))) * (O) are equal.
    Connections:
    a -> c
    a -> d
    a -> e
    b -> c
    b -> d
    b -> e
    c -> e
    d -> e
    sinks:  sources: 
    vs:
    a -> c
    a -> d
    a -> e
    b -> c
    b -> d
    b -> e
    c -> e
    d -> e
    sinks:  sources: 
    --- PASS: TestExampleForDocumentation (0.00s)
    PASS
    ok  	category/categorytest	0.001s



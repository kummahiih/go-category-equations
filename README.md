# go-category-equations

An attemp to rewrite the https://github.com/kummahiih/python-category-equations -library in Golang. 
The library should be fluent to use and have the same tests or more. 
The full API documentation can be found from https://godoc.org/github.com/kummahiih/go-category-equations .

## Usage Example

Implement 'category.Connectable' (here 'NewConnectable') and 'category.Operator' (here 'NewConnectionPrinter'). 
This has been done for tests in 'category/categorytest' 
and the usage can be seen for example from a test 'TestExampleForDocumentation':

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



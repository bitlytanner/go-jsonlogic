# Fork notes

This is a fork of [HuanTeng/go-jsonlogic](https://github.com/HuanTeng/go-jsonlogic). This is designed to be a temporary fork until the behavior I want can be proposed/merged into the original repo. 

This repo has two key differences from the original:
1. When an object with multiple keys is found, we return it to the original caller. It is basically saying "we found an object with 2 or more fields. Let's return it."
2. The `===` and `!==` have been changed to `==` and `!=` respectively

# go-jsonlogic

[![Travis CI](https://travis-ci.org/HuanTeng/go-jsonlogic.svg?branch=master)](https://travis-ci.org/HuanTeng/go-jsonlogic)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitlytanner/go-jsonlogic)](https://goreportcard.com/report/github.com/bitlytanner/go-jsonlogic)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3e9df51b227c47b6b903a2a78ae62072)](https://www.codacy.com/app/the729/go-jsonlogic?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=HuanTeng/go-jsonlogic&amp;utm_campaign=Badge_Grade)

Golang implementation of JsonLogic (jsonlogic.com), which is an abstract syntax tree (AST) represented as a JSON object. 

Custom operators are supported.

## Example

```golang
var rule, data interface{}

json.Unmarshal([]byte(`
{"if": [
	{"<": [{"var":"temp"}, 0] }, "freezing",
	{"<": [{"var":"temp"}, 100] }, "liquid",
	"gas"
]}
`), &rule)

json.Unmarshal([]byte(`{"temp":55}`), &data)

got, err := jsonlogic.Apply(rule, data)
if err != nil {
	// handle error
}

fmt.Println(got)
// Output: liquid
```

## Custom operators

You can add your own operator to `jsonlogic` by implementing `jsonlogic.Operator` interface and registering it to a `jsonlogic` instance.

The following example shows how to implement a greatest common divisor (gcd) operator.

```golang

// gcd operator
type gcd struct{}

// Operate implements `jsonlogic.Operator`
func (gcd) Operate(applier jsonlogic.LogicApplier, data jsonlogic.DataType, params []jsonlogic.RuleType) (jsonlogic.DataType, error) {
	if len(params) != 2 {
		return nil, errors.New("only support 2 params")
	}

	var (
		p0, p1 interface{}
		err    error
	)
	// apply jsonlogic to each parameters recursively
	if p0, err = applier.Apply(params[0], data); err != nil {
		return nil, err
	}
	if p1, err = applier.Apply(params[1], data); err != nil {
		return nil, err
	}
	p0f, ok0 := p0.(float64)
	p1f, ok1 := p1.(float64)
	if !ok0 || !ok1 {
		return nil, errors.New("params should be numbers")
	}

	// recursive GCD function
	var gcdFunc func(a int, b int) int
	gcdFunc = func(a int, b int) int {
		if b == 0 {
			return a
		}
		return gcdFunc(b, a%b)
	}
	out := gcdFunc(int(p0f), int(p1f))

	// to output a number, always use float64
	return float64(out), nil
}

// Create an instance of `jsonlogic`, and register gcd operator
jl := jsonlogic.NewJSONLogic()
jl.AddOperation("gcd", gcd{})

// Use `gcd` as an operator to calculate: gcd(14+1, 25)
var rule interface{}
json.Unmarshal([]byte(`{"gcd": [{"+": [14, 1]}, 25]}`), &rule)

got, err := jl.Apply(rule, nil)
if err != nil {
	// handle error
}

fmt.Println(got)
// Output: 5
```

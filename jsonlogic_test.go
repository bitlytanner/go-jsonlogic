package jsonlogic_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	jsonlogic "github.com/bitlytanner/go-jsonlogic"
)

type TestCase struct {
	name   string
	rule   string
	data   string
	expect interface{}
}

var errorCases = []TestCase{
	{
		name:   "empty op",
		rule:   `{"": 1}`,
		data:   `{}`,
		expect: fmt.Errorf("operator  not found"),
	},
	{
		name:   "not found op",
		rule:   `{"not_found": 1}`,
		data:   `{}`,
		expect: fmt.Errorf("operator not_found not found"),
	},
	{
		name:   "not quick access",
		rule:   `{"$id": 1}`,
		data:   `{}`,
		expect: fmt.Errorf("quick access op not defined"),
	},
}

func runTestCases(cases []TestCase, t *testing.T) {
	var rule, data interface{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(c.rule), &rule); err != nil {
				t.Errorf("rule error: %s", err)
			}
			if err := json.Unmarshal([]byte(c.data), &data); err != nil {
				t.Errorf("data error: %s", err)
			}
			got, err := jsonlogic.Apply(rule, data)
			if err != nil {
				t.Errorf("apply error: %s", err)
			} else if !reflect.DeepEqual(got, c.expect) {
				t.Errorf("expect %+v got %+v", c.expect, got)
			}
		})
	}
}

func TestErrorCases(t *testing.T) {
	var rule, data interface{}

	jl := jsonlogic.NewJSONLogic()

	jl.AddOperation("_quick_access", nil)

	for _, c := range errorCases {
		if err := json.Unmarshal([]byte(c.rule), &rule); err != nil {
			t.Errorf("Case %s: rule error: %s", c.name, err)
		}
		if err := json.Unmarshal([]byte(c.data), &data); err != nil {
			t.Errorf("Case %s: data error: %s", c.name, err)
		}
		_, err := jl.Apply(rule, data)
		if !reflect.DeepEqual(err, c.expect) {
			t.Errorf("Case %s: expect error %+v got %+v", c.name, c.expect, err)
		}
	}
}


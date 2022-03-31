package goway

import (
	"net/http"
	"reflect"
	"testing"
)

// test slice dups removing.
func TestRemoveDuplicatesFromSlice(t *testing.T) {
	var slice = []string{"1", "2", "3", "4", "1", "2", "3", "4"}
	var expectedSlice = []string{"1", "2", "3", "4"}
	var result = removeDuplicateValues(slice)
	var isSame = reflect.DeepEqual(result, expectedSlice)
	if !isSame {
		t.Fatal("expected same slices")
	}
}

func TestIsVar(t *testing.T) {
	var trueCase = "{id}"
	var res, name = isRouteVar(trueCase)
	if !res || name != "id" {
		t.Fatal("wrong result")
	}

	var falseCase = "{id"
	res, name = isRouteVar(falseCase)
	if res {
		t.Fatal("wrong result")
	}
}

// test request set/get variables.
func TestGetSetVars(t *testing.T) {
	var req, err = http.NewRequest(http.MethodGet, "http://127.0.0.1", nil)
	if err != nil {
		t.Fatal(err)
	}
	addVarToContext(req, "hello", "world")
	var vars = Vars(req)
	var varVal, ok = vars["hello"]
	if !ok {
		t.Fatal("failed to get var")
	}
	if varVal != "world" {
		t.Fatalf("wrong var value")
	}
}

func TestIsMethodAllowed(t *testing.T) {
	type caser struct {
		num       int
		methods   []string
		reqMethod string
		expected  bool
	}
	var cases = []caser{
		// true.
		{
			num:       1,
			methods:   []string{http.MethodGet, http.MethodPut},
			reqMethod: http.MethodPut,
			expected:  true,
		},
		{
			num:       2,
			methods:   nil,
			reqMethod: http.MethodPut,
			expected:  true,
		},

		// false.
		{
			num:       3,
			methods:   []string{http.MethodGet},
			reqMethod: http.MethodDelete,
			expected:  false,
		},
		{
			num:       4,
			methods:   []string{},
			reqMethod: http.MethodPut,
			expected:  false,
		},
	}
	for _, cased := range cases {
		var result = isMethodAllowed(cased.methods, cased.reqMethod)
		if result != cased.expected {
			t.Fatalf("case num: %v | expected: %v | got: %v", cased.num, cased.expected, result)
		}
	}
}

func TestRemoveSlashStartEnd(t *testing.T) {
	type caser struct {
		num      int
		value    string
		expected string
	}
	var cases = []caser{
		{
			num:      1,
			value:    "",
			expected: "",
		},
		{
			num:      2,
			value:    "/hello/world/",
			expected: "hello/world",
		},
		{
			num:      3,
			value:    "hello/world",
			expected: "hello/world",
		},
		{
			num:      4,
			value:    "hello//////world////////",
			expected: "hello/world",
		},
	}
	for _, cased := range cases {
		var result = removeSlashStartEnd(cased.value)
		if result != cased.expected {
			t.Fatalf("case num: %v | expected: %v | got: %v", cased.num, cased.expected, result)
		}
	}
}

func TestPathToStandart(t *testing.T) {
	type caser struct {
		num      int
		value    string
		expected string
	}
	var cases = []caser{
		{
			num:      1,
			value:    "",
			expected: "",
		},
		{
			num:      2,
			value:    "///hello////world///",
			expected: "/hello/world",
		},
		{
			num:      3,
			value:    "hello/world",
			expected: "/hello/world",
		},
		{
			num:      4,
			value:    "hello/world///",
			expected: "/hello/world",
		},
	}
	for _, cased := range cases {
		var result = pathToStandart(cased.value)
		if result != cased.expected {
			t.Fatalf("case num: %v | expected: %v | got: %v", cased.num, cased.expected, result)
		}
	}
}

func TestProcessAllowedMethods(t *testing.T) {
	type caser struct {
		num        int
		methods    []string
		addMethods []string
		expected   []string
	}
	var cases = []caser{
		{
			num:        1,
			methods:    []string{http.MethodGet, http.MethodDelete, http.MethodGet},
			addMethods: []string{http.MethodGet, http.MethodPut},
			expected:   []string{http.MethodGet, http.MethodDelete, http.MethodPut},
		},
		{
			num:        2,
			methods:    nil,
			addMethods: []string{http.MethodDelete, http.MethodPut, http.MethodGet, http.MethodDelete},
			expected:   []string{http.MethodDelete, http.MethodPut, http.MethodGet},
		},
		{
			num:        3,
			methods:    []string{http.MethodGet},
			addMethods: nil,
			expected:   []string{http.MethodGet},
		},
	}
	for _, cased := range cases {
		var result = processAllowedMethods(cased.methods, cased.addMethods...)
		if !reflect.DeepEqual(cased.expected, result) {
			t.Fatalf("case num: %v | expected: %v | got: %v", cased.num, cased.expected, result)
		}
	}
}

package set

import (
	"fmt"
	"strings"
	"testing"
)

// NewStringSet() is called and the result is validated.
func TestNewStringSet(t *testing.T) {
	if ss := NewStringSet(); !ss.IsEmpty() {
		t.Fatalf("expected: true, got: false")
	}
}

// CreateStringSet() is called and the result is validated.
func TestCreateStringSet(t *testing.T) {
	ss := CreateStringSet("foo")
	if str := ss.String(); str != `[foo]` {
		t.Fatalf("expected: %s, got: %s", `["foo"]`, str)
	}
}

// StringSet.Add() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetAdd(t *testing.T) {
	testCases := []struct {
		name           string
		value          string
		expectedResult string
	}{
		// Test first addition.
		{"test1", "foo", `[foo]`},
		// Test duplicate addition.
		{"test2", "foo", `[foo]`},
		// Test new addition.
		{"test3", "bar", `[bar foo]`},
	}

	ss := NewStringSet()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ss.Add(testCase.value)
			if str := ss.String(); str != testCase.expectedResult {
				t.Fatalf("test %v expected: %s, got: %s", testCase.name, testCase.expectedResult, str)
			}
		})
	}
}

// StringSet.Remove() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetRemove(t *testing.T) {
	ss := CreateStringSet("foo", "bar")
	testCases := []struct {
		name           string
		value          string
		expectedResult string
	}{
		// Test removing non-existen item.
		{"test1", "baz", `[bar foo]`},
		// Test remove existing item.
		{"test2", "foo", `[bar]`},
		// Test remove existing item again.
		{"test2", "foo", `[bar]`},
		// Test remove to make set to empty.
		{"test3", "bar", `[]`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ss.Remove(testCase.value)
			if str := ss.String(); str != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, str)
			}
		})
	}
}

// StringSet.Contains() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetContains(t *testing.T) {
	ss := CreateStringSet("foo")
	testCases := []struct {
		name           string
		value          string
		expectedResult bool
	}{
		// Test to check non-existent item.
		{"test1", "bar", false},
		// Test to check existent item.
		{"test2", "foo", true},
		// Test to verify case sensitivity.
		{"test3", "Foo", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := ss.Contains(testCase.value); result != testCase.expectedResult {
				t.Fatalf("expected: %t, got: %t", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.FuncMatch() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetFuncMatch(t *testing.T) {
	ss := CreateStringSet("foo", "bar")
	testCases := []struct {
		name           string
		matchFn        func(string, string) bool
		value          string
		expectedResult string
	}{
		// Test to check match function doing case insensive compare.
		{"test1", func(setValue string, compareValue string) bool {
			return strings.EqualFold(setValue, compareValue)
		}, "Bar", `[bar]`},
		// Test to check match function doing prefix check.
		{"test2", func(setValue string, compareValue string) bool {
			return strings.HasPrefix(compareValue, setValue)
		}, "foobar", `[foo]`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := ss.FuncMatch(testCase.matchFn, testCase.value)
			if result := s.String(); result != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.ApplyFunc() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetApplyFunc(t *testing.T) {
	ss := CreateStringSet("foo", "bar")
	testCases := []struct {
		name           string
		applyFn        func(string) string
		expectedResult string
	}{
		// Test to apply function prepending a known string.
		{"test1", func(setValue string) string { return "mybucket/" + setValue }, `[mybucket/bar mybucket/foo]`},
		// Test to apply function modifying values.
		{"test2", func(setValue string) string { return setValue[1:] }, `[ar oo]`},
	}

	for _, testCase := range testCases {
		s := ss.ApplyFunc(testCase.applyFn)
		if result := s.String(); result != testCase.expectedResult {
			t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
		}
	}
}

// StringSet.Equals() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetEquals(t *testing.T) {
	testCases := []struct {
		name           string
		set1           StringSet
		set2           StringSet
		expectedResult bool
	}{
		// Test equal set
		{"test1", CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar"), true},
		// Test second set with more items
		{"test2", CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar", "baz"), false},
		// Test second set with less items
		{"test3", CreateStringSet("foo", "bar"), CreateStringSet("bar"), false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := testCase.set1.Equals(testCase.set2); result != testCase.expectedResult {
				t.Fatalf("expected: %t, got: %t", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.Intersection() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetIntersection(t *testing.T) {
	testCases := []struct {
		name           string
		set1           StringSet
		set2           StringSet
		expectedResult StringSet
	}{
		// Test intersecting all values.
		{"test1", CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar")},
		// Test intersecting all values in second set.
		{"test2", CreateStringSet("foo", "bar", "baz"), CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar")},
		// Test intersecting different values in second set.
		{"test3", CreateStringSet("foo", "baz"), CreateStringSet("baz", "bar"), CreateStringSet("baz")},
		// Test intersecting none.
		{"test4", CreateStringSet("foo", "baz"), CreateStringSet("poo", "bar"), NewStringSet()},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := testCase.set1.Intersection(testCase.set2); !result.Equals(testCase.expectedResult) {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.Difference() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetDifference(t *testing.T) {
	testCases := []struct {
		name           string
		set1           StringSet
		set2           StringSet
		expectedResult StringSet
	}{
		// Test differing none.
		{"test1", CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar"), NewStringSet()},
		// Test differing in first set.
		{"test2", CreateStringSet("foo", "bar", "baz"), CreateStringSet("foo", "bar"), CreateStringSet("baz")},
		// Test differing values in both set.
		{"test3", CreateStringSet("foo", "baz"), CreateStringSet("baz", "bar"), CreateStringSet("foo")},
		// Test differing all values.
		{"test4", CreateStringSet("foo", "baz"), CreateStringSet("poo", "bar"), CreateStringSet("foo", "baz")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := testCase.set1.Difference(testCase.set2); !result.Equals(testCase.expectedResult) {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.Union() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetUnion(t *testing.T) {
	testCases := []struct {
		name           string
		set1           StringSet
		set2           StringSet
		expectedResult StringSet
	}{
		// Test union same values.
		{"test1", CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar")},
		// Test union same values in second set.
		{"test2", CreateStringSet("foo", "bar", "baz"), CreateStringSet("foo", "bar"), CreateStringSet("foo", "bar", "baz")},
		// Test union different values in both set.
		{"test2", CreateStringSet("foo", "baz"), CreateStringSet("baz", "bar"), CreateStringSet("foo", "baz", "bar")},
		// Test union all different values.
		{"test2", CreateStringSet("foo", "baz"), CreateStringSet("poo", "bar"), CreateStringSet("foo", "baz", "poo", "bar")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result := testCase.set1.Union(testCase.set2); !result.Equals(testCase.expectedResult) {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.MarshalJSON() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetMarshalJSON(t *testing.T) {
	testCases := []struct {
		name           string
		set            StringSet
		expectedResult string
	}{
		// Test set with values.
		{"test1", CreateStringSet("foo", "bar"), `["bar","foo"]`},
		// Test empty set.
		{"test2", NewStringSet(), "[]"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if result, _ := testCase.set.MarshalJSON(); string(result) != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, string(result))
			}
		})
	}
}

// StringSet.UnmarshalJSON() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name           string
		data           []byte
		expectedResult string
	}{
		// Test to convert JSON array to set.
		{"test1", []byte(`["bar","foo"]`), `[bar foo]`},
		// Test to convert JSON string to set.
		{"test2", []byte(`"bar"`), `[bar]`},
		// Test to convert JSON empty array to set.
		{"test3", []byte(`[]`), `[]`},
		// Test to convert JSON empty string to set.
		{"test4", []byte(`""`), `[]`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var set StringSet
			set.UnmarshalJSON(testCase.data)
			if result := set.String(); result != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, result)
			}
		})
	}
}

// StringSet.String() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetString(t *testing.T) {
	testCases := []struct {
		name           string
		set            StringSet
		expectedResult string
	}{
		// Test empty set.
		{"test1", NewStringSet(), `[]`},
		// Test set with empty value.
		{"test2", CreateStringSet(""), `[]`},
		// Test set with value.
		{"test3", CreateStringSet("foo"), `[foo]`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if str := testCase.set.String(); str != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, str)
			}
		})
	}
}

// StringSet.ToSlice() is called with series of cases for valid and erroneous inputs and the result is validated.
func TestStringSetToSlice(t *testing.T) {
	testCases := []struct {
		name           string
		set            StringSet
		expectedResult string
	}{
		// Test empty set.
		{"test1", NewStringSet(), `[]`},
		// Test set with empty value.
		{"test2", CreateStringSet(""), `[]`},
		// Test set with value.
		{"test3", CreateStringSet("foo"), `[foo]`},
		// Test set with value.
		{"test4", CreateStringSet("foo", "bar"), `[bar foo]`},
	}

	for _, testCase := range testCases {
		t.Run("testCase.name", func(t *testing.T) {
			sslice := testCase.set.ToSlice()
			if str := fmt.Sprintf("%s", sslice); str != testCase.expectedResult {
				t.Fatalf("expected: %s, got: %s", testCase.expectedResult, str)
			}
		})
	}
}

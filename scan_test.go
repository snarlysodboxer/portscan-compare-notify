package main

import (
	"fmt"
	"strings"
	"testing"
)

// convertStringToIntSlice
type testStringToIntSlice struct {
	str  string
	ints []int
}

var stringToIntSliceTests = []testStringToIntSlice{
	{"432 123\r1 5234234\n213", []int{432, 123, 1, 5234234, 213}},
	{"789", []int{789}},
	{"789 ", []int{789}},
	{" 789 ", []int{789}},
	{" 789", []int{789}},
	{"5", []int{5}},
}

func TestConvertStringToIntSlice(t *testing.T) {
	for _, test := range stringToIntSliceTests {
		converted := convertStringToIntSlice(&test.str)
		converted_string := fmt.Sprintf("%d", converted)
		test_string := fmt.Sprintf("%d", test.ints)
		if converted_string != test_string {
			t.Error(
				"For", strings.Replace(strings.Replace(test.str, "\n", " ", -1), "\r", " ", -1),
				"expected", test.ints,
				"got", converted,
			)
		}
	}
}

// compare
type testComparePorts struct {
	expected         []int
	found            []int
	expected_unfound []int
	unexpected_found []int
}

//// expected, found, expected_unfound, unexpected_found
var comparePortsTests = []testComparePorts{
	{[]int{80, 443, 5432}, []int{80, 443}, []int{5432}, []int{}},
	{[]int{80, 443, 5432}, []int{80, 443, 9876}, []int{5432}, []int{9876}},

	{[]int{80, 443}, []int{80, 443, 5432}, []int{}, []int{5432}},
	{[]int{443, 80}, []int{80, 443, 5432}, []int{}, []int{5432}},
	{[]int{443, 80}, []int{5432, 443, 80}, []int{}, []int{5432}},

	{[]int{80, 443}, []int{443, 80}, []int{}, []int{}},
	{[]int{80, 443}, []int{80, 443}, []int{}, []int{}},
}

func TestComparePorts(t *testing.T) {
	for _, test := range comparePortsTests {
		expected_unfound, unexpected_found := comparePorts(test.expected, test.found)
		found_euf := fmt.Sprintf("%v", expected_unfound) != fmt.Sprintf("%v", test.expected_unfound)
		found_uef := fmt.Sprintf("%v", unexpected_found) != fmt.Sprintf("%v", test.unexpected_found)
		if found_euf || found_uef {
			t.Error(
				"For", test.expected, "and", test.found,
				"expected", test.expected_unfound, "and", test.unexpected_found,
				"got", expected_unfound, "and", unexpected_found,
			)
		}
	}
}

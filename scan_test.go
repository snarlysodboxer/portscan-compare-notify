package main

import (
	"fmt"
	"strings"
	"testing"
)

func removeNewLines(str string) string {
	return strings.Replace(strings.Replace(str, "\n", " ", -1), "\r", " ", -1)
}
func singleQuote(str string) string {
	return fmt.Sprintf("'%s'", str)
}

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
		convertedString := fmt.Sprintf("%d", converted)
		testString := fmt.Sprintf("%d", test.ints)
		if convertedString != testString {
			t.Error(
				"For", singleQuote(removeNewLines(test.str)),
				"expected", test.ints,
				"got", converted,
			)
		}
	}
}

// compare
type testComparePorts struct {
	expected        []int
	found           []int
	expectedUnfound []int
	unexpectedFound []int
}

//// expected, found, expectedUnfound, unexpectedFound
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
		expectedUnfound, unexpectedFound := comparePorts(test.expected, test.found)
		foundEuf := fmt.Sprintf("%v", expectedUnfound) != fmt.Sprintf("%v", test.expectedUnfound)
		foundUef := fmt.Sprintf("%v", unexpectedFound) != fmt.Sprintf("%v", test.unexpectedFound)
		if foundEuf || foundUef {
			t.Error(
				"For", test.expected, "and", test.found,
				"expected", test.expectedUnfound, "and", test.unexpectedFound,
				"got", expectedUnfound, "and", unexpectedFound,
			)
		}
	}
}

// grepNmap
type testGrepNmap struct {
	input          string
	output         string
	notShown       bool
	notShownNumber []string
}

var grepNmapTests = []testGrepNmap{
	{"asdf\nfdsaak sldddd\n\rsld sdlbv \n2 \n239bn2 \n\n\r\n\n02jk45 \n43\n30489 \n34", "2 239 02 43 30489 34", false, []string{}},
	{"asdf\nfdsaak sldddd\n\rsld sdlbv \n2\nNot shown: (closed) 234234\n239bn2 02jk45\n43 30489 34", "2 239 43", true, []string{"234234"}},
	{"asdf\nsdlbv \n2\nNot shown: 234234\n239bn2 02jk45\n43 30489 34", "2 239 43", false, []string{"234234"}},
	{"asdf\nfdsaak sldddd\n\rsld sdlbv 2Not shown: 234234\n239bn2 \n02j 98k45\n43 \n30489\n34", "239 02 43 30489 34", false, []string{}},
	{"asdf\nfdsaak sldddd\n\rsld sdlbv \n2\nNot shown: 234234 closed 32\n239bn2 \n\n02jk45\n43 \n30489 \n", "2 239 02 43 30489", true, []string{"234234", "32"}},
	{"asdf\nfdsaak sldddd\n\rsld \n34", "34", false, []string{}},
	{"\n23asdf\nfdsaak sldddd\n6\rsld", "23 6", false, []string{}},
}

func TestGrepNmap(t *testing.T) {
	for _, test := range grepNmapTests {
		output, notShown, notShownQuantity := grepNmap(test.input)
		nsns := fmt.Sprintf("%s", test.notShownNumber) != fmt.Sprintf("%s", notShownQuantity)
		if test.output != output || test.notShown != notShown || nsns {
			t.Error(
				"For", singleQuote(removeNewLines(test.input)),
				"expected", singleQuote(test.output), "and", test.notShown, "and", fmt.Sprintf("%v,", test.notShownNumber),
				"got", singleQuote(output), "and", notShown, "and", notShownQuantity,
			)
		}
	}
}

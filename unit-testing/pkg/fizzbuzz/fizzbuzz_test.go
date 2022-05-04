package fizzbuzz

import (
	"fmt"
	"testing"
)

func intToIntPtr(a int) *int {
	return &a
}

func TestFizzBuzz(t *testing.T) {
	categoryToTestCaseMap := map[string][]struct {
		total                  int
		exptectedTotalOverride *int
		fizzAt                 int
		buzzAt                 int
		expected               []string
	}{
		"SimpleCase": {
			{total: 1, fizzAt: 10, buzzAt: 10, expected: []string{"1"}},
			{total: 5, fizzAt: 10, buzzAt: 10, expected: []string{"1", "2", "3", "4", "5"}},
			{total: 1, fizzAt: 1, buzzAt: 1, expected: []string{"FizzBuzz"}},
			{total: 6, fizzAt: 2, buzzAt: 3, expected: []string{"1", "Fizz", "Buzz", "Fizz", "5", "FizzBuzz"}},
			{total: 5, fizzAt: 1, buzzAt: 5, expected: []string{"Fizz", "Fizz", "Fizz", "Fizz", "FizzBuzz"}},
		},

		"TotalEdgeCase": {
			{total: 0, fizzAt: 1, buzzAt: 1, expected: []string{}},
			{total: -1, exptectedTotalOverride: intToIntPtr(0), fizzAt: 1, buzzAt: 1, expected: []string{}},
		},

		"DivideByZeroEdgeCase": {
			{total: 1, fizzAt: 0, buzzAt: 2, expected: []string{"1"}},
			{total: 1, fizzAt: 2, buzzAt: 0, expected: []string{"1"}},
			{total: 1, fizzAt: 1, buzzAt: 0, expected: []string{"Fizz"}},
			{total: 1, fizzAt: 0, buzzAt: 1, expected: []string{"Buzz"}},
			{total: 1, fizzAt: 0, buzzAt: 0, expected: []string{"1"}},
		},

		"NegativeIntegersCase": {
			{total: 1, fizzAt: -1, buzzAt: 10, expected: []string{"Fizz"}},
			{total: 1, fizzAt: 10, buzzAt: -1, expected: []string{"Buzz"}},
			{total: 1, fizzAt: -1, buzzAt: -1, expected: []string{"FizzBuzz"}},
			{total: 1, fizzAt: -1, buzzAt: 1, expected: []string{"FizzBuzz"}},
			{total: 6, fizzAt: -2, buzzAt: -3, expected: []string{"1", "Fizz", "Buzz", "Fizz", "5", "FizzBuzz"}},
		},
	}

	for category, testCases := range categoryToTestCaseMap {
		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%s(total=%d,fizzAt=%v,buzzAt=%v)", category, testCase.total, testCase.fizzAt, testCase.buzzAt), func(t *testing.T) {
				//validate test case setup
				expectedTotal := len(testCase.expected)

				total := testCase.total

				if testCase.exptectedTotalOverride != nil {
					total = *testCase.exptectedTotalOverride
				}

				if total != expectedTotal {
					t.Fatalf("invalid test case: total != len(expected) ( total = %d, len(expected) = %d )", total, expectedTotal)
				}

				actual := FizzBuzz(testCase.total, testCase.fizzAt, testCase.buzzAt)

				actualTotal := len(actual)

				if total != actualTotal {
					t.Errorf("invalid total: expected %d but got %d", total, actualTotal)
				}

				for i := 0; i < actualTotal; i++ {
					if testCase.expected[i] != actual[i] {
						t.Errorf("invalid result[%d]: expected '%s' but got '%s'", i, testCase.expected[i], actual[i])
					}
				}
			})
		}
	}
}

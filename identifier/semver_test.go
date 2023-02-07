/*
 * Copyright 2023 Asim Ihsan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package identifier

import (
	"strings"
	"testing"
)

func TestDoesSemverMatch(t *testing.T) {
	tests := []struct {
		version     string
		requirement string
		expected    bool
	}{
		// exact / caret requirements
		{"1.2.3", "1.2.3", true},
		{"1.2.3", "^1.2.3", true},
		{"1.2.3", "1.2", false},

		// tilde requirements
		{"1.2.3", "~1.2.3", true},
		{"1.2.4", "~1.2.3", true},
		{"1.3", "~1.2.3", false},
		{"1.2.3", "~1.2", true},
		{"1.2.4", "~1.2", true},
		{"1.3", "~1.2", false},
		{"2.0", "~1.2", false},
		{"1.2.3", "~1", true},
		{"1.2.4", "~1", true},
		{"1.3", "~1", true},
		{"2.0", "~1", false},
		{"1.2.3", "~2", false},

		// comparison requirements, just a single > or >=
		{"1.2.3", ">=1.2", true},
		{"1.2.3", "> 1.2", true},
		{"1.1", ">= 1.2", false},
	}

	for _, test := range tests {
		actual := Satisfies(test.version, test.requirement)
		if actual != test.expected {
			t.Errorf("Satisfies(%s, %s) = %t, want %t", test.version, test.requirement, actual, test.expected)
		}
	}
}

func TestRegressionFuzzDoesSemverMatch_01(t *testing.T) {
	actual := Satisfies("1", "~1.0")
	if actual != true {
		t.Errorf("Satisfies(1, ~1.0) = %t, want %t", actual, true)
	}
}

func TestRegressionFuzzDoesSemverMatch_02(t *testing.T) {
	actual := Satisfies("1", "~1.0.0")
	if actual != true {
		t.Errorf("Satisfies(1, ~1.0.0) = %t, want %t", actual, true)
	}
}

func FuzzDoesSemverMatch(f *testing.F) {
	// seed the corpus. each testcase is space delimited, first element is version, second is requirement.
	for _, testcase := range []string{
		"1.2.3 1.2.3",
		"1.2.3 ^1.2.3",
		"1.2.3 1.2",
		"1.2.3 ~1.2.3",
		"1.2.4 ~1.2.3",
		"1.3 ~1.2.3",
		"1.2.3 ~1.2",
		"1.2.4 ~1.2",
		"1.3 ~1.2",
		"2.0 ~1.2",
		"1.2.3 ~1",
		"1.2.4 ~1",
		"1.3 ~1",
		"2.0 ~1",
		"1.2.3 ~2",
		"1.2.3 >=1.2",
		"1.2.3 >1.2",
		"1.2.3 <1.2",
		"1.2.3 <=1.2",
		"1.2.3 ==1.2",
	} {
		f.Add([]byte(testcase))
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		// split the input into version and requirement
		split := strings.Split(string(input), " ")
		if len(split) != 2 {
			t.Skip("invalid input")
		}
		versionString := split[0]
		requirement := split[1]

		// run the test
		Satisfies(versionString, requirement)
	})
}

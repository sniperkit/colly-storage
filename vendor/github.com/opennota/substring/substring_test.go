// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package substring

import (
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"
)

const longString = "All rolling hash functions are linear in the number of characters, but their complexity with respect to the length of the window (k) varies."

var testCases = []struct {
	s      string
	substr string
	want   bool
}{
	{"", "", true},
	{"abc", "", true},
	{"", "abc", false},
	{"a", "a", true},
	{"a", "A", true},
	{"A", "a", true},
	{"abc", "a", true},
	{"abc", "A", true},
	{"ABC", "b", true},
	{"abc", "B", true},
	{"ABC", "c", true},
	{"abc", "C", true},
	{"abc", "d", false},
	{"abc", "D", false},
	{"abc", "aa", false},
	{"abc", "Aa", false},
	{"abc", "aaA", false},
	{"ABC", "abc", true},
	{"abc", "aBc", true},
	{"abc", "aBcD", false},
	{"aBC", "ab", true},
	{"abc", "AB", true},
	{"abc", "Ab", true},
	{"abc", "aB", true},
	{"abc", "abd", false},
	{"abc", "ac", false},
	{"abc", "BC", true},
	{"abc", "BC!", false},
}

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randomizeCase(s string) string {
	rr := []rune(s)
	for i, r := range rr {
		switch rnd.Intn(2) {
		case 0:
			rr[i] = unicode.ToLower(r)
		case 1:
			rr[i] = unicode.ToUpper(r)
		}
	}
	return string(rr)
}

func TestMatch(t *testing.T) {
	for _, tc := range testCases {
		m := NewMatcher(tc.substr)
		got := m.Match(tc.s)
		if got != tc.want {
			t.Errorf("NewMatcher(%q).Match(%q) = %v, want %v", tc.substr, tc.s, got, tc.want)
		}
	}

	for n := 1; n <= len(longString); n++ {
		for i := 0; i < len(longString)-n+1; i++ {
			substr := randomizeCase(longString[i : i+n])
			m := NewMatcher(substr)
			if !m.Match(longString) {
				t.Fatalf("NewMatcher(%q) must match %q", substr, longString)
			}
		}
	}
}

var boolResult bool

func BenchmarkMatch(b *testing.B) {
	m := NewMatcher("wInDoW")
	for i := 0; i < b.N; i++ {
		boolResult = m.Match(longString)
	}
}

func BenchmarkStringsContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		boolResult = strings.Contains(longString, "window")
	}
}

func BenchmarkStringsContainsWithToLower(b *testing.B) {
	substr := strings.ToLower("wInDoW")
	for i := 0; i < b.N; i++ {
		boolResult = strings.Contains(strings.ToLower(longString), substr)
	}
}

func BenchmarkRegexpMatchString(b *testing.B) {
	r := regexp.MustCompile("(?i)" + regexp.QuoteMeta("wInDoW"))
	for i := 0; i < b.N; i++ {
		boolResult = r.MatchString(longString)
	}
}

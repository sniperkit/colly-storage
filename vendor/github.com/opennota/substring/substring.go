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

// Package substring implements case-insensitive substring matching
// using Rabin-Karp string searching algorithm.
package substring

import (
	"unicode"
	"unicode/utf8"
)

// Matcher is a case-insensitive substring matcher.
type Matcher struct {
	needle []rune
	n      int
	hash   uint32
	pow    uint32
	rbuf   []rune
}

// prime is the prime base used in Rabin-Karp algorithm.
const prime = 16777619

func pow(n int) uint32 {
	pw := uint32(1)
	sq := uint32(prime)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			pw *= sq
		}
		sq *= sq
	}
	return pw
}

// NewMatcher returns a new substring matcher which case-insensitively matches substring s.
func NewMatcher(s string) *Matcher {
	rr := []rune(s)
	n := len(rr)

	var h uint32
	for i, r := range rr {
		r = unicode.ToLower(r)
		rr[i] = r
		h = h*prime + uint32(r)
	}

	return &Matcher{
		needle: rr,
		n:      n,
		hash:   h,
		pow:    pow(n),
		rbuf:   make([]rune, n),
	}
}

func equal(rr, rbuf []rune, i int) bool {
	for _, r := range rr {
		if r != rbuf[i] {
			return false
		}
		i++
		if i == len(rbuf) {
			i = 0
		}
	}
	return true
}

// Match reports whether s contains the matcher substring (ignoring case).
func (m *Matcher) Match(s string) bool {
	i := 0
	var h uint32
	for j := 0; j < m.n; j++ {
		r, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			return false
		}
		i += size
		r = unicode.ToLower(r)

		h = h*prime + uint32(r)

		m.rbuf[j] = r
	}

	if h == m.hash && equal(m.needle, m.rbuf, 0) {
		return true
	}

	j := 0
	for _, r := range s[i:] {
		r = unicode.ToLower(r)

		h = h*prime + uint32(r)
		h -= m.pow * uint32(m.rbuf[j])

		m.rbuf[j] = r
		j++
		if j == m.n {
			j = 0
		}

		if h == m.hash && equal(m.needle, m.rbuf, j) {
			return true
		}
	}
	return false
}

// Contains reports whether substr is within s (ignoring case).
func Contains(s, substr string) bool {
	return NewMatcher(substr).Match(s)
}

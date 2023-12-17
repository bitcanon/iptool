/*
Copyright © 2024 Mikael Schultz <bitcanon@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"
	"io"
	"sort"
)

// getLongestKeyLength returns the length of the longest key in a map.
// This is used to align the values when printing the variables.
func getLongestKeyLength(m map[string]string) int {
	maxKeyLength := 0
	for key := range m {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}
	return maxKeyLength
}

func PrintMap(out io.Writer, m map[string]string) {
	// Extract and sort the keys
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Find the length of the longest key
	maxKeyLength := getLongestKeyLength(m)

	// Print the variables and pad the keys with spaces to align the values
	for _, key := range keys {
		value := m[key]
		// Left align (-) with padding spaces (*)
		fmt.Printf("%-*s : %s\n", maxKeyLength, key, value)
	}
}

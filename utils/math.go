/*
Copyright © 2024 Mikael Schultz <mikael@conf-t.se>

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

import "math"

// ClosestLargerPowerOfTwo returns the smallest power of 2 that is greater than or equal to n.
func ClosestLargerPowerOfTwo(n int) int {
	// The smallest power of 2 greater than or equal to any number less than 1 is 1 (2^0)
	if n < 1 {
		return 1
	}

	// Convert n to float64 for the logarithm function
	log2 := math.Log2(float64(n))

	// Ceil to find the next whole power of 2 if n is not already a power of 2
	ceilLog2 := math.Ceil(log2)

	// Convert back to int and raise 2 to this power
	return int(math.Pow(2, ceilLog2))
}

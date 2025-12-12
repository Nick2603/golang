package main

import (
	"fmt"
	"strconv"
)

func FibonacciIterative(n int) int {
	if n < 0 {
		return n
	}

	if n == 0 {
		return 0
	}

	if n == 1 {
		return 1
	}

	prev2 := 0
	prev1 := 1
	current := 0

	for i := 2; i <= n; i++ {
		current = prev1 + prev2
		prev2 = prev1
		prev1 = current
	}

	return current
}

func FibonacciRecursive(n int) int {
	if n < 0 {
		return n
	}

	if n == 0 {
		return 0
	}

	if n == 1 {
		return 1
	}

	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}

	if n == 2 {
		return true
	}

	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func IsBinaryPalindrome(n int) bool {
	if n < 0 {
		return false
	}

	binary := strconv.FormatInt(int64(n), 2)

	left := 0
	right := len(binary) - 1

	for left < right {
		if binary[left] != binary[right] {
			return false
		}

		left++
		right--
	}

	return true
}

func ValidParentheses(s string) bool {
	stack := []rune{}

	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]

			if top != pairs[char] {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func Increment(num string) int {
	if num == "" {
		return 0
	}

	for _, char := range num {
		if char != '0' && char != '1' {
			return 0
		}
	}

	decimal, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		return 0
	}

	return int(decimal) + 1
}

func main() {
	// Невеликі демонстраційні виклики (для наочного запуску `go run .`)
	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10))
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10))

	fmt.Println("IsPrime(2):", IsPrime(2))
	fmt.Println("IsPrime(15):", IsPrime(15))
	fmt.Println("IsPrime(29):", IsPrime(29))

	fmt.Println("IsBinaryPalindrome(7):", IsBinaryPalindrome(7))
	fmt.Println("IsBinaryPalindrome(6):", IsBinaryPalindrome(6))

	fmt.Println(`ValidParentheses("[]{}()"):`, ValidParentheses("[]{}()"))
	fmt.Println(`ValidParentheses("[{]}"):`, ValidParentheses("[{]}"))

	fmt.Println(`Increment("101") ->`, Increment("101"))
}

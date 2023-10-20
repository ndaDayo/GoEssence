package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Divisor must not be zero")
	}

	return a / b, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Number of arguments must be 2")
		os.Exit(1)
	}

	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "First argument must be float value: %v", err)
		os.Exit(1)
	}
	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Second argument must be float value: %v", err)
		os.Exit(1)
	}

	result, err := divide(a, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid argument: %v", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

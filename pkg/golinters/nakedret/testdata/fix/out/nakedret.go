//golangcitest:args -Enakedret
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func NakedretIssue() (a int, b string) {
	if a > 0 {
		return a, b
	}

	fmt.Println("nakedret")

	if b == "" {
		return 0, "0"
	}

	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...

	// len of this function is 33
	return a, b
}

func NoNakedretIssue() (a int, b string) {
	if a > 0 {
		return
	}

	if b == "" {
		return 0, "0"
	}

	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...

	// len of this function is 30
	return
}

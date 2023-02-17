package version

import "fmt"

func ExampleGetVersion() {
	SetVersion("1.2.3")
	fmt.Println(GetVersion())
	// Output: 1.2.3
}

func ExampleSetVersion() {
	SetVersion("1.23.1")
	fmt.Println(GetVersion())
	// Output: 1.23.1
}

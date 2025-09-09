package prt

import (
	"fmt"
	"os"
)

func ExampleNewStringPrinter() {
	p := NewStringPrinter()

	p.Println("Tho", "mas ", 1976)
	p.Println("Mause ", 1971)

	fmt.Print(p.String())
	// Output:
	// Thomas 1976
	// Mause 1971
}

func ExampleBuildString() {
	s := BuildString(func(p *Printer) {
		p.Println(10, " + ", 10, " = ", 10+10)
	})
	fmt.Print(s)
	// Output:
	// 10 + 10 = 20
}

func ExampleBuildFile() {
	err := BuildFile("buildfile.txt", func(p *Printer) {
		p.Printfln("foobar %d", 10)
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := os.ReadFile("buildfile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(string(data))
	
	err = BuildFile("/nopermission.txt", func(p *Printer) {
		p.Println("i'm not allowed!")
	})
	fmt.Println(err)
	
	// Output:
	// foobar 10
	// open /nopermission.txt: permission denied
}

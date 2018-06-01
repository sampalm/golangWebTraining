package main

import (
	"fmt"
	"os"
)

type dollar float64
type euro float64

func (d dollar) String() string { return fmt.Sprintf("$%.2f", d) }
func (e euro) String() string   { return fmt.Sprintf("€%.2f", e) }

func main() {
	d := dollar(50)
	// Format Dollar to String
	fmt.Fprintf(os.Stdout, "Dollar: %.2f\nFDollar: %s", d, d)
	e := euro(65)
	// Format Euro to String
	fmt.Fprintf(os.Stdout, "\nDollar: %.2f\nFDollar: %s", e, e)
}

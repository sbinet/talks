// hello is a simple package
package hello

import "fmt"

// Hello greets someone.
func Hello(name string) string {
	return fmt.Sprintf("hello %q from Go", name)
}

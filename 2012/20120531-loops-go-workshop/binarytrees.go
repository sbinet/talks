package main

import "code.google.com/p/go-tour/tree"
import (
 "fmt"
 _ "time"
)
 
//type Tree struct {
//Left  *Tree
//Value int
//Right *Tree
//}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	var walk func(t *tree.Tree)
	walk = func(t *tree.Tree) {
		if t.Left != nil {
			walk(t.Left)
		}
     ch <- t.Value
		if t.Right != nil {
			walk(t.Right)
		}
	}
	walk(t)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
ch2 := make(chan int)
go Walk(t1, ch1)
go Walk(t2, ch2)
    for {
	v1,ok1 := <-ch1
v2,ok2 := <-ch2
       if !ok1 && !ok2 {
          break // no more nodes in trees
}
if ok1 != ok2 {
	   // trees with different sizes
	   return false
}
if v2 != v1 {
          return false
}
}
return true
}

func main() {
	ch := make(chan int)
go Walk(tree.New(1), ch)
for v := range ch {
	fmt.Print(v, " ")
}
fmt.Println()
o11 := Same(tree.New(1), tree.New(1))
o12 := Same(tree.New(1), tree.New(2))
o21 := Same(tree.New(2), tree.New(1))
tt := tree.New(1)
tt.Right = tree.New(1)
fmt.Println("o11:",o11)
fmt.Println("o12:",o12)
fmt.Println("o21:",o21)
fmt.Println(":",Same(tt, tt))
fmt.Println(":",Same(tt, tree.New(1)))
}

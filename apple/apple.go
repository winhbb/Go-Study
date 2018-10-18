package apple

import "fmt"

type Green struct {
	color int
}

func Beluck() {
	g := Green{1}
	fmt.Printf("color is %d", g.color)
}

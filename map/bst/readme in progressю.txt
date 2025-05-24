
Пример

package main




import (
	"fmt"
	"github.com/Qwental/go-containers/map/bst"
)

func compareInt(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
func main() {
	m := bst.NewBSTMap[int, int](compareInt)
	m.Put(1, 2)
	fmt.Println(m.Get(1))
}

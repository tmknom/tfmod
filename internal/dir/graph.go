package dir

import "fmt"

type Graph[SRC, DST Path] struct {
	items map[string][]*DST
}

func NewGraph[SRC, DST Path]() *Graph[SRC, DST] {
	return &Graph[SRC, DST]{
		items: make(map[string][]*DST, 64),
	}
}

func (d *Graph[SRC, DST]) Add(src *SRC, dst *DST) {
	key := (*src).Rel()
	d.items[key] = append(d.items[key], dst)
}

func (d *Graph[SRC, DST]) Include(src string) bool {
	_, ok := d.items[src]
	return ok
}

func (d *Graph[SRC, DST]) ListDst(src string) []*DST {
	result, _ := d.items[src]
	return result
}

func (d *Graph[SRC, DST]) String() string {
	return fmt.Sprintf("%v", d.items)
}

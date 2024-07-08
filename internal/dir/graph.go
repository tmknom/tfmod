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

func (d *Graph[SRC, DST]) Include(src Path) bool {
	_, ok := d.items[src.Rel()]
	return ok
}

func (d *Graph[SRC, DST]) ListDst(src *SRC) []*DST {
	result, _ := d.items[(*src).Rel()]
	return result
}

func (d *Graph[SRC, DST]) String() string {
	return fmt.Sprintf("%v", d.items)
}

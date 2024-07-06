package collection

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type TreeSet struct {
	items map[string]bool
}

func NewTreeSet() *TreeSet {
	return &TreeSet{
		items: make(map[string]bool, 64),
	}
}

func (s *TreeSet) Add(item string) {
	s.items[item] = true
}

func (s *TreeSet) Slice() []string {
	result := make([]string, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}

func (s *TreeSet) ToJson() string {
	bytes, err := json.Marshal(s.Slice())
	if err != nil {
		log.Fatalf("not json marshaled: %#v", s)
	}
	return string(bytes)
}

func (s *TreeSet) String() string {
	return fmt.Sprintf("%v", s.Slice())
}

func (s *TreeSet) GoString() string {
	return fmt.Sprintf("%#v", s.Slice())
}

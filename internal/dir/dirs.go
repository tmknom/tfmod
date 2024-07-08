package dir

type Dirs []*Dir

func (d Dirs) Slice() []string {
	result := make([]string, 0, len(d))
	for _, item := range d {
		result = append(result, item.Rel())
	}
	return result
}

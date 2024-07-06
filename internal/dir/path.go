package dir

type Path interface {
	Abs() string
	Rel() string
}

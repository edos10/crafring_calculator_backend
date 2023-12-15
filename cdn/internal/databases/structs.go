package databases

type Database interface {
	// to fix to structs
	GetPath(id int) (string, error)
}

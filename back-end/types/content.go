package types

type FolderObject struct {
	ID         string
	Name       string
	Permission int
}

type FileObject struct {
	ID         string
	Name       string
	Permission int
	Sensitive  string
	Hash       string
}

type Content struct {
	Permission int
	Folders    []FolderObject
	Files      []FileObject
}

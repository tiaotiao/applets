package fileserver


// LocalFiles module manages local files only. 
// It doesn't deal with golbal lock and replication
// TODO cache file handlers for future use.
type LocalFiles struct {
	// TODO
}

func NewLocalFiles(root string) *LocalFiles {
	// TODO
	return nil
}

func (l *LocalFiles) Read(file string, offset int, n int) ([]byte, error) {
	// TODO
	return nil, nil
}

func (l *LocalFiles) Append(file string, data []byte) error {
	// TODO
	return nil, nil
}

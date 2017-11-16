package sdk

const PermRead string = "read"
const PermWrite string = "write"

type LockClient struct {
	// TODO
}

func NewLockClient() *LockClient {
	l := LockClient{}
	// TODO
	return &l
}

func (l *LockClient) Connect(addr string) error {
	// TODO
	return nil
}

func (l *LockClient) Require(path string, perm string, timeout int) error {
	// TODO
	return nil
}

func (l *LockClient) Release(path string) error {
	// TODO
	return nil
}

func (l *LockClient) Close() error {
	// TODO
	return nil
}

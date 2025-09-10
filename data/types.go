package data

type User struct {
	ID             string
	Username       string
	HashedPassword string
	HomeDir        string
	Roles          []string
}

type LimitType string

const (
	LimitStorage LimitType = "storage"
	LimitFiles   LimitType = "files"
	LimitTraffic LimitType = "traffic"
)

type Quota struct {
	ID     string
	UserID string
	Type   LimitType
	Max    int64
	Used   int64
}

package models

type ACL struct {
	UserID     int
	Resource   string
	Permission string
}

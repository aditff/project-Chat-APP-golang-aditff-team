// model/user.go
package model

import "time"

type User struct {
	ID       string
	Name     string
	Email    string
	Online   bool
	LastSeen *time.Time
}

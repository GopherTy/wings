package db

import "time"

// Adminstrator adminstrator table
type Adminstrator struct {
	ID        uint64    `xorm:"pk autoincr 'id'"`
	User      string    `xorm:"varchar(16) 'user'"`
	Password  string    `xorm:"varchar(18) 'passwd'"`
	IPAddress string    `xorm:" 'ip_address' "`
	LastLogin time.Time `xorm:"last_login"`
}

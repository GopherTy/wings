package db

import "time"

// Administrator table columns
const (
	AdministratorColID        = "id"
	AdministratorColUser      = "user"
	AdministratorColPassword  = "passwd"
	AdministratorColIPAddress = "ip_address"
	AdministratorColLastLogin = "last_login"
)

// Administrator adminstrator table
type Administrator struct {
	ID        uint64    `xorm:"pk autoincr 'id'"`
	User      string    `xorm:"varchar(16) unique 'user' "`
	Password  string    `xorm:"varchar(18) 'passwd'"`
	IPAddress string    `xorm:"'ip_address'"`
	LastLogin time.Time `xorm:"'last_login'"`
}

// TableName .
func (Administrator) TableName() string {
	return "administrator"
}

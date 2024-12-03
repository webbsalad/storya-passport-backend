package pg

import "time"

type User struct {
	ID           string `db:"id"`
	EmailId      string `db:"email_id"`
	Name         string `db:"name"`
	PassworgHash string `db:"password_hash"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

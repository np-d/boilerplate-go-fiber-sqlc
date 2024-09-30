// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID          int32            `db:"id" json:"id"`
	DisplayName string           `db:"display_name" json:"display_name"`
	Username    string           `db:"username" json:"username"`
	Email       string           `db:"email" json:"email"`
	Password    string           `db:"password" json:"password"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt   pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
}

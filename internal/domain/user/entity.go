package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID  `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

func NewUser(email, password, firstName, lastName string) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:        email,
		PasswordHash: string(bytes),
		FirstName:    firstName,
		LastName:     lastName,
	}, nil
}

func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return err == nil, err
}

func (u *User) Delete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

func (u *User) Restore() {
	u.DeletedAt = nil
	u.UpdatedAt = time.Now()
}

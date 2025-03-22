package user

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UID       string    `json:"uid" gorm:"column:uid;type:varchar(50);not null;primaryKey"`
	Name      string    `json:"name" gorm:"column:user_name;type:varchar(50);not null;"`
	Userid    string    `json:"id" gorm:"column:user_id;type:varchar(50);not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdateAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	IsDeleted bool      `json:"-" gorm:"column:is_deleted"`
	Avatar    string    `json:"avatar" gorm:"column:avatar"`
	Email     string    `json:"email" gorm:"column:user_email;"`
	Password  string    `json:"-" gorm:"column:password;"`
}

// GetUid implementation jwtUser interface
func (u *User) GetUid() string {
	return u.UID
}

// BeforeCreate generate uuid, crypt password, check is exist
func (u *User) BeforeCreate(tx *gorm.DB) error {
	err := tx.Model(&User{}).Where("user_email = ? and is_deleted = false", u.Email).First(&User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		return errors.New("email already taken")
	}
	u.CreatedAt = time.Now()
	u.UpdateAt = time.Now()
	u.Password = cryptor.Sha256(u.Password)
	u.UID = uuid.New().String()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdateAt = time.Now()
	return nil
}

func (*User) TableName() string {
	return "users" // customer table name
}

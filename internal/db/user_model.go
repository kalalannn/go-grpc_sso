package db

import (
	"grpc-sso/internal/crypt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"` // encrypted
}

func (user *User) EQPassword(plainPassword string) bool {
	return crypt.CMPHashedPlainPassword(user.Password, plainPassword)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(connectionString string) (*UserService, error) {
	db, err := InitGorm(connectionString)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})

	return &UserService{
		db: db,
	}, nil
}

func (us *UserService) CreateUser(email, username, plainPassword string) (*User, error) {
	hashedPassword, err := crypt.HashedPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := us.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) findOneBy(where string, by ...interface{}) (*User, error) {
	var user User
	if err := us.db.Where(where, by...).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) FindOneByUsername(username string) (*User, error) {
	return us.findOneBy("username = ?", username)
}

func (us *UserService) FindOneByEmail(email string) (*User, error) {
	return us.findOneBy("email = ?", email)
}

func (us *UserService) AuthUserWithUsername(username, plainPassword string) (*User, error) {
	user, err := us.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.EQPassword(plainPassword) {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (us *UserService) AuthUserWithEmail(email, plainPassword string) (*User, error) {
	user, err := us.FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.EQPassword(plainPassword) {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

package store

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserNotValid     = errors.New("Incorrect parametrs!")
	UserAlreadyExist = errors.New("user already exist!")
	AuthFiled        = errors.New("auth failed")
)

// User ...
type User struct {
	gorm.Model
	// ID        uint `json:"id";gorm:"primary_key"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	FirstName   string `json:"firstname";gorm:"not null"`
	LastName    string `json:"lastname";gorm:"not null`
	Email       string `json:"email";gorm:"unique,not null"`
	Phone       string `json:"phone"`
	Password    string `json:"password,omitempty"`
	NewPassword string `json:"newpassword"`
	Admin       string `json:"admin";gorm:"default:false"`
}

/** CRUD Methods **/
func AddUser(u *User) error {

	if u.Valid() != true {
		return UserNotValid
	}

	if IsUserExist(u.Email) {
		return UserAlreadyExist
	}

	if password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err == nil {
		u.Password = string(password)
	}

	return x.Create(u).Error
}

func FindUser() []*User {
	var users []*User
	x.Raw("SELECT id, first_name, last_name, phone, email, admin, created_at, updated_at, deleted_at FROM users").Scan(&users)

	return users
}

func FindUserLimit(page string) []*User {
	var users []*User

	offset, _ := strconv.Atoi(page)
	offset = offset - 1

	x.Raw("SELECT id, first_name, last_name, phone, email, admin, created_at, updated_at, deleted_at FROM users").Offset(offset * 2).Limit(2).Scan(&users)

	return users
}

func UpdateUser(u *User) error {
	return x.Save(u).Error
}

func DeleteUser(u *User) error {
	return x.Unscoped().Delete(&User{}, "email LIKE ?", u.Email).Error
}

func DeleteUserByID(id uint) error {
	if id > 0 {
		return x.Unscoped().Delete(&User{}, "id = ?", id).Error
	}
	return errors.New("id is not valid")
}

/** Helper functions **/
// func FindByEmail2(email string) (*User, error) {
// u := new(User)
// err := x.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&u).Error
// if err == nil {
// 	return UserAlreadyExist
// }
// if err == gorm.ErrRecordNotFound {
// 	return nil
// }
// return err
// }

func FindByEmail(email string) (*User, error) {
	u := new(User)
	if err := x.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func FindByID(id string) *User {
	u := new(User)
	x.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&u)
	return u
}

// TODO Изучить как работает форматирование дат
func (u *User) FormatDate(field string) string {
	if field == "CreatedAt" {
		return u.CreatedAt.Format("02.01.2006 15:04")
	} else {
		return u.UpdatedAt.Format("02.01.2006 15:04")
	}
}

// Valid ...
func (u *User) Valid() bool {
	// TODO Сделать номальную валидацию
	return strings.Contains(u.Email, "@") && len(u.FirstName) > 0 && len(u.LastName) > 0 && len(u.Phone) > 0
}

func IsUserExist(email string) bool {
	return x.Raw(
		"SELECT * FROM users WHERE @mail",
		sql.Named(
			"mail",
			strings.ToLower(email),
		),
	).Scan(&User{}).RowsAffected > 0
}

// ComparePass ..
func (u *User) ComparePass(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	return err == nil
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPhoneNumber() string {
	return u.Phone
}

/** HOOKS Database **/
// func (u *User) BeforeSave(scope *gorm.Scope) (err error) {

// 	// Generate Password!
// 	temporaryPass, _ := password.Generate(10, 4, 0, false, false)

// 	scope.SetColumn("EncryptedPassword", temporaryPass) // TODO Deleted after test

// 	// Send Message with password for user
// 	if err := utils.Send(u.Email, temporaryPass); err != nil {
// 		return err
// 	}

// 	// TODO Переделать, на hash и solt пароля
// 	if password, err := bcrypt.GenerateFromPassword([]byte(temporaryPass), 0); err == nil {
// 		scope.SetColumn("password", password)
// 	}

// 	return nil
// }

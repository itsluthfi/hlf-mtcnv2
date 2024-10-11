package models

import (
	"errors"

	"html"
	"strings"

	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/utils/token"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:225;not null;unique" json:"username"`
	Password  string `gorm:"size:225;not null" json:"password"`
	AccountID string `gorm:"size:255" json:"accountId"`
	Email     string `gorm:"size:225" json:"email"`
	Name      string `gorm:"size:225" json:"name"`
	Phone     string `gorm:"size:255" json:"phone"`
}

type Transactions struct {
	gorm.Model
	Sender   string `gorm:"size:255" json:"sender"`
	Receiver string `gorm:"size:255" json:"receiver"`
	Date     string `gorm:"size:255" json:"date"`
	Value    string `gorm:"size:255" json:"value"`
}

func Migrate() {
	DB.AutoMigrate(&User{}, &Transactions{})
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User, error) {
	var u User

	DB.First(&u, uid)
	u.PrepareGive()

	return u, nil
}

func GetUsernameByID(uid uint) (string, error) {
	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return "", errors.New("user not found")
	}

	return u.Username, nil
}

func GetAccountIDByUsername(username string) (string, error) {
	var u User

	if err := DB.Where("username = ?", username).First(&u).Error; err != nil {
		return "", errors.New("user not found")
	}

	return u.AccountID, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type GetTransactionsResponse struct {
	From  string `gorm:"size:255" json:"from"`
	To    string `gorm:"size:255" json:"to"`
	Date  string `gorm:"size:255" json:"date"`
	Value string `gorm:"size:255" json:"value"`
	Type  string `gorm:"size:255" json:"type"`
}

func GetTransactions(username string, accountID string) []GetTransactionsResponse {
	t1 := []Transactions{}
	t2 := []Transactions{}
	t3 := []Transactions{}

	DB.Find(&t1, "receiver = ?", username)
	DB.Find(&t2, "sender = ?", username)
	DB.Find(&t3, "receiver = ?", accountID)

	var JsonForm []GetTransactionsResponse
	var JsonElement GetTransactionsResponse

	t := append(t1, t2...)
	t = append(t, t3...)

	for i := 0; i < len(t); i++ {
		JsonElement.Date = t[i].Date
		JsonElement.From = t[i].Sender
		JsonElement.To = t[i].Receiver
		JsonElement.Value = t[i].Value

		if t[i].Sender == "System" {
			JsonElement.Type = "0"
		}
		if t[i].Sender == username {
			JsonElement.Type = "1"
		}
		if t[i].Receiver == accountID {
			JsonElement.Type = "2"
		}

		JsonForm = append(JsonForm, JsonElement)
	}

	return JsonForm
}

func LoginCheck(username string, password string) (string, error) {
	var err error
	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID, u.Email, u.AccountID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *User) SaveUser() (*User, error) {
	var err error = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (t *Transactions) SaveTransaction() (*Transactions, error) {
	var err error = DB.Create(&t).Error
	if err != nil {
		return &Transactions{}, err
	}

	return t, nil
}

func (u *User) BeforeSave() error {
	// Turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	// remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

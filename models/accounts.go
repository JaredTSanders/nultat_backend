package models

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	u "github.com/JaredTSanders/nultat_backend/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	FName       string `json:"first_name"`
	LName       string `json:"last_name"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	Standing    string `json:"standing"`
	Password    string `json:"password"`
	Token       string `json:"token";sql:"-"`
	Role        string `json:"role"`
	MFA_Enabled string `json:"mfa_enabled"`
	AccType     string `json:"account_type"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 8 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.Role = "user"
	account.Status = "active"
	account.Standing = "good"
	account.MFA_Enabled = "no"
	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string, w http.ResponseWriter) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	sessionToken := uuid.NewV4().String()
	_, err = Cache.Do("SETEX", sessionToken, "300", account.Email)

	if err != nil {
		return u.Message(false, "error setting cache")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(300 * time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: 2,
	})

	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	// tk := &Token{UserId: account.ID}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	return resp
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	response, err := Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewV4().String()
	_, err = Cache.Do("SETEX", newSessionToken, "300", fmt.Sprintf("%s", response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = Cache.Do("DEL", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    newSessionToken,
		Expires:  time.Now().Add(300 * time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: 2,
	})
}

// func GetCurrentUser(id uint) *BasicAccount {
// 	account := &BasicAccount{}
// 	err := GetDB().Table("accounts").Where("id = ?", id).First(account).Error
// 	if err != nil {
// 		return nil
// 	}
// 	return account
// }

func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func GetUserByEmail(e string) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("email = ? ", e).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}

func UpdateCurrentUser(account *Account) *Account {
	// if resp, ok := account.Validate(); !ok {
	// 	return resp
	// }
	return nil

	// GetDB().Update(account)

	// if account.ID <= 0 {
	// 	return u.Message(false, "Failed to create account, connection error.")
	// }

	// account.Password = "" //delete password

	// response := u.Message(true, "Account has been updated")
	// response["account"] = account
	// return response
}

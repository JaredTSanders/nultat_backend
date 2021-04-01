package models

import (
	"fmt"
	"net/http"
	"context"
	"strings"
	"time"
	"errors"

	u "github.com/JaredTSanders/nultat_backend/utils"
	uuid "github.com/satori/go.uuid"
	
	"github.com/getsentry/sentry-go"
    "github.com/honeycombio/libhoney-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/client-go/tools/clientcmd"


	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/kubernetes/pkg/api/v1"

	"k8s.io/client-go/kubernetes"
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
	Token       string `gorm:"-" ; json:"token"`
	Role        string `json:"role"`
	MFA_Enabled string `json:"mfa_enabled"`
	AccType     string `json:"account_type"`
	Namespace   string 
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
		sentry.CaptureMessage("Database connection failed")
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

	config, err := clientcmd.BuildConfigFromFlags("", *Kubeconfig)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig not found, thrown in BuildConfigFromFlags in accounts.go"))
        })
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig error, thrown in in NewForConfig in accounts.go"))
        })
		panic(err)
	}

	ns := "arkpc-" +fmt.Sprint(account.ID)

	GetDB().Table("accounts").Where("email = ?", account.Email).First(account).Update("namespace", ns)

	account.Namespace = ns

	kubeclient := client.CoreV1().Namespaces()

	// Create resource object
	object := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
		Spec: corev1.NamespaceSpec{},
	}

	// Manage resource
	_, err = kubeclient.Create(context.TODO(), object, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Namespace " + ns + " created successfully!")

	// nameSpace := 


	//Create new JWT token for the newly registered account
	// tk := &Token{UserId: account.ID}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	// account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string, w http.ResponseWriter) map[string]interface{} {

	
	start := time.Now()

	params := map[string]interface{}{
		"hostname": "api.jaredtsanders.com",
		"built": false,
		"user_id": -1,
	}

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return u.Message(false, "Email address not found")
		}
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Database connection failed, thrown in Login block"))
        })
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		w.WriteHeader(http.StatusUnauthorized)
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	sessionToken := uuid.NewV4().String()
	_, err = Cache.Do("SETEX", sessionToken, "3000", account.Email)

	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Redis SETEX failed, thrown in Login block"))
        })
		return u.Message(false, "error setting cache")
	}

	
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(3000 * time.Second),
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

	libhoney.Add(params)

	builder := libhoney.NewBuilder()

	builder.AddField("built", true)

	ev := builder.NewEvent()

	t := time.Now()

	elapsed := t.Sub(start)

	ev.AddField("user_id", account.Email)
	ev.AddField("latency_ms", elapsed)
	ev.Timestamp = time.Now()

    ev.Send()

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
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Redis DEL failed, thrown in Refresh block"))
        })
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
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Redis DEL failed, thrown in Refresh block"))
        })
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = Cache.Do("DEL", sessionToken)
	if err != nil {

		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Redis DEL failed, thrown in Refresh block"))
        })
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

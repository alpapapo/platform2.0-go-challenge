package models

import (
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	Assets	[]Asset	`gorm:"many2many:user_asset;"`

}

func (user *User) SetToken() {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationMinuntes, _ := strconv.Atoi( os.Getenv("JWT_TOKEN_EXPIRY_TIME_IN_MINUTES"))
	expirationTime := time.Now().Add(time.Duration(expirationMinuntes) * time.Minute)
	// Create the JWT claims, which includes the username and expiry time

	tk := &Token{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_TOKEN_PASSWORD")))
	user.Token = tokenString
}

//Validate incoming user details...
func (user *User) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return c.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 8 {
		return c.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return c.Message(false, "Email address already in use by another user."), false
	}

	return c.Message(false, "Requirement passed"), true
}

func (user *User) Create() (map[string]interface{}) {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword := c.SHA256OfString(user.Password)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return c.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	user.SetToken()

	user.Password = "" //delete password

	response := c.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func Login(email, password string) (map[string]interface{}) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Message(false, "Email address not found")
		}
		return c.Message(false, "Connection error. Please retry")
	}

	if c.SHA256OfString(password) != user.Password { //Password does not match!
		failedLoginMessage := "Unsuccessful login attempt from user: " + user.Username
		log.Println(failedLoginMessage)
		return c.Message(false, failedLoginMessage)
	}
	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	user.SetToken()

	successfulLoginMessage := user.Username + " logged in!"
	log.Println(successfulLoginMessage)
	resp := c.Message(true, successfulLoginMessage)
	resp["token"] = user.Token
	return resp
}


func GetUser(u uint) *User {

	user := &User{}
	GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Email == "" { //User not found!
		return nil
	}

	user.Password = ""
	return user
}


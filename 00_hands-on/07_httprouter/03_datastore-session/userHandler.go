package webapp

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type keyHash []byte

// User is how users accounts are formated
type User struct {
	Name       string
	Email      string
	Password   string
	Subscribed time.Time
	Role       int // 1 = User	 2 = Admin 	3 = Owner
}

// CheckAccount Is used to verify is an email is already used or it is available to be used. This returns true if the email is already registered.
func CheckAccount(r *http.Request, email string) (bool, error) {
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "User", email, 0, nil)
	if query, err := datastore.NewQuery("User").Ancestor(key).Count(ctx); err == nil {
		if query > 0 {
			return true, errors.New("Email is already been used")
		}
	} else {
		return false, err
	}

	return false, nil
}

// LoginAccount is used to execute a LOGIN action for any user account. This will return TRUE if every statement execute without any error and also returns the user Role.
func LoginAccount(r *http.Request, email string, password string) (logged bool, role int, err error) {
	password = HashPassword(password, nil)
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "User", email, 0, nil)
	// Check if email exists and match
	found, _ := CheckAccount(r, email)
	if found == true {
		var us User
		_, err := datastore.NewQuery("User").Ancestor(key).Run(ctx).Next(&us)
		if err != nil {
			return false, 0, err
		}
		if &us != nil && us.Password == password {
			log.Println("Account Match")
			return true, us.Role, nil
		}
		return false, 0, errors.New("Email or account invalid")
	}
	log.Println("Account not found")
	//Return nil if non error occurs
	return false, 0, errors.New("Account doesn't exists or maybe you have forggoten something")
}

// SetAccount is used to CREATE user account and insert it into DATASTORE
func SetAccount(r *http.Request, user User) error {
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "User", user.Email, 0, nil)
	// Check if email is already been used
	if _, account := CheckAccount(r, user.Email); account != nil {
		return account
	}
	// Encode user password
	user.Password = HashPassword(user.Password, nil)
	// Put user account into Datastore
	_, err := datastore.Put(ctx, key, &user)
	if err != nil {
		return err
	}
	return nil
}

// GetAccount Gets only one account per request and return an *User type with all information
func GetAccount(r *http.Request, email string) (us User, err error) {
	if found, err := CheckAccount(r, email); found == false {
		return us, err
	}
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "User", email, 0, nil)
	if err := datastore.Get(c, key, &us); err != nil {
		return us, err
	}
	return us, nil
}

// GetAllAccount TODO: This func will get all accounts from datastore and return it only for users with "Role" higher then 1
func GetAllAccount() {

}

// UpdateAccount This function will overwrite the entire Entity in datastore with the given data.
func UpdateAccount(r *http.Request, user User) (bool, error) {
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "User", user.Email, 0, nil)
	// Check if email exists and match
	found, _ := CheckAccount(r, user.Email)
	if found == true {
		// Put user account into Datastore
		if _, err := datastore.Put(ctx, key, &user); err != nil {
			return false, err
		}
		// All user data was updated
		return true, nil
	}
	//Return nil if non error occurs
	return false, errors.New("Account doesn't exists or maybe you have forggoten something")
}

// HashPassword Generate a secure hash for password
func HashPassword(password string, key keyHash) string {
	if key == nil {
		key = []byte("secret-key-here")
	}
	hash := append([]byte(password), key...)
	hasher := sha1.New()
	hasher.Write(hash)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

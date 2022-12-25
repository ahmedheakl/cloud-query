package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type Query struct {
	Sql    string `json:"query"`
	Schema string `json:"schema"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Purchase struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type Cart struct {
	ItemsRow
	Quantity int `json:"quantity"`
}

func (c *Cart) SetQuantity(quantity int) {
	c.Quantity = quantity
}

// check if provided credentials are signed up
func IsValidLogin(user User) (bool, error) {
	result, err := ReturnQuery(
		`SELECT 1
		FROM users
		WHERE email = $1 AND sha256((salt || $2)::bytea) = hashed_password::bytea`, "verification",
		user.Email, user.Password)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	var response Verification
	_ = json.Unmarshal(result, &response)
	return response.Value == 1, nil
}

// check that string provided is a valid email
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}

// add user to database
func AddUser(user User) error {
	salt := uuid.NewString()

	err := executeQuery(
		`INSERT INTO users(email,salt,hashed_password)
		VALUES($1,$2,sha256(($2 || $3)::bytea))`,
		user.Email, salt, user.Password)

	return err
}

// generate uuid string cookie and map cookie to email in redis cache
func SetCookie(w *http.ResponseWriter, email string) string {
	val := uuid.NewString()
	RegisterCookie(val, email)
	http.SetCookie(*w, &http.Cookie{
		Name:  "goCookie",
		Value: val,
		Path:  "/",
	})
	return val
}

// get value of cookie, ask to login if not found
func RequireCookie(w *http.ResponseWriter, r *http.Request) (string, error) {
	writer := *w
	c, err := r.Cookie("goCookie")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Please login first!"))
		return "", err
	}
	return c.Value, nil
}

// add item with quantity to user's cart
func AddToCart(pur Purchase, cookie string) (int, error) {
	return modifyCart(pur, cookie)
}

// remove quantity of item from cart,
// update quantity to zero or less to remove from cart
func RemoveFromCart(pur Purchase, cookie string) (int, error) {
	pur.Quantity = -pur.Quantity
	return modifyCart(pur, cookie)
}

// translate cookie to user email, and write the carts info to database
func WriteThrough(cookie string, items map[string]string, timestamp time.Time) error {
	email, err := TranslateCookie(cookie)
	if err != nil {
		return err
	}

	for item, quantity := range items {
		err = executeQuery(
			`INSERT INTO purchases(user_id, item_id, quantity, purchase_date)
			VALUES((SELECT id FROM users WHERE email = $1), $2, $3, $4)`,
			email, item, quantity, timestamp)
		if err != nil {
			return err
		}
	}
	return nil
}

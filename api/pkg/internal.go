package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	log "github.com/MohamedAbdeen21/cloud-store/logger"
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
	Item     int `json:"item"`
	Quantity int `json:"quantity"`
}

type Cart struct {
	ItemsRow
	Quantity int `json:"quantity"`
}

func (c *Cart) SetQuantity(quantity int) {
	c.Quantity = quantity
}

func InitConnections() error {
	_, err := getOrCreateCache()
	if err != nil {
		return err
	}

	_, err = getOrCreateCarts()
	if err != nil {
		return err
	}

	_, err = getOrCreateDB()
	if err != nil {
		return err
	}
	return nil
}

// check if provided credentials are signed up
func IsValidLogin(user User) (bool, string, error) {
	result, err := ReturnQuery(
		`SELECT 1
		FROM users
		WHERE email = $1 AND sha256((salt || $2)::bytea) = hashed_password::bytea`, "verification",
		user.Email, user.Password)
	if err != nil {
		log.Error.Printf("error while executing login query for email error '%s'", user.Email)
		return false, "internal error", err
	}

	var response Verification
	_ = json.Unmarshal(result, &response)

	log.Info.Println(response.Value)
	// see if the problem is from email or password
	if response.Value != 1 {
		result, err = ReturnQuery(
			`SELECT 1
				FROM users
				WHERE email = $1`, "verification", user.Email)
		if err != nil {
			return false, "internal error", err
		}

		var email_response Verification
		_ = json.Unmarshal(result, &email_response)
		// email exists, problem in password
		if email_response.Value == 1 {
			return false, "Incorrect password", nil
		} else {
			return false, "Email is not registered", nil
		}
	}
	return true, "", nil
}

// check that string provided is a valid email
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		log.Info.Printf("email provided '%s' is not a valid email", email)
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

	if err != nil {
		log.Info.Printf("sign up failed. Email '%s' already exists", user.Email)
	}

	return err
}

// generate uuid string cookie and map cookie to email in redis cache
func SetCookie(w *http.ResponseWriter, email string) string {
	val := uuid.NewString()
	if email != "remove" {
		RegisterCookie(val, email)
		http.SetCookie(*w, &http.Cookie{
			Name:     "goCookie",
			Value:    val,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Path:     "/",
		})
	} else {
		http.SetCookie(*w, &http.Cookie{
			Name:  "goCookie",
			Value: "",
			// SameSite: http.SameSiteNoneMode,
			// Secure:   true,
			Path:    "/",
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
		})
	}
	return val
}

// add cookie-to-email mapping to redis
func RegisterCookie(cookie string, email string) {
	rdb, _ := getOrCreateCache()
	rdb.Set(context.Background(), cookie, email, time.Hour*240)
}

func RemoveCookie(cookie string) {
	rdb, _ := getOrCreateCache()
	rdb.Del(context.Background(), cookie)
}

// get value of cookie, ask to login if not found
func RequireCookie(w *http.ResponseWriter, r *http.Request) (string, error) {
	writer := *w
	c, err := r.Cookie("goCookie")
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "Please login first!"})
		writer.Write(resp)
		return "", err
	}
	return c.Value, nil
}

// add item with quantity to user's cart
func AddToCart(pur Purchase, cookie string) (int, error) {
	return modifyCart(pur, cookie)
}

func GetCart(cookie string) (map[string]string, error) {
	rdb, _ := getOrCreateCarts()
	return rdb.HGetAll(context.Background(), cookie).Result()
}

// remove quantity of item from cart,
// update quantity to zero or less to remove from cart
func RemoveFromCart(pur Purchase, cookie string) (int, error) {
	pur.Quantity = -pur.Quantity
	return modifyCart(pur, cookie)
}

// translate cookie to user email, and write the carts info to database
func WriteCart(cookie string, items map[string]string, timestamp time.Time) error {
	email, err := translateCookie(cookie)
	if err != nil {
		return err
	}

	for item, quantity := range items {
		err = executeQuery(
			`INSERT INTO purchases(user_id, item_id, quantity, purchase_date)
			VALUES((SELECT id FROM users WHERE email = $1), $2, $3, $4)`,
			email, item, quantity, timestamp)
		if err != nil {
			log.Error.Printf("error while executing cart query for {%s,%s,%s,%s}", email, item, quantity, timestamp)
			return err
		}
	}

	RemoveCart(cookie)
	return nil
}

// remove the cart completely
func RemoveCart(cookie string) {
	rdb, _ := getOrCreateCarts()
	rdb.Del(context.Background(), cookie)
}

func SetHeaders(w *http.ResponseWriter, origin string, method string) {
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, withCredentials")
	(*w).Header().Set("Access-Control-Allow-Methods", method)
	(*w).Header().Set("Content-Type", "application/json")
}

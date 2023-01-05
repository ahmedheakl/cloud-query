package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/MohamedAbdeen21/cloud-store/logger"
	"github.com/MohamedAbdeen21/cloud-store/pkg"
	"github.com/lib/pq"
)

var frontendOrigin string = strings.Split(os.Getenv("frontendOrigin"), "=")[0]
var selfOrigin string = strings.Split(os.Getenv("APIOrigin"), "=")[0]

func InitConnections() error {
	return pkg.InitConnections()
}

// needs all queries to return all columns of a table AND in the same order,
// check structs at pkg/internal for the right order
// TODO: return only the required columns and in any order
func CustomQuery(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "POST")
	// Extract query from request body
	var query pkg.Query
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error.Printf("can't decode body of request, error: %s", err.Error())
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "invalid request body"})
		w.Write(resp)
		return
	}

	log.Info.Printf("Executing query %s", query.Sql)

	// execute query and return result
	resp, err := pkg.ReturnQuery(query.Sql, query.Schema)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error.Printf("can't execute query:%s error:%s", query.Sql, err.Error())
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "internal error"})
		w.Write(resp)
		return
	}

	w.Write(resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "POST")
	// decode request body into User struct instance `user`
	var user pkg.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error.Printf("can't decode body of request, error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "invalid request body"})
		w.Write(resp)
		return
	}

	log.Info.Printf("logging in for user '%s' ...", user.Email)

	// validate given email
	if !pkg.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "the email provided is not a valid email"})
		w.Write(resp)
		return
	}

	// validate login credentials
	isValid, errorMessage, err := pkg.IsValidLogin(user)
	if err != nil {
		log.Error.Printf("verification query for '%s' failed", user.Email)
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "internal error"})
		w.Write(resp)
		return
	}

	var resp []byte
	if !isValid {
		w.WriteHeader(http.StatusForbidden)
		log.Info.Printf("login for user '%s' failed, %s", user.Email, errorMessage)
		resp, _ = json.Marshal(map[string]any{"response": false, "message": errorMessage})
		// w.Write(resp)
	} else {
		// get value of cookie, generate new cookie if not found.
		var cookie string
		c, err := r.Cookie("goCookie")
		if err != nil {
			log.Info.Printf("login successful for user '%s', generated cookie '%s'", user.Email, cookie)
		} else {
			cookie = c.Value
			pkg.RegisterCookie(cookie, user.Email)
			log.Info.Printf("login successful for user '%s', reused cookie '%s'", user.Email, cookie)
		}

		if isValid {
			// http.Redirect(w, r, fmt.Sprintf("/cookie/%s/", user.Email), http.StatusSeeOther)
			w.WriteHeader(http.StatusAccepted)
			resp, _ = json.Marshal(map[string]any{"response": true, "message": "user logged in successfully"})
			// w.Write(resp)
		}
	}
	w.Write(resp)
}

// Logging out removes the cart
func Logout(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "GET")
	// Must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}
	pkg.RemoveCart(cookie)
	pkg.RemoveCookie(cookie)
	// pkg.SetCookie(&w, "remove")
	resp, _ := json.Marshal(map[string]any{"response": true, "message": "user logged out successfully"})
	w.Write(resp)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "POST")
	// decode request body into User struct instance `user`
	var user pkg.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error.Printf("can't decode body of request, error: %s", err.Error())
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "invalid request body"})
		w.Write(resp)
		return
	}

	// validate given email
	if !pkg.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "the email provided is not a valid email"})
		w.Write(resp)
		return
	}

	// add User, returns error if email already exists
	err = pkg.AddUser(user)

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "email already exists"})
		w.Write(resp)
	} else {
		// pkg.SetCookie(&w, user.Email)
		w.WriteHeader(http.StatusCreated)
		resp, _ := json.Marshal(map[string]any{"response": true, "message": "user created successfully"})
		w.Write(resp)
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "POST")
	// Decode request body into Purchase struct
	var pur pkg.Purchase
	err := json.NewDecoder(r.Body).Decode(&pur)
	if err != nil {
		log.Error.Printf("can't decode body of request, error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "invalid request body"})
		w.Write(resp)
		return
	}

	// Must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// add to cart, don't care about new quantity
	_, err = pkg.AddToCart(pur, cookie)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "internal error"})
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(map[string]any{"response": true, "message": "item added successfully"})
		w.Write(resp)
	}
}

func RemoveItem(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "POST")
	// decode request into Purchase struct
	var pur pkg.Purchase
	err := json.NewDecoder(r.Body).Decode(&pur)
	if err != nil {
		log.Error.Printf("can't decode body of request, error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "invalid request body"})
		w.Write(resp)
		return
	}

	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// remove from cart, don't care about new quantity
	_, err = pkg.RemoveFromCart(pur, cookie)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "internal error"})
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(map[string]any{"response": true, "message": "item removed successfully"})
		w.Write(resp)
	}
}

func CheckItems(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "GET")
	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// get all items and return as json
	cart, err := pkg.GetCart(cookie)
	ids := []int{}
	for id := range cart {
		int_id, _ := strconv.Atoi(id)
		ids = append(ids, int_id)
	}

	if len(ids) == 0 {
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(ids)
		w.Write(resp)
		return
	}

	// get all item info
	items, _ := pkg.ReturnQuery("SELECT * FROM items WHERE id = ANY($1)", "items", pq.Array(ids))
	var cartItems []pkg.Cart
	json.Unmarshal(items, &cartItems)

	for index, item := range cartItems {
		for id, quantity := range cart {
			int_id, _ := strconv.Atoi(id)
			int_quantity, _ := strconv.Atoi(quantity)
			if item.Id == int_id {
				cartItems[index].SetQuantity(int_quantity)
			}
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error.Printf("error fetching cart for user %s", cookie)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "internal error"})
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(cartItems)
		w.Write(response)
	}
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, frontendOrigin, "GET")
	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// get all items in shopping cart from redis
	items, _ := pkg.GetCart(cookie)
	if len(items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "cart is empty"})
		w.Write(resp)
		return
	}

	// write the items to the database and remove cart from cache
	err = pkg.WriteCart(cookie, items, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal(map[string]any{"response": false, "message": "error writing cart"})
		w.Write(resp)
		return
	}
	resp, _ := json.Marshal(map[string]any{"response": true, "message": "checkout successful"})
	w.Write(resp)
}

func GenerateCookie(w http.ResponseWriter, r *http.Request) {
	pkg.SetHeaders(&w, selfOrigin, "GET")
	var email = strings.Split(r.URL.Path, "/")[2]
	pkg.SetCookie(&w, email)
	if email != "remove" {
		http.Redirect(w, r, frontendOrigin, http.StatusFound)
	} else {
		http.Redirect(w, r, frontendOrigin+"/templates/auth.html", http.StatusSeeOther)
	}
}

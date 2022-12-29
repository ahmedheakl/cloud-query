package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MohamedAbdeen21/cloud-store/pkg"
	"github.com/lib/pq"
)

// needs all queries to return all columns of a table AND in the same order,
// check structs at pkg/internal for the right order
// TODO: return only the required columns and in any order
func CustomQuery(w http.ResponseWriter, r *http.Request) {
	// Extract query from request body
	var query pkg.Query
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
	}
	log.Printf("Executing query %s", query.Sql)

	// execute query and return result
	resp, err := pkg.ReturnQuery(query.Sql, query.Schema)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")

	// decode request body into User struct instance `user`
	var user pkg.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// get value of cookie, generate new cookie if not found.
	var cookie string
	c, err := r.Cookie("goCookie")
	if err != nil {
		cookie = pkg.SetCookie(&w, user.Email)
	} else {
		cookie = c.Value
		pkg.RegisterCookie(cookie, user.Email)
	}

	// validate given email
	if !pkg.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]string{"response": "the email provided is not a valid email"})
		w.Write(resp)
		return
	}

	// validate login credentials
	isValid, err := pkg.IsValidLogin(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid login")
		return
	}

	var resp []byte
	if isValid {
		w.WriteHeader(http.StatusAccepted)
		resp, _ = json.Marshal(map[string]any{"response": true, "cookie": cookie})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ = json.Marshal(map[string]bool{"response": false})
	}
	w.Write(resp)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	// decode request body into User struct instance `user`
	var user pkg.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	// validate given email
	if !pkg.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]string{"response": "the email provided is not a valid email"})
		w.Write(resp)
		return
	}

	// add User, returns error if email already exists
	err = pkg.AddUser(user)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		pkg.SetCookie(&w, user.Email)
		w.WriteHeader(http.StatusCreated)
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	// Decode request body into Purchase struct
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	var pur pkg.Purchase
	err := json.NewDecoder(r.Body).Decode(&pur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// add to cart, don't care about new quantity
	_, err = pkg.AddToCart(pur, cookie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func RemoveItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	// decode request into Purchase struct
	var pur pkg.Purchase
	err := json.NewDecoder(r.Body).Decode(&pur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// remove from cart, don't care about new quantity
	_, err = pkg.RemoveFromCart(pur, cookie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func CheckItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rdb := pkg.GetOrCreateCarts()

	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// get all items and return as json
	cart, err := rdb.HGetAll(context.Background(), cookie).Result()
	ids := []int{}
	for id := range cart {
		int_id, _ := strconv.Atoi(id)
		ids = append(ids, int_id)
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(cartItems)
		w.Write(response)
	}
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Cpmtemt-Type, withCredentials")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rdb := pkg.GetOrCreateCarts()

	// must have cookie
	cookie, err := pkg.RequireCookie(&w, r)
	if err != nil {
		return
	}

	// get all items in shooping cart from redis
	items := rdb.HGetAll(context.Background(), cookie).Val()
	for k, v := range items {
		println(k, v)
	}

	if len(items) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cart is empty"))
		return
	}

	// write the items to the database
	err = pkg.WriteThrough(cookie, items, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}

	// remove cart from redis if write was successful
	pkg.RemoveCart(cookie)
}

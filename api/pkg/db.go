package pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/MohamedAbdeen21/cloud-store/logger"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // import declerations/initializations only
)

var host string = strings.Split(os.Getenv("host"), "=")[0]
var user string = strings.Split(os.Getenv("username"), "=")[0]
var password string = strings.Split(os.Getenv("password"), "=")[0]
var database string = strings.Split(os.Getenv("database"), "=")[0]
var redisConnection string = strings.Split(os.Getenv("redis"), "=")[0]

var conn string = fmt.Sprintf("host=%s port=5432 user=%s password=%s database=%s", host, user, password, database)

var db *sql.DB
var dbCache *redis.Client
var cartsCache *redis.Client

// `users` table
type UserRow struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	Salt           string `json:"salt"`
	HashedPassword string `json:"hashed_password"`
}

// `purchases` table
type PurchasesRow struct {
	User_id       int       `json:"user_id"`
	Item_id       int       `json:"item_id"`
	Quantity      int       `json:"quantity"`
	Purchase_date time.Time `json:"purchase_date"`
}

// `items` table
type ItemsRow struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
}

// to store sign-in return value
type Verification struct {
	Value int `json:"value"`
}

// general struct to return result of a query
type Response struct {
	Users    []UserRow
	Purchase []PurchasesRow
	Items    []ItemsRow
	Verify   Verification
}

// general RDS postgreSQL database
func getOrCreateDB() (*sql.DB, error) {
	var err error
	if db == nil {
		log.Info.Println("connecting to database ...")
		db, err = sql.Open("postgres", conn)
		if err != nil {
			log.Fatal.Printf("can't connect to database, error: %s", err.Error())
			return nil, err
		}
		_, err := db.Exec("SELECT 1")
		if err != nil {
			log.Fatal.Printf("can't connect to database, error: %s", err.Error())
			return nil, err
		}
		log.Info.Println("successfully connected to database!")
	}
	return db, nil
}

// used for cookie-to-email translation
func getOrCreateCache() (*redis.Client, error) {
	if dbCache == nil {
		dbCache = redis.NewClient(&redis.Options{
			Addr:     redisConnection,
			Password: "",
			DB:       0, // default database, takes values from 0-15 inclusive
		})
		_, err := dbCache.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal.Printf("can't connect to redis, error: %s", err.Error())
			return nil, fmt.Errorf("can't connect to redis: %s", err.Error())
		}
		log.Info.Println("successfully connected to redis!")
	}
	return dbCache, nil
}

// used to cache shopping carts, shopping carts are stored in
// primary database only when the purchase is made
func getOrCreateCarts() (*redis.Client, error) {
	if cartsCache == nil {
		cartsCache = redis.NewClient(&redis.Options{
			Addr:     redisConnection,
			Password: "",
			DB:       1, // Database for shopping carts
		})
		_, err := cartsCache.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal.Printf("can't connect to redis, error: %s", err.Error())
			return nil, fmt.Errorf("can't connect to redis: %s", err.Error())
		}
		log.Info.Println("successfully connected to redis!")
	}
	return cartsCache, nil
}

// used to execute a query that returns a result
func ReturnQuery(query string, table string, args ...any) ([]byte, error) {
	var response Response

	db, err := getOrCreateDB()
	if err != nil {
		return nil, err
	}

	// execute query, pass arguments if given
	var rows *sql.Rows
	if len(args) != 0 {
		rows, err = db.Query(query, args...)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}

	// which struct to use, expects all columns to be returned from the query
	// return all columns AND in order to avoid errors
	for rows.Next() {
		switch table {
		case "users":
			row := UserRow{}
			if err := rows.Scan(&row.Id, &row.Email, &row.Salt, &row.HashedPassword); err != nil {
				val, _ := json.Marshal(response.Users)
				return val, err
			}
			response.Users = append(response.Users, row)
		case "purchases":
			row := PurchasesRow{}
			if err := rows.Scan(&row.User_id, &row.Item_id, &row.Quantity, &row.Purchase_date); err != nil {
				val, _ := json.Marshal(response.Users)
				return val, err
			}
			response.Purchase = append(response.Purchase, row)
		case "items":
			row := ItemsRow{}
			if err := rows.Scan(&row.Id, &row.Name, &row.Brand, &row.Description, &row.Price, &row.Image); err != nil {
				val, _ := json.Marshal(response.Users)
				return val, err
			}
			response.Items = append(response.Items, row)
		case "verification":
			row := Verification{}
			if err := rows.Scan(&row.Value); err != nil {
				return nil, err
			}
			response.Verify = row
		default:
			return nil, fmt.Errorf("schema provided is unknown")
		}
	}

	if err = rows.Err(); err != nil {
		println(err.Error())
	}

	// which struct to return
	switch table {
	case "users":
		return json.Marshal(response.Users)
	case "purchases":
		return json.Marshal(response.Purchase)
	case "items":
		return json.Marshal(response.Items)
	case "verification":
		return json.Marshal(response.Verify)
	default:
		return nil, fmt.Errorf("schema provided is unknown")
	}
}

// used to execute a query that doesn't return a result
func executeQuery(query string, args ...any) error {
	db, err := getOrCreateDB()
	if err != nil {
		return err
	}

	// execute query and pass arguments if given
	if len(args) != 0 {
		_, err = db.Exec(query, args...)
	} else {
		_, err = db.Exec(query)
	}
	return err
}

// modify content of cart, can be used to add or remove items
func modifyCart(pur Purchase, cookie string) (int, error) {
	rdb, _ := getOrCreateCarts()
	newVal, err := rdb.HIncrBy(context.Background(), cookie, strconv.Itoa(pur.Item), int64(pur.Quantity)).Result()

	// cart not found, create it
	if err == redis.Nil {
		_, err = rdb.HSet(context.Background(), cookie, pur.Item, pur.Quantity).Result() // returns 1 regardless of initial value
		if err != nil {
			log.Error.Printf("error creating cart: %s", err.Error())
			return 0, err
		}
		newVal = int64(pur.Quantity)
		rdb.Expire(context.Background(), cookie, time.Hour*24) // expires in 24 hours
	}

	// if quantity of an item is <= 0, remove the item from the cart
	if newVal <= 0 {
		rdb.HDel(context.Background(), cookie, strconv.Itoa(pur.Item))
	}

	return int(newVal), nil
}

// translate cookie into email
func translateCookie(cookie string) (string, error) {
	rdb, _ := getOrCreateCache()
	val, err := rdb.Get(context.Background(), cookie).Result()
	if err == redis.Nil {
		return val, fmt.Errorf("can't find email for cookie %s", cookie)
	}
	return val, err
}

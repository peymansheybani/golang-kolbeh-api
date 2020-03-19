package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Hotel struct {
	Id            int
	Name          string
	Checkin_time  string
	Checkout_time string
	Address       string
	City          City
	User          User
	Description   string
	Period        interface{}
	Infos         interface{}
	General       interface{}
	Possible      interface{}
	Services      interface{}
	ServicesImage interface{}
	Rooms         []Room
}

type City struct {
	City_id   int
	Code      string
	City_name string
}

type User struct {
	User_id int
	Email   string
}

type Room struct {
	Id                 int
	Hotel_id           int
	Name               string
	Size               string
	Extra              int
	Description        string
	Meals              string
	Dynamic_data       interface{}
	Availability_free  interface{}
	Availability_price interface{}
	Images             interface{}
	ShowImage          interface{}
	Board              interface{}
	Price              string
	Possibles          interface{}
}

func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "21312"
	dbName := "kolbeh"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetRooms(Hotel_id int) []Room {
	db := DbConn()

	fmt.Println("HOTEL ID IS =", Hotel_id)
	selDB, err := db.Query(`SELECT r.id,r.name,r.extra,r.size,
							r.description,r.meals,r.dynamic_data,
							r.availability_free,r.availability_price,r.hotel_id
							FROM rooms AS r WHERE r.hotel_id = ?`, Hotel_id)
	if err != nil {
		panic(err.Error())
	}
	emp := Room{}
	res := []Room{}
	for selDB.Next() {
		var id int
		var name string
		var extra int
		var size string
		var meals string
		var description string
		var dynamic_data string
		var availability_free string
		var availability_price string
		var hotel_id int
		var err error
		err = selDB.Scan(&id, &name, &extra, &size, &meals, &description, &dynamic_data,
			&availability_free, &availability_price, &hotel_id)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("room data => ", id, name)
		emp.Id = id
		emp.Name = name
		emp.Extra = extra
		emp.Hotel_id = hotel_id
		emp.Meals = meals
		emp.Size = size
		emp.Availability_free = availability_free
		emp.Availability_price = availability_price
		emp.Description = description

		var result map[string]interface{}
		json.Unmarshal([]byte(dynamic_data), &result)
		emp.Dynamic_data = result
		emp.Images = result["images"]
		emp.ShowImage = result["shoImage"]
		emp.Possibles = result["possibles"]

		res = append(res, emp)
	}

	defer db.Close()

	return res
}

func GetHotel() []Hotel {
	db := DbConn()
	hotels, err := db.Query(`SELECT h.id,h.name,h.checkin_time,h.checkout_time,
							h.address,h.city_id, c.name AS city_name,c.code,
							h.user_id,u.email,h.description,h.dynamic_data
							FROM hotels AS h JOIN cities AS c ON h.city_id = c.id
							INNER JOIN users AS u ON u.id = h.user_id`)

	if err != nil {
		panic(err.Error())
	}

	customHotels := Hotel{}
	response := []Hotel{}

	for hotels.Next() {
		var id int
		var name, checkin_time, checkout_time, address string
		// var []city{id : "",name : "", code : ""}
		var city_id int
		var city_name string
		var code string
		var user_id int
		var email string
		var description string
		var dynamic_data string
		// var room_id int
		// var room_name string
		// var extra int
		// var size int
		// var meals int
		// var room_description string
		// var d_data string
		// var availability_free string
		// var availability_price string
		err := hotels.Scan(&id, &name, &checkin_time, &checkout_time,
			&address, &city_id, &city_name, &code, &user_id,
			&email, &description, &dynamic_data)
		if err != nil {
			panic(err.Error())
		}

		customHotels.Id = id
		customHotels.Name = name
		customHotels.Checkout_time = checkout_time
		customHotels.Checkin_time = checkin_time
		customHotels.Address = address
		customHotels.City.City_id = city_id
		customHotels.City.City_name = city_name
		customHotels.City.Code = code
		customHotels.User.User_id = user_id
		customHotels.User.Email = email
		customHotels.Description = description

		var result map[string]interface{}
		json.Unmarshal([]byte(dynamic_data), &result)

		customHotels.Period = result["period"]
		customHotels.Infos = result["infos"]
		customHotels.General = result["general"]
		customHotels.Possible = result["possible"]
		customHotels.Services = result["services"]
		customHotels.ServicesImage = result["servicesImage"]

		rooms := GetRooms(id)
		fmt.Println("Rooms Pami joon IS =>>>>>>", rooms)
		customHotels.Rooms = rooms
		response = append(response, customHotels)
	}

	defer db.Close()

	return response
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
	fmt.Println(GetHotel())
	json.NewEncoder(w).Encode(GetHotel())
}

func main() {
	fmt.Println(GetHotel())
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

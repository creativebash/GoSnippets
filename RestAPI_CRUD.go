package main

// import packages
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	// loading the driver anonymously,
	// *aliasing its package qualifier to _
	// *so none of its exported names are visible to our code.
	_ "github.com/lib/pq"
)

// declare a struct of type User
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Sex          string `json:"sex"`
	Date_created string `json:"date_created"`
}

// declare and assign postgres db connection detail
// *Constants are declared like variables, but with the const keyword.
// *Constants can be character, string, boolean, or numeric values.
// *Constants cannot be declared using the := syntax.
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "peacemaker"
	dbname   = "godb"
)

// a function that sets up the database and establishes connection
func setupDB(user, password, host, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// the main function that get executed first when the program is run. the start point of the app
func main() {
	// creating the database object and opening connection
	db, err := setupDB(user, password, host, dbname)
	if err != nil {
		panic(err)
	}

	// Defering the closing of the database connection to give room for further query
	defer db.Close()

	// executing an SQL query statement
	_, err = db.Exec("DROP table if EXISTS users")
	if err != nil {
		panic(err)
	}

	// To check right away that the database is available and accessible
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully pinged the db")

	// an SQL query statement that creates 'users' table if not exists
	createUserTable := `
		CREATE table IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username varchar NOT NULL,
		email varchar NOT NULL,
		firstname varchar,
		lastname varchar,
		sex varchar,
		date_created timestamptz DEFAULT CURRENT_TIMESTAMP
		);
	`
	// executing the SQL query statement
	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(err)
	}
	fmt.Println("users table created successfully!")

	// inserting records into users table
	// assigning the SQL query to insertUser variable
	insertUser := `
		INSERT into users (
			username,
			email,
			firstname,
			lastname,
			sex
		) VALUES
		('bash', 'anakobembash@gmail.com', 'Bashir', 'Anakobe', 'male'),
		('teemah', 'teemah247@gmail.com', 'Fatimah', 'Muhammed', 'female'),
		('wasman', 'wasman01@gmail.com', 'Abdulwasiu', 'Anakobe', 'male'),
		('medo', 'ahmed123@gmail.com', 'Ahmed', 'Ibrahim', 'male'),
		('zain', 'zainyray@gmail.com', 'Zainab', 'Idris', 'female'),
		('stacia', 'cheerfulann@gmail.com', 'Anastasia', 'Ugwu', 'female')
		;
	`
	// executing the query
	_, err = db.Exec(insertUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("users added successfully!")

	// view records from the users table on the command line
	// * This script uses the db.Query function to execute
	// * a SELECT statement to retrieve all records from
	// * the "users" table. The rows.Next function is then
	// * used in a loop to iterate through the returned rows,
	// * and the rows.Scan function is used to extract
	// * the values for each column into variables.
	// * You can also use sql.QueryRow if you expect only one row
	// * to be returned by your query.
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var email string
		var firstname string
		var lastname string
		var sex string
		var date_created string
		// copying columns (fields) in the current row into the address destination
		// * Scan copies the columns in the current row into the values pointed at by dest.
		// * The number of values in dest must be the same as the number of columns in Rows.
		// * Scan converts columns read from the database into the following common
		// * Go types and special types provided by the sql package:
		err = rows.Scan(&id, &username, &email, &firstname, &lastname, &sex, &date_created)
		if err != nil {
			panic(err)
		}
		fmt.Println("ID:", id, "Username:", username, "Email:", email, "Fisrtname:", firstname, "Lastname:", lastname, "Sex:", sex, "Date Created:", date_created)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	// * http.HandleFunc et al apply your handlers to a package-global
	// * instance of the ServeMux held in the http package,
	// * which http.ListenAndServe then starts.

	// http.HandleFunc("/users", handleViewUsers)
	// http.HandleFunc("/newuser", handleCreateUser)
	// fmt.Println("Server listening on port 3000...")
	// http.ListenAndServe(":3000", nil)

	// * You can also create your own instance which gives you some
	// * more control and makes it easier to unit test.
	// * ServeMux is an HTTP request multiplexer. It is used for request routing and dispatching.
	// * The NewServeMux function allocates and returns a new ServeMux.

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handleViewUsers)
	mux.HandleFunc("/newuser", handleCreateUser)
	mux.HandleFunc("/userupdate", handleUpdateUser)
	mux.HandleFunc("/deleteuser", handleDeleteUser)
	fmt.Println("Server listening on port 3000...")
	http.ListenAndServe(":3000", mux)
} // end of func main()

// creating users slice to store user data
// * "users" variable becomes a slice of struct User
var users []User

// creating a handler function for users' GET request
func handleViewUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// connecting to database
	db, err := setupDB(user, password, host, dbname)
	if err != nil {
		panic(err)
	}
	// Defering the closing of the database connection to give room for further query
	defer db.Close()

	// SELECT statement to retrieve all records from the "users" table
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterating through the returned rows and extracting the values for each column into variables
	for rows.Next() {
		var user User // * user is instance of Sruct User
		// checking if error occurs while concurrently extracting data
		// * "Scan the rows and assign the values to the variables in user struct,
		// * if there is an error return an internal server error"
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.Sex, &user.Date_created); err != nil {
			// sending an HTTP error response with the appropriate status code and message,
			//  and then returning from the function to stop further execution
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Encoding the users data to json and sending it as response
	// * json.NewEncoder(w).Encode(users) is used to encode the users variable
	// * as JSON and write it to the io.Writer type variable w. The json.NewEncoder(w)
	// * creates a new encoder that writes to the w variable, and
	// * the .Encode(users) method is called on that encoder to
	// * write the JSON encoding of users to w. This is typically used to
	// * write the JSON representation of a Go data structure
	// * to an HTTP response body or to a file on disk.
	json.NewEncoder(w).Encode(users)

	// setting slice users to empty to avoid duplicate append on reload
	users = nil
}

// creating a handler function for new user POST request
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// * API keys are a common way to authenticate and authorize access to your API.
	// * One way to implement API keys in your Go application would be to
	// * include a header in the HTTP request called "API-KEY" and check its value
	// * in your request handler before processing the request.
	// The following code checks for the presence of
	// the "API-KEY" header and its value before processing further request

	apiKey := r.Header.Get("API-KEY")
	if apiKey != "your_api_key" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Connecting to the database
	db, err := setupDB(user, password, host, dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Parsing the request body as json
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Inserting the new user into the "users" table
	_, err = db.Exec("INSERT INTO users (username, email, firstname, lastname, sex) VALUES ($1, $2, $3, $4, $5)", newUser.Username, newUser.Email, newUser.Firstname, newUser.Lastname, newUser.Sex)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Sending a response indicating success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// *Note! You cannot set the request body r.Body directly in the
// *handleCreateUser handler to test. The r.Body is an io.ReadCloser
// *and is populated by the http server when it receives a request.
// *To test the handleCreateUser handler, you would need to simulate
// *an HTTP request to the server, and include the JSON payload in the
// *request body. This can be done using a testing framework like
// *net/http/httptest package in go or a HTTP client library
// *like net/http or github.com/golang/go/httptest.

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Connecting to the database
	db, err := setupDB(user, password, host, dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Parsing the request body as json
	var userUpdateData User
	if err := json.NewDecoder(r.Body).Decode(&userUpdateData); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Update the record with ID 3
	// u := User{ID: 3, Username: "newusername", Email: "newemail@example.com", Firstname: "New", Lastname: "Name", Sex: "male"}
	err = updateUser(db, userUpdateData)
	if err != nil {
		panic(err)
	}
	fmt.Println("User updated successfully!")

	// Sending a response indicating success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userUpdateData)
}

func updateUser(db *sql.DB, u User) error {
	query := `
		UPDATE users SET
			username = $1,
			email = $2,
			firstname = $3,
			lastname = $4,
			sex = $5
		WHERE id = $6;
	`

	_, err := db.Exec(query, u.Username, u.Email, u.Firstname, u.Lastname, u.Sex, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Connecting to the database
	db, err := setupDB(user, password, host, dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Parsing the request body as json
	var userDateDel User
	if err := json.NewDecoder(r.Body).Decode(&userDateDel); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Update the record with ID 3
	// u := User{ID: 3, Username: "newusername", Email: "newemail@example.com", Firstname: "New", Lastname: "Name", Sex: "male"}
	err = deleteUser(db, userDateDel)
	if err != nil {
		panic(err)
	}
	fmt.Println("User updated successfully!")

	// Sending a response indicating success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userDateDel)
}

func deleteUser(db *sql.DB, u User) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, u.ID)
	if err != nil {
		return err
	}

	return nil
}

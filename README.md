### Reflect ORM 

[![reflect-orm](https://github.com/WolvenSpirit/reflect-orm/actions/workflows/go.yml/badge.svg)](https://github.com/WolvenSpirit/reflect-orm/actions/workflows/go.yml)

Minimal proof of concept for an object relational mapper that can map data from tables without having to interact with the database beforehand.

It works by providing a `struct` that describes what is believed to be the data in the table.

It's meant to be simple with minimal boilerplate code.


```go

import (
    "fmt"
    "github.com/WolvenSpirit/reflectorm"
)

var db *sql.DB

func connectDB(user, pass, host, database, sslMode string) {
	var err error
	driverName := "postgres"
	url := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", driverName, user, pass, host, database, sslMode)
	if db, err = sql.Open(driverName, url); err != nil {
		fmt.Println("sql.Open: ", err.Error())
	}
}

// Define your table data, the struct doesn't need to be comprehensive
type Users struct {
    Id          int64
    Name        string
    Email       string
    // ...
}

func init() {
    // ... init db with credentials and other logic
}

func main() {
    // We use 'users' as our table definition
    var users Users
    // To be explicit regarding what is returned we define the result beforehand
    var result []Users
    // Get the users from the table
    result = Get(users, nil)
    // That's it!
    fmt.Println(result)
}

```


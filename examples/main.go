package reflectorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type TT struct {
	Id        int64
	Publisher string
	Channel   string
	Consumer  string
	Ack       bool
	Data      string
	Created   time.Time
	Duration  []uint8
	Completed time.Time
}
type notify_events TT

func connectDB(user, pass, host, database, sslMode string) {
	var err error
	driverName := "postgres"
	url := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", driverName, user, pass, host, database, sslMode)
	if db, err = sql.Open(driverName, url); err != nil {
		fmt.Println("sql.Open: ", err.Error())
	}
}

func MapFields[T interface{}](rows *sql.Rows, v T, rawValues []interface{}) T {
	names, _ := rows.ColumnTypes()
	t := reflect.TypeOf(v)
	rv := reflect.ValueOf(&v).Elem()
	nf := t.NumField()
	m := make(map[string]interface{}, nf)
	// Assuming rawValues are returned consistently, in the same order as the names
	for k := range names {
		// fmt.Printf("matching %s %s\n", names[k].Name(), *rawValues[k].(*interface{}))
		m[strings.Title(names[k].Name())] = *rawValues[k].(*interface{})
	}
	for k := range names {
		name := strings.Title(names[k].Name())
		sf, ok := t.FieldByName(strings.Title(name))
		sv := (rv).FieldByName(strings.Title(name))
		if !ok {
			fmt.Println("!! Failed to get struct field")
		}
		if sf.Type == reflect.TypeOf(m[name]) && sf.Name == name && sv.CanAddr() {
			sv.Set(reflect.ValueOf(m[name]))
		} else {
			fmt.Println("!!", sf.Type, reflect.TypeOf(m[name]))
		}
	}
	return v
}

func Get[T interface{}](model T, where []string) (rows []T) {
	var modelName string
	var fields []string
	queryStart := "select "
	if t := reflect.TypeOf(model); t.Kind() == reflect.Struct {
		modelName = t.Name()
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, t.Field(i).Name)
			if i+1 != t.NumField() {
				queryStart += fields[i] + ", "
			} else {
				queryStart += fields[i] + " "
			}
		}
		queryStart += "from " + modelName + " "
		if where != nil {
			queryStart += "where "
			for k := range where {
				queryStart += where[k]
			}
		}
	}
	//---
	rws, err := db.Query(queryStart)
	if err != nil {
		fmt.Println(err.Error())
	}
	var results []T
	for rws.Next() {
		var t T
		names, err := rws.Columns()

		columns := make([]interface{}, len(names))
		values := make([]interface{}, len(names))
		for k := range names {
			values[k] = &columns[k]
		}

		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println("result", names)
		rws.Scan(values...)
		t = MapFields[T](rws, t, values)
		fmt.Printf("%+v", t)
		results = append(results, t)
	}
	fmt.Println(queryStart)
	return results
}

func main() {
	connectDB("postgres", "", "localhost:5432", "shared_db01", "disable")
	/*rows, err := db.Query(`select
	id, publisher, channel, consumer,
	ack, data, created, duration, completed
	 from notify_events`)
	if err != nil {
		fmt.Println(err.Error())
	}
	var results []TT
	for rows.Next() {
		t := TT{}
		names, err := rows.Columns()

		columns := make([]interface{}, len(names))
		values := make([]interface{}, len(names))
		for k := range names {
			values[k] = &columns[k]
		}

		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println("result", names)
		rows.Scan(values...)
		t = MapFields[TT](rows, t, values)
		fmt.Printf("%+v", t)
		results = append(results, t)
	}*/
	var tt notify_events
	Get(tt, nil)
}

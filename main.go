package reflectorm

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

var store sync.Map

type JoinTable[model interface{}] struct {
	Model model
	On    []string
}

type QueryOptions[model interface{}] struct {
	Where []string
	Join  JoinTable[model]
}

type CacheResult[T interface{}] struct {
	Result T
	Error  error
}

/*
MapFields attempts to map each row data as an instance of type T
*/
func MapFields[T interface{}](rows *sql.Rows, v T, rawValues []interface{}) T {
	names, _ := rows.ColumnTypes()
	t := reflect.TypeOf(v)
	rv := reflect.ValueOf(&v).Elem()
	nf := t.NumField()
	m := make(map[string]interface{}, nf)
	// Assuming rawValues are returned consistently, in the same order as the names
	for k := range names {
		m[strings.Title(names[k].Name())] = *rawValues[k].(*interface{})
	}
	for k := range names {
		name := strings.Title(names[k].Name())
		sf, ok := t.FieldByName(strings.Title(name))
		sv := (rv).FieldByName(strings.Title(name))
		if !ok {
			os.Stderr.Write([]byte("!! Failed to get struct field"))
		}
		typeOfVal := reflect.TypeOf(m[name])
		if sf.Type == typeOfVal && sf.Name == name && sv.CanAddr() {
			sv.Set(reflect.ValueOf(m[name]))
		} else {
			fmt.Println("!! mismatch ", sf.Type, typeOfVal)
		}
	}
	return v
}

/*
Get constructs the query from the provided definition struct, performs the query and attempts to call MapFields for each row returned.
*/
func Get[T interface{}](db *sql.DB, model T, where []string) (rows []T, err error) {
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
	if rws != nil {
		defer rws.Close()
	}
	if err != nil {
		return nil, err
	}
	var results []T
	for rws.Next() {
		var t T
		names, err := rws.Columns()
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
			return nil, err
		}
		columns := make([]interface{}, len(names))
		values := make([]interface{}, len(names))
		for k := range names {
			values[k] = &columns[k]
		}
		rws.Scan(values...)
		t = MapFields[T](rws, t, values)
		results = append(results, t)
	}
	return results, nil
}

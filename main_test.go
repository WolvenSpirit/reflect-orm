package reflectorm

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type args[T any] struct {
	model T
	where []string
}
type Test[T any] struct {
	name     string
	args     args[T]
	wantRows []T
	wantErr  bool
}

type test_table struct {
	id     int
	Field1 string
	Field2 time.Time
}

type test_table3 struct {
	id     int
	Field2 string
	Field3 int
	Field4 []uint8
	Field5 float64
}

type wrong_table test_table3

var db *sql.DB

func sqliteNewDB() {
	var err error
	if db, err = sql.Open("sqlite3", "./test.db"); err != nil {
		panic(err.Error())
	}
}

func TestGet(t *testing.T) {
	sqliteNewDB()
	var ne test_table
	tests := []Test[test_table]{
		{
			args: args[test_table]{model: ne,
				where: []string{"id=1"}},
			wantRows: []test_table{{Field1: "test3"}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, err := Get(db, tt.args.model, tt.args.where)
			if err != nil && !tt.wantErr {
				t.Errorf(err.Error())
			}
			if tt.wantErr && err == nil {
				t.Fail()
			}
			if err == nil && gotRows[0].Field1 != tt.wantRows[0].Field1 {
				t.Errorf("Get() = %v, want %v", gotRows, tt.wantRows)
			}
		})
	}
	db.Close()
}

func TestGet2(t *testing.T) {
	sqliteNewDB()
	var ne test_table3
	tests := []Test[test_table3]{
		{
			args: args[test_table3]{model: ne,
				where: []string{"id=5"}},
			wantRows: []test_table3{{Field2: "test3", Field3: 3.0}},
			wantErr:  false,
		}, {
			args: args[test_table3]{model: ne,
				where: []string{"wrong"}},
			wantRows: []test_table3{{Field2: "test3", Field3: 3.0}},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, err := Get(db, tt.args.model, tt.args.where)
			if err != nil && !tt.wantErr {
				t.Errorf(err.Error())
			}
			if tt.wantErr && err == nil {
				t.Error("wantErr=true but no error")
			}
			if err == nil {
				if gotRows[0].Field2 != tt.wantRows[0].Field2 {
					t.Errorf("Get() = %v, want %v", gotRows, tt.wantRows)
				}
			}
		})
	}
	db.Close()
}

func TestGet3(t *testing.T) {
	sqliteNewDB()
	var ne wrong_table
	tests := []Test[wrong_table]{
		{
			args: args[wrong_table]{model: ne,
				where: []string{"id=5"}},
			wantRows: []wrong_table{{Field2: "test3", Field3: 3.0}},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, err := Get(db, tt.args.model, tt.args.where)
			if err != nil && !tt.wantErr {
				t.Errorf(err.Error())
			}
			if tt.wantErr && err == nil {
				t.Error("wantErr=true but no error")
			}
			if err == nil && gotRows[0].Field2 != tt.wantRows[0].Field2 {
				t.Errorf("Get() = %v, want %v", gotRows, tt.wantRows)
			}
		})
	}
	db.Close()
}

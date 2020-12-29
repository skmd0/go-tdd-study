package setup

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	// 0. flag.Parse() if you need flags
	exitCode := run(m)
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	const (
		dropDB          = `DROP DATABASE IF EXISTS test_user_store;`
		createDB        = `CREATE DATABASE test_user_store;`
		createUserTable = `CREATE TABLE users (
												 id SERIAL PRIMARY KEY,
												 name TEXT,
												 email TEXT UNIQUE NOT NULL
											 );`
	)

	psql, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=testtest sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}

	_, err = psql.Exec(dropDB)
	if err != nil {
		panic(fmt.Errorf("psql.Exec() err = %s", err))
	}
	_, err = psql.Exec(createDB)
	if err != nil {
		panic(fmt.Errorf("psql.Exec() err = %s", err))
	}
	// teardown
	defer func() {
		_, err = psql.Exec(dropDB)
		if err != nil {
			panic(fmt.Errorf("psql.Exec() err = %s", err))
		}
	}()

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=testtest sslmode=disable dbname=test_user")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}
	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(fmt.Errorf("db.Exec() err = %s", err))
	}
	return m.Run()
}

func TestUserStore(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=testtest sslmode=disable dbname=test_user")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err = %s", err))
	}
	us := &UserStore{
		sql: db,
	}
	t.Run("Find", testUserStore_Find(us))
	t.Run("Create", testUserStore_Find(us))
	t.Run("Delete", testUserStore_Find(us))
	t.Run("Subscribe", testUserStore_Find(us))
	// teardown
}

func testUserStore_Find(us *UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		usr := &User{
			Name:  "Domen Skamlic",
			Email: "domen@test.com",
		}
		err := us.Create(usr)
		if err != nil {
			t.Errorf("us.Create() err = %s", err)
		}
		defer func() {
			err := us.Delete(usr.ID)
			if err != nil {
				t.Errorf("us.Delete() err = %s", err)
			}
		}()

		tests := []struct {
			name    string
			id      int
			want    *User
			wantErr error
		}{
			{"Found", usr.ID, usr, nil},
			{"Not Found", -1, nil, ErrNotFound},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, err := us.Find(tc.id)
				if err != tc.wantErr {
					t.Errorf("us.Find() err = %s", err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("us.Find() = %+v, want %+v", got, tc.want)
				}
			})
		}
	}
}

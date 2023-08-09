package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type postgresDb[T Storable] struct {
	db        *sql.DB
	tableName string
}

func (p *postgresDb[T]) Insert(obj T) {
	_, err := p.Get(obj.GetPrimaryKey())
	if err == nil {
		return
	}
	jsonObject, _ := json.Marshal(obj)
	fmt.Printf("data is being saved: %v\n", obj)
	r, err := p.db.Query(fmt.Sprintf("insert into %s values (%d, '%v')", p.tableName, obj.GetPrimaryKey(), string(jsonObject)))
	fmt.Printf("rows returned %v, error %v\n", r, err)
	fmt.Printf("insert into %s values (%d, '%v')", p.tableName, obj.GetPrimaryKey(), string(jsonObject))
}

func (p *postgresDb[T]) Update(obj T) {
	_, err := p.Get(obj.GetPrimaryKey())
	if err != nil {
		return
	}
	jsonObject, _ := json.Marshal(obj)
	fmt.Printf("data is being saved: %v\n", obj)
	r, err := p.db.Query(fmt.Sprintf("update %s SET data='%v' where id=%d", p.tableName, string(jsonObject), obj.GetPrimaryKey()))
	fmt.Printf("rows returned %v, error %v\n", r, err)
	fmt.Printf("update %s SET data='%v' where id=%d", p.tableName, string(jsonObject), obj.GetPrimaryKey())
}

func (p *postgresDb[T]) Delete(key int) T {
	r, err := p.Get(key)
	if err != nil {
		return r
	}
	p.db.Query(fmt.Sprintf("delete from %s where id = %d", p.tableName, key))
	return r
}

func (p *postgresDb[T]) Get(key int) (T, error) {
	var id int
	var d string
	var data T
	r, err := p.db.Query(fmt.Sprintf("select * from %s where id = %d", p.tableName, key))
	if err != nil {
		return data, fmt.Errorf("value not found: %v", err)
	}
	defer r.Close()
	r.Next()
	err = r.Scan(&id, &d)
	if err != nil {
		fmt.Println("error while getting")
		return data, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	if err := json.Unmarshal([]byte(d), &data); err != nil {
		fmt.Println("error while unmarshalling")
		return data, err
	}
	fmt.Printf("%v\n", data)
	return data, nil
}

func CreateDB[T Storable](name string, table string) (DB[T], error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://postgres:postgres@localhost:5432/%s?sslmode=disable", name))
	if err != nil {
		return nil, err
	}
	if _, err := db.Query(fmt.Sprintf("select 1 from %s", table)); err != nil {
		if _, err := db.Query(fmt.Sprintf("create table %s(id int primary key not null, data json);", table)); err != nil {
			log.Fatalf("error while creating table: %v", err)
			return nil, err
		}
	}
	return &postgresDb[T]{db: db, tableName: table}, err
}

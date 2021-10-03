package repository

import (
	"github.com/bernardosecades/feeder/pkg/value"

	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
	"strings"
)

type Sku interface {
	Persist(block map[string]value.Sku) (int64, error)
	Delete(block map[string]value.Sku) (int64, error)
}

type skuPostgreSQL struct {
	SQL *sql.DB
}

// NewSkuPostgreSQL create new instance of repository.Sku with postgresSQL implementation
func NewSkuPostgreSQL(host, port, user, pass, db string) Sku {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db)
	d, err := sql.Open("postgres", dbSource)
	if err != nil {
		panic(err)
	}

	err = d.Ping() // Need to do this to check that the connection is valid
	if err != nil {
		panic(err)
	}


	return &skuPostgreSQL{SQL: d}
}

// Persist save block of value.sku in records table and will ignore the insert if sku already exist
// It will return number of skus inserted
func (r *skuPostgreSQL) Persist(block map[string]value.Sku) (int64, error) {
	if len(block) == 0 {
		return 0, nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	i := 1
	for _, w := range block {
		element := fmt.Sprintf("($%d)", i)
		valueStrings = append(valueStrings, element)
		valueArgs = append(valueArgs, w.String())
		i++
	}

	smt := `INSERT INTO records (sku) VALUES %s ON CONFLICT (sku) DO NOTHING`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	result, err := r.SQL.Exec(smt, valueArgs...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete remove block of value.sku in records table and it will return number of skus deleted
func (r *skuPostgreSQL) Delete(block map[string]value.Sku) (int64, error) {
	if len(block) == 0 {
		return 0, nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	i := 1
	for _, w := range block {
		element := fmt.Sprintf("$%d", i)
		valueStrings = append(valueStrings, element)
		valueArgs = append(valueArgs, w.String())
		i++
	}

	smt := `DELETE FROM records WHERE sku IN (%s)`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	result, err := r.SQL.Exec(smt, valueArgs...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}


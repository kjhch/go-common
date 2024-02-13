package db

import (
	"database/sql/driver"
	"encoding/json"
)

type JsonArray[T any] []T

// sql.Scanner
func (j *JsonArray[T]) Scan(src any) error {
	return json.Unmarshal(src.([]byte), j)
}

// driver.Valuer
func (j JsonArray[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j)
	return raw, err
}

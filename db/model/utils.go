package model

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	Set bool
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)

}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	ns.Set = true

	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil

	return err
}

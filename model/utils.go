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
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte("null"), nil
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	ns.Set = true
	if string(b) == "null" {
		ns.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

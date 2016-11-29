package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		i := sql.NullInt64{}
		if err := json.Unmarshal(data, &i); err != nil {
			return err
		}
		v.NullInt64 = i
		return nil
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		s := sql.NullString{}
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v.NullString = s
		return nil
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

func Quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func QuoteMulti(names []string) []string {
	res := make([]string, 0, len(names))
	for _, name := range names {
		res = append(res, fmt.Sprintf("`%s`", name))
	}
	return res
}

func QuoteString(name string) string {
	return fmt.Sprintf("\"%s\"", name)
}

func Unescape(s string) string {
	return strings.Replace(s, `\`, ``, -1)
}

package sql

import (
	"database/sql"

	jsonUtils "github.com/hdget/utils/json"
	"github.com/spf13/cast"
	"github.com/sqlc-dev/pqtype"
)

func GetNullString(filters map[string]string, key string) sql.NullString {
	if v, ok := filters[key]; ok {
		return sql.NullString{String: v, Valid: true}
	}
	return sql.NullString{}
}

func GetNullInt32(filters map[string]string, key string) sql.NullInt32 {
	if v, ok := filters[key]; ok {
		return sql.NullInt32{Int32: cast.ToInt32(v), Valid: true}
	}
	return sql.NullInt32{}
}

func GetNullInt64(filters map[string]string, key string) sql.NullInt64 {
	if v, ok := filters[key]; ok {
		return sql.NullInt64{Int64: cast.ToInt64(v), Valid: true}
	}
	return sql.NullInt64{}
}

func ToNullString(val string) sql.NullString {
	if val == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: val, Valid: true}
}

func ToNullInt32(val int32) sql.NullInt32 {
	if val == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: val, Valid: true}
}

func ToNullInt64(val int64) sql.NullInt64 {
	if val == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: val, Valid: true}
}

func ToNullJsonObject(val any) pqtype.NullRawMessage {
	if val != nil {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonObject(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

func ToNullJsonArray(val any) pqtype.NullRawMessage {
	if val != nil {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonArray(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

//
//// NullRawMessage represents a json.RawMessage that may be null.
//// NullRawMessage implements the Scanner interface so
//// it can be used as a scan destination, similar to NullString.
//type NullRawMessage struct {
//	RawMessage json.RawMessage
//	Valid      bool // Valid is true if RawMessage is not NULL
//}
//
//// Scan implements the Scanner interface.
//func (n *NullRawMessage) Scan(src interface{}) error {
//	if src == nil {
//		n.Valid = false
//		return nil
//	}
//	switch src := src.(type) {
//	case []byte:
//		srcCopy := make([]byte, len(src))
//		copy(srcCopy, src)
//		n.RawMessage, n.Valid = srcCopy, true
//	default:
//		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, []byte{})
//	}
//	return nil
//}
//
//// Value implements the driver Valuer interface.
//func (n NullRawMessage) Value() (driver.Value, error) {
//	if !n.Valid {
//		return nil, nil
//	}
//	return []byte(n.RawMessage), nil
//}

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

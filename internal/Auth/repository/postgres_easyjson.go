// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package repository

import (
	sql "database/sql"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA814f5e5Decode20192NextLevelInternalAuthRepository(in *jlexer.Lexer, out *PostgresRepo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "DB":
			if in.IsNull() {
				in.Skip()
				out.DB = nil
			} else {
				if out.DB == nil {
					out.DB = new(sql.DB)
				}
				easyjsonA814f5e5DecodeDatabaseSql(in, out.DB)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA814f5e5Encode20192NextLevelInternalAuthRepository(out *jwriter.Writer, in PostgresRepo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"DB\":"
		out.RawString(prefix[1:])
		if in.DB == nil {
			out.RawString("null")
		} else {
			easyjsonA814f5e5EncodeDatabaseSql(out, *in.DB)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostgresRepo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA814f5e5Encode20192NextLevelInternalAuthRepository(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostgresRepo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA814f5e5Encode20192NextLevelInternalAuthRepository(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostgresRepo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA814f5e5Decode20192NextLevelInternalAuthRepository(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostgresRepo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA814f5e5Decode20192NextLevelInternalAuthRepository(l, v)
}
func easyjsonA814f5e5DecodeDatabaseSql(in *jlexer.Lexer, out *sql.DB) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA814f5e5EncodeDatabaseSql(out *jwriter.Writer, in sql.DB) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
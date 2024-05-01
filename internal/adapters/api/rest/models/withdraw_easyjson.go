// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
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

func easyjson4f4a6fc6DecodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(in *jlexer.Lexer, out *Withdraw) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "order":
			out.Order = string(in.String())
		case "sum":
			out.Sum = float64(in.Float64())
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
func easyjson4f4a6fc6EncodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(out *jwriter.Writer, in Withdraw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix[1:])
		out.String(string(in.Order))
	}
	{
		const prefix string = ",\"sum\":"
		out.RawString(prefix)
		out.Float64(float64(in.Sum))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdraw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4f4a6fc6EncodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdraw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4f4a6fc6EncodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdraw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4f4a6fc6DecodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdraw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4f4a6fc6DecodeGithubComK0st1aGophermartInternalAdaptersApiRestModels(l, v)
}

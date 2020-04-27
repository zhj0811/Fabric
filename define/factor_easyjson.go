// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package define

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

func easyjson5155cd71DecodeGithubComZhj0811FabricDefine(in *jlexer.Lexer, out *ResData) {
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
		case "resCode":
			out.ResCode = int(in.Int())
		case "resMsg":
			out.ResMsg = string(in.String())
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
func easyjson5155cd71EncodeGithubComZhj0811FabricDefine(out *jwriter.Writer, in ResData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"resCode\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ResCode))
	}
	{
		const prefix string = ",\"resMsg\":"
		out.RawString(prefix)
		out.String(string(in.ResMsg))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5155cd71EncodeGithubComZhj0811FabricDefine(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5155cd71EncodeGithubComZhj0811FabricDefine(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5155cd71DecodeGithubComZhj0811FabricDefine(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5155cd71DecodeGithubComZhj0811FabricDefine(l, v)
}
func easyjson5155cd71DecodeGithubComZhj0811FabricDefine1(in *jlexer.Lexer, out *Factory) {
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
		case "key":
			out.Key = string(in.String())
		case "value":
			out.Value = string(in.String())
		case "expand1":
			out.Expand1 = string(in.String())
		case "expand2":
			out.Expand2 = string(in.String())
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
func easyjson5155cd71EncodeGithubComZhj0811FabricDefine1(out *jwriter.Writer, in Factory) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key\":"
		out.RawString(prefix[1:])
		out.String(string(in.Key))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.String(string(in.Value))
	}
	if in.Expand1 != "" {
		const prefix string = ",\"expand1\":"
		out.RawString(prefix)
		out.String(string(in.Expand1))
	}
	if in.Expand2 != "" {
		const prefix string = ",\"expand2\":"
		out.RawString(prefix)
		out.String(string(in.Expand2))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Factory) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5155cd71EncodeGithubComZhj0811FabricDefine1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Factory) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5155cd71EncodeGithubComZhj0811FabricDefine1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Factory) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5155cd71DecodeGithubComZhj0811FabricDefine1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Factory) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5155cd71DecodeGithubComZhj0811FabricDefine1(l, v)
}

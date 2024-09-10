// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package responses

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

func easyjson576be6ddDecodeGithubComVilinvilTaskMessaggioPkgResponses(in *jlexer.Lexer, out *ResponseSuccessful) {
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
		case "body":
			out.Body = string(in.String())
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
func easyjson576be6ddEncodeGithubComVilinvilTaskMessaggioPkgResponses(out *jwriter.Writer, in ResponseSuccessful) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix[1:])
		out.String(string(in.Body))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseSuccessful) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComVilinvilTaskMessaggioPkgResponses(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseSuccessful) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComVilinvilTaskMessaggioPkgResponses(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseSuccessful) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComVilinvilTaskMessaggioPkgResponses(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseSuccessful) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComVilinvilTaskMessaggioPkgResponses(l, v)
}

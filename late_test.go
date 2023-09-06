package late

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"golang.org/x/text/language"
)

var (
	_ json.Unmarshaler = (*Language)(nil)
	_ json.Marshaler   = (*Language)(nil)
	_ fmt.Stringer     = (*Language)(nil)
)

func TestLanguageJSON(t *testing.T) {
	type TestT struct {
		Lang Language `json:"lang"`
	}

	buf := []byte(`{"lang":"en-US"}`)

	t.Run("Unmarshal", func(st *testing.T) {
		val := TestT{}
		if err := json.Unmarshal(buf, &val); err != nil {
			st.Fatalf("json unmarhsal: %v", err)
		}

		expected := Language(language.AmericanEnglish)
		if val.Lang != expected {
			st.Errorf("%v != %v", val.Lang, expected)
		}
	})

	t.Run("Marshal", func(st *testing.T) {
		b, err := json.Marshal(TestT{
			Lang: Language(language.AmericanEnglish),
		})
		if err != nil {
			st.Fatalf("json marshal: %v", err)
		}

		if !bytes.Equal(b, buf) {
			st.Errorf("%q != %q", string(b), string(buf))
		}
	})
}

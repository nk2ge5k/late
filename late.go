package late

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/text/language"
)

// Language code
// TODO(nk2ge5k): not sure if it is really necessary to use golang.org/x/text/language
type Language language.Tag

// MarshalJSON returns json represetations of the language.
// Implements json.Marhshaler interface.
func (l Language) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(language.Tag(l).String())), nil
}

// String return string represetation of the language.
func (l Language) String() string {
	return language.Tag(l).String()
}

// UnmarshalJSON parses json into language.
// Implements json.Unarhshaler interface.
func (l *Language) UnmarshalJSON(b []byte) error {
	if len(b) <= 2 {
		return fmt.Errorf("invalid language string representation %q", string(b))
	}

	str := string(b[1 : len(b)-1])

	t, err := language.All.Parse(str)
	if err != nil {
		return fmt.Errorf("parse %q: %w", str, err)
	}

	*l = Language(t)

	return nil
}

// Request contains information needed for translation to be made
type Request struct {
	// List of arguments that should be placed inside of string during
	// translation.
	Args map[string]string `json:"args"`

	// Code of the language into which text is being translated (BCP 47)
	// https://en.wikipedia.org/wiki/IETF_language_tag
	Lang string `json:"lang"`

	// Translation key
	Key string `json:"key"`

	// The name of the keyset where translation key has been stored.
	// TODO(nk2ge5k): i believe this field should be omitted or hidden for
	//	public APIs.
	Keyset string `json:"keyset"`

	// The number of items for plural forms.
	Count int `json:"count"`
}

// Translation represents data as it stored inside of database/cache
type Translation struct {
	// Translation key
	Key string `json:"key"`

	// Any metadata
	Metadata any `json:"metadata"`

	// Translated text for the singular form and for every plural form if necessary.
	// NOTE(nk2ge5k): may be usefull to use array here.
	Value []string `json:"value"`

	// Code of the language into which text is being translated (BCP 47)
	// https://en.wikipedia.org/wiki/IETF_language_tag
	Lang Language `json:"lang"`
}

// LangurageConfiguration contains configuration for every supported language.
type LanguageConfiguration struct {
	// TODO(nk2ge5k): some fileds here
}

// LookupFormIndex tries to find index of the plural form based on the language
// and number of items, such index can be then used to extract translation from
// Translation.Value field
func (lc *LanguageConfiguration) LookupFormIndex(lang Language, n int) (int, error) {
	return 0, errors.New("not implemented yet")
}

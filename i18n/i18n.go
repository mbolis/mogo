package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed *.json
var messages embed.FS

var bundle *i18n.Bundle
var localizer *i18n.Localizer

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := messages.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		_, err := bundle.LoadMessageFileFS(messages, f.Name())
		if err != nil {
			panic(err)
		}
	}

	localizer = i18n.NewLocalizer(bundle)
}

func SetLang(l language.Tag) {
	localizer = i18n.NewLocalizer(bundle, l.String())
}

func T(id string) string {
	s, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	return s
}

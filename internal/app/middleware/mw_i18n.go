package middleware

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var serverLangs = []language.Tag{
	language.SimplifiedChinese, // zh-Hans fallback
	language.AmericanEnglish,   // en-US
	language.Japanese,          // jp
}

func getAcceptLanguage(acceptLanguate string) (lang string) {
	var matcher = language.NewMatcher(serverLangs)
	t, _, _ := language.ParseAcceptLanguage(acceptLanguate)
	_, idx, _ := matcher.Match(t...)

	return serverLangs[idx].String()
}

func createI18nBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, l := range serverLangs {
		messageFile := fmt.Sprintf("pkg/locales/active.%s.toml", strings.ToUpper(l.String()))
		logrus.Printf("LANGUAGE:%s", l.String())
		bundle.MustLoadMessageFile(messageFile)
	}

	return bundle
}

// Localizer Localizer
func Localizer(c *gin.Context) *i18n.Localizer {
	val, ok := c.Get("localizer")

	if !ok {
		return &i18n.Localizer{}
	}

	return val.(*i18n.Localizer)
}

// I18nMiddleware i18n middleware
func I18nMiddleware() gin.HandlerFunc {
	// NOTE: Create a go-i18n Bundle to use for the lifetime of your application.
	bundle := createI18nBundle()

	return func(c *gin.Context) {
		locale := c.Query("locale")
		if locale != "" {
			c.Request.Header.Set("Accept-Language", locale)
		}
		lang := getAcceptLanguage(c.GetHeader("Accept-Language"))

		// NOTE: On June 2012, the deprecation of recommendation to use the "X-" prefix has become official as RFC 6648.
		// https://stackoverflow.com/questions/3561381/custom-http-headers-naming-conventions
		// c.Request.Header.Set("I18n-Language", lang)
		c.Set("i18n", lang)

		// NOTE: Create a go-i18n Localizer to use for a set of language preferences.
		localizer := i18n.NewLocalizer(bundle, lang, c.GetHeader("Accept-Language"))
		c.Set("localizer", localizer)

		c.Next()
	}
}

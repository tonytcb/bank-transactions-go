package handler

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	englishTranslator := en.New()
	uni = ut.New(englishTranslator, englishTranslator)

	validate = validator.New()

	t, _ := uni.GetTranslator("en")
	trans = t
	entranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func translateValidations(e validator.ValidationErrors) map[string]string {
	var errs = make(map[string]string, 0)

	formatFieldName := func(v string) string {
		v = strings.ToLower(v)
		parts := strings.Split(v, ".")
		return strings.Join(parts[1:], ".")
	}

	for _, v := range e {
		fieldName := formatFieldName(v.Namespace())
		value := v.Translate(trans)

		errs[fieldName] = value
	}

	return errs
}

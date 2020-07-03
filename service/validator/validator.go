package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"log"
)

var (
	uni   *ut.UniversalTranslator
	Trans ut.Translator
)

func New() *validator.Validate {
	var found bool

	v := validator.New()

	translator := en.New()
	uni = ut.New(translator, translator)

	Trans, found = uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	if err := entranslations.RegisterDefaultTranslations(v, Trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", Trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", Trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("passwrd", Trans, func(ut ut.Translator) error {
		return ut.Add("passwrd", "{0} is not strong enough", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwrd", fe.Field())
		return t
	})

	_ = v.RegisterValidation("passwrd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	return v
}

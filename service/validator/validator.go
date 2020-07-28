package validator

import (
	"log"

	"github.com/go-playground/locales/en"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var Trans ut.Translator

func New() *validator.Validate {
	var found bool

	english := en.New()
	uni := ut.New(english, english)

	if Trans, found = uni.GetTranslator("en"); !found {
		log.Fatal("en not found")
	}

	validate := validator.New()

	if err := entranslations.RegisterDefaultTranslations(validate, Trans); err != nil {
		log.Fatal(err)
	}

	return validate
}

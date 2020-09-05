package validatorx

import (
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enLang "github.com/go-playground/validator/v10/translations/en"
	zhLang "github.com/go-playground/validator/v10/translations/zh"
	"github.com/storezhang/gox"
)

var (
	translator *ut.UniversalTranslator
	validate   *validator.Validate
)

type Validator struct {
	validator *validator.Validate
}

func init() {
	validate = validator.New()
	if err := initValidation(); nil != err {
		panic(err)
	}
}

func (cv *Validator) Validate(i interface{}) (err error) {
	return cv.validator.Struct(i)
}

func New() (validator *Validator) {
	translator = ut.New(en.New(), en.New(), zh.New())
	if english, success := translator.GetTranslator("en"); success {
		if err := enLang.RegisterDefaultTranslations(validate, english); nil != err {
			return
		}
	}
	if chinese, success := translator.GetTranslator("zh"); success {
		if err := zhLang.RegisterDefaultTranslations(validate, chinese); nil != err {
			return
		}
	}

	validator = &Validator{validator: validate}

	return
}

func I18n(lang string, errs validator.ValidationErrors) (i18n validator.ValidationErrorsTranslations) {
	sep := "_"
	if strings.Contains(lang, "-") {
		sep = "-"
	}

	splits := strings.Split(lang, sep)
	for i := 0; i < len(splits); i++ {
		if t, s := translator.FindTranslator(lang); s {
			i18n = errs.Translate(t)

			break
		} else {
			if index := strings.LastIndex(lang, sep); -1 == index {
				break
			} else {
				lang = lang[0:index]
			}
		}
	}
	if nil == i18n {
		if t, s := translator.GetTranslator("zh"); s {
			i18n = errs.Translate(t)
		}
	}

	// 得到的国际化字符串是一个带请求体的键值，类似于LoginReq.Password：错误消息
	// 而我们需要的是password: 错误消息
	for field, msg := range i18n {
		newField := gox.CamelName(string(field[strings.IndexRune(field, '.')]))
		i18n[newField] = msg
		// 删除原来的错误消息，避免前端混乱
		delete(i18n, field)
	}

	return
}

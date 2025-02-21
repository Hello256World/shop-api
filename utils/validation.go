package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	regex := `^(\+?\d{1,3}[-. ]?)?(\(?\d{3}\)?[-. ]?)?\d{3}[-. ]?\d{4}$`
	matched, _ := regexp.MatchString(regex, phone)
	return matched
}

func Validation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phone", phoneValidation)
	}
}

func FormValidation(err string, fileds map[string]string) []string {
	errorsText := map[string]string{
		"gte":      "فیلد %v باید بزرگ تر یا برابر صفر باشد",
		"required": "فیلد %v اجباری می باشد",
		"email":    "ایمیل وارد شده معتبر نمی باشد",
		"phone":    "تلفن همراه معتبر نمی باشد",
		"eqfield":  "تکرار رمز عبور باید با رمز عبور مطابقت داشته باشد",
		"gt":       "%v باید بزرگ تر از صفر باشد",
	}

	var final []string
	myError := strings.Split(err, "\n")

	for _, newError := range myError {
		for validate, message := range errorsText {
			if strings.Contains(newError, validate) {
				for filed, filedName := range fileds {
					if strings.Contains(newError, filed) {
						if validate == "email" || validate == "phone" || validate == "eqfield" {
							final = append(final, message)
						} else {
							final = append(final, fmt.Sprintf(message, filedName))
						}
						break
					}
				}
				break
			}
		}
	}

	return final
}

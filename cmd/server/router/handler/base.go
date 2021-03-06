package handler

import (
	"github.com/aaronchen2k/tester/internal/pkg/utils"
	"github.com/aaronchen2k/tester/internal/server/biz/validate"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

type BaseCtrl struct {
	Ctx iris.Context
}

func (c *BaseCtrl) Validate(s interface{}, ctx iris.Context) bool {
	err := validate.Validate.Struct(s)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validate.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(_utils.ApiRes(400, e, nil))
				return true
			}
		}
	}

	return false
}

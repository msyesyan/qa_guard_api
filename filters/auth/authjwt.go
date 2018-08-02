package filters

import (
	"regexp"
	"qa_guard_api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/astaxie/beego/context"
)

func renderError(ctx *context.Context) {
	ctx.Output.SetStatus(401)
	var err interface{}
	err = "Not authorized"
	ctx.Output.JSON(err, false, false)
}

// AuthJwtToken check jwt
func AuthJwtToken(ctx *context.Context) {
	//TODO should be more strict
	if shouldSkip, _ := regexp.MatchString("/sign_[in|up]", ctx.Input.URL()); shouldSkip {
		return
	}

	// Extract jwt token from request header
	tokenStr, err := request.HeaderExtractor{"Authorization"}.ExtractToken(ctx.Request)

	if err != nil {
		// TODO ctx.Abort will not trigger ErrorHandler
		// ctx.Abort(401, "Not authorized")

		renderError(ctx)
	}

	var keyFunc = func(t *jwt.Token) (interface{}, error) {
		return []byte("qa_guard_api"), nil
	}

	token, err := (&jwt.Parser{UseJSONNumber: true}).ParseWithClaims(tokenStr, &jwt.StandardClaims{}, keyFunc)

	if (err != nil) {
		renderError(ctx)
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		_, err := models.GetUserByUsernameOrEmail("", claims.Subject)
		if err != nil {
			renderError(ctx)
		}
	} else {
		renderError(ctx)
	}
}

package controllers

import (
	"time"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go"
	"encoding/json"
	"errors"
	"qa_guard_api/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//  UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.AddUser(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title signup
// @Description signup user
// @Param body body models.UserSignup true 	"body for user content"
// @Success 200 {object} models.User
// @Failure 403 user exists
// @router /sign_up [post]
func (u *UserController) SignUp() {
	var user models.UserSignup
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid, err := models.SignUp(&user)

	if err != nil {
		u.Ctx.Output.SetStatus(400)
		u.Data["json"] = err.Error()
	} else {
	  result, _ := models.GetUserById(uid)
		u.Data["json"] = result
	}

	u.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetCurrentUser
// @Title Get current login in user
// @Description Get current login in user
// @Param	Authorization	header	string true	"jwt token"
// @Sucess 200 {object} models.User
// @Failure 401
// @router /current
func (c *UserController) GetCurrentUser() {
	token, _ := request.HeaderExtractor{"Authorization"}.ExtractToken(c.Ctx.Request)
	user, err := models.GetUserFromToken(token)

	if err != nil {
		c.Abort("401")
	}

	c.Data["json"] = user
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.User{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateUserById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteUser(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// SignIn
// @Title SignIn
// @Description sign in with email and password
// @Param	body	body	models.UserSignIn	true "body for user sign in"
// @Sucess 200 json{ jwt }
// @Failure 401 authenticate error
// @router /sign_in [post]
func (c *UserController) SignIn() {
	var userSignIn models.UserSignIn
	json.Unmarshal(c.Ctx.Input.RequestBody, &userSignIn)

	// TODO get user by email
	user, err := models.GetUserByUsernameOrEmail("", userSignIn.Email)

	if err != nil {
		c.Data["json"] = "user authenticate fail"
		c.Ctx.Output.SetStatus(401)
	}

	// TODO validate password

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Subject: user.Email,
		ExpiresAt: time.Now().Add(time.Hour  * 24 * 7).Unix(),
	})

	ss, jwtErr := token.SignedString([]byte("qa_guard_api"))

	if jwtErr != nil {
		c.Data["json"] = jwtErr.Error()
	} else {
		c.Data["json"] = map[string]string{"jwt": ss}
	}

	c.ServeJSON()
}

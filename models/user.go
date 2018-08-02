package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int64  `orm:"auto"`
	Username string `orm:"size(128)"`
	Email    string `orm:"size(128)"`
	Projects []*Project `orm:"reverse(many)"`
	PrjectUsers []*ProjectUser `orm:"reverse(many)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime);"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime);"`
}

func (u *User) TableName() string {
    return "users"
}

// UserSignup with username and email
type UserSignup struct {
	Username string
	Email    string
}

// UserSignIn with email and password
type UserSignIn struct {
	Email string
	Password string
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// SignUp signup user
func SignUp(user *UserSignup) (int64, error) {
	u, _ := GetUserByUsernameOrEmail(user.Username, user.Email)

	if (u.Id > 0) {
		return u.Id, errors.New("username or email is already exist")
	}

	return AddUser(&User{
		Email: user.Email,
		Username: user.Username,
	})
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.QueryTable(new(User)).Filter("Id", id).RelatedSel().One(v); err == nil {
		o.LoadRelated(v, "Projects")
		return v, nil
	}
	return nil, err
}

// GetUserByUsernameOrEmail find user by username or email
func GetUserByUsernameOrEmail(username string, email string) (user User, err error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond.And("username", username).Or("email", email)

	err = o.QueryTable(new(User)).SetCond(cond).One(&user)

	if err == nil {
		o.LoadRelated(&user, "Projects")
	}

	return user, err
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []User
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// GetUserFromToken get user from jwt token
func GetUserFromToken(tokenStr string) (user User, err error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	var keyFunc = func(t *jwt.Token) (interface{}, error) {
		return []byte("qa_guard_api"), nil
	}

	token, err := (&jwt.Parser{UseJSONNumber: true}).ParseWithClaims(tokenStr, &jwt.StandardClaims{}, keyFunc)

	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		user, err := GetUserByUsernameOrEmail("", claims.Subject)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	return user, errors.New("Not found")
}

// AddUserProject add a project to user
func AddUserProject(user *User, project *Project) (result *Project, err error){
	// TODO transaction
	project.User = user

	id, err := AddProject(project)

	if err != nil {
		return project, err
	}

	result, err = GetProjectById(id)

	if err != nil {
		return result, err
	}

	projectUser := &ProjectUser {
		User: user,
		Project: result,
	}

	_, err = AddProjectUser(projectUser)

	if err != nil {
		return result, err
	}

	o := orm.NewOrm()
	o.LoadRelated(result, "ProjectUsers")

	return result, err
}
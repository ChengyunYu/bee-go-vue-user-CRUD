package controllers

import (
	"bee-go-vue/conf"
	"bee-go-vue/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/ldap.v3"
	"log"
	"strconv"
)

var userResults = models.UserResultCollection{}
var userResultsEmpty = models.UserResultCollection{}
var totalResults int = 0
var totalPage int = 1
var status = conf.FAILURE

type H map[string]interface{}

func init() {
	emptyPage := models.UserCollection{}
	userResultsEmpty.UserResult = append(userResultsEmpty.UserResult, emptyPage)
	userResults = userResultsEmpty
}

type UserController struct {
	beego.Controller
}

func (u *UserController) URLMapping() {
	u.Mapping("Get", u.GetUsers)
	u.Mapping("Post", u.PostUser)
	u.Mapping("Put", u.PutUser)
	u.Mapping("Delete", u.DeleteUser)
	u.Mapping("Get", u.ChangeUserPassword)
}

func (u *UserController) ShowIndex() {
	u.TplName = "index.html"
}

func (u *UserController) ChangeUserPassword() {
	password := u.bindPassword()
	uid := u.Ctx.Input.Param(":uid")
	org := u.Ctx.Input.Param(":org")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.HOST_NAME, conf.PORT_NUMBER))
	if err != nil {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": err.Error(),
		}
		u.ServeJSON()

		l.Close()
		return
	}
	defer l.Close()

	res, message := models.LogUserIn(conf.ADMIN_UID, conf.ADMIN_ORG, conf.ADMIN_PASSWORD, l)
	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	res, message = models.ModifyUserPassword(uid, org, password.Password1, l)
	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	u.Data["json"] = H{
		"result": res,
		"message": message,
	}
	u.ServeJSON()
}

func GetKeywordUsers(keyword string) string {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.HOST_NAME, conf.PORT_NUMBER))
	if err != nil {
		log.Println(err)
		userResults = userResultsEmpty
		status = conf.FAILURE
		return err.Error()
	}
	defer l.Close()

	res, message, num, pages, results := models.GetUsers(keyword, l)
	status = res

	if status != conf.SUCCESS {
		status = conf.FAILURE
		userResults = userResultsEmpty
		log.Println(message)
		return message
	}

	if num == 0 {
		results = userResultsEmpty
		pages = 1
	}

	status = conf.SUCCESS
	totalResults = num
	totalPage = pages
	userResults = results
	return message
}

func (u *UserController) GetUsers() {
	keyword := u.Ctx.Input.Param(":keyword")
	currentPage, _ := strconv.Atoi(u.Ctx.Input.Param(":page"))
	message := GetKeywordUsers(keyword)

	if currentPage > totalPage {
		currentPage = totalPage
	}

	u.Data["json"] = H{
		"result": status,
		"message": message,
		"results": userResults.UserResult[currentPage - 1],
		"total_page": totalPage,
		"total_results": totalResults,
		"current_page":currentPage,
	}

	u.ServeJSON()
	return
}

func (u *UserController) GetAllUsers() {
	currentPage, _ := strconv.Atoi(u.Ctx.Input.Param(":page"))
	message := GetKeywordUsers("")

	if currentPage > totalPage {
		currentPage = totalPage
	}

	u.Data["json"] = H{
		"result": status,
		"message": message,
		"results": userResults.UserResult[currentPage - 1],
		"total_page": totalPage,
		"total_results": totalResults,
		"current_page":currentPage,
	}

	u.ServeJSON()
	return
}

func (u *UserController) PostUser() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.HOST_NAME, conf.PORT_NUMBER))
	if err != nil {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": err.Error(),
		}
		u.ServeJSON()

		l.Close()
		return
	}
	defer l.Close()

	res, message := models.LogUserIn(conf.ADMIN_UID, conf.ADMIN_ORG, conf.ADMIN_PASSWORD, l)

	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	user := u.bind()

	res, message = models.PostUser(user, l)

	u.Data["json"] = H{
		"result": res,
		"message": message,
	}
	u.ServeJSON()
}

func (u *UserController) PutUser() {
	user := u.bind()
	uid := u.Ctx.Input.Param(":uid")
	org := u.Ctx.Input.Param(":org")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.HOST_NAME, conf.PORT_NUMBER))
	if err != nil {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": err.Error(),
		}
		u.ServeJSON()

		l.Close()
		return
	}
	defer l.Close()

	res, message := models.LogUserIn(conf.ADMIN_UID, conf.ADMIN_ORG, conf.ADMIN_PASSWORD, l)
	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	res, message = models.ModifyUidOrg(uid, org, user.Mail, user.Org, l)
	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	res, message = models.PutUser(user, l)
	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		res, message = models.ModifyUidOrg(user.Mail, user.Org, uid, org, l)
		return
	}

	u.Data["json"] = H{
		"result": res,
		"message": message,
	}
	u.ServeJSON()
}

func (u *UserController) DeleteUser() {
	uid := u.Ctx.Input.Param(":uid")
	org := u.Ctx.Input.Param(":org")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.HOST_NAME, conf.PORT_NUMBER))
	if err != nil {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": err.Error(),
		}
		u.ServeJSON()

		l.Close()
		return
	}
	defer l.Close()

	res, message := models.LogUserIn(conf.ADMIN_UID, conf.ADMIN_ORG, conf.ADMIN_PASSWORD, l)

	if res != conf.SUCCESS {
		log.Println(err)
		u.Data["json"] = H{
			"result": conf.FAILURE,
			"message": message,
		}
		u.ServeJSON()
		return
	}

	res, message = models.DeleteUser(uid, org, l)

	u.Data["json"] = H{
		"result": res,
		"message": message,
	}
	u.ServeJSON()
}

func (u *UserController) bind() (us models.User) {
	err := json.NewDecoder(u.Ctx.Request.Body).Decode(&us)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func (u *UserController) bindPassword() (ps models.Password) {
	err := json.NewDecoder(u.Ctx.Request.Body).Decode(&ps)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

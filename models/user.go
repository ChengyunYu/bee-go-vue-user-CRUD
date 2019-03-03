package models

import (
	"bee-go-vue/conf"
	"fmt"
	"github.com/astaxie/beego/orm"
	"gopkg.in/ldap.v3"
	"log"
	"math"
	"sort"
)

type User struct {
	Surname		string	`orm:"size(255)" json:"surname"`
	GivenName	string	`orm:"size(255)" json:"given_name"`
	Password	string	`orm:"size(255)" json:"password_1"`
	Password2	string	`orm:"size(255)" json:"password_2"`
	Mail		string	`orm:"pk" orm:"size(255)" json:"email"`
	Org			string	`orm:"size(255)" json:"org"`
}

type Password struct {
	Password1	string	`orm:"size(255)" json:"password_1"`
	Password2	string	`orm:"size(255)" json:"password_2"`
}

type UserList []User

type UserCollection struct {
	Users []User `json:"users"`
}

type UserResultCollection struct {
	UserResult []UserCollection
}

func (u UserList) Len() int {
	return len(u)
}

func (u UserList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UserList) Less(i, j int) bool {
	return u[i].Surname < u[j].Surname
}

func init() {
	orm.RegisterModel(new(User))
}

func GetUsers(keyword string, l *ldap.Conn) (string, string, int, int, UserResultCollection) {
	keywordPart := ""
	if keyword != "" {
		keywordPart = keyword + "*"
	}

	filter := "(&(|(uid=*" + keywordPart + ")" +
		"(sn=*" + keywordPart + ")" +
		"(givenName=*" + keywordPart + ")" +
		"(ou=*" + keywordPart + "))" + "(objectClass=inetOrgPerson)" +
		"(!(uid=cheyu)))"

	searchRequest := ldap.NewSearchRequest(
		conf.BASE_DN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, []string{"uid", "sn", "givenName", "mail", "ou"},
		nil)

	result := UserResultCollection{}
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return conf.FAILURE, err.Error(), 0, 1, result
	}

	page := UserCollection{}
	counter := 0
	allPages := UserList{}

	for _, entry := range sr.Entries {
		org := entry.GetAttributeValue("ou")
		data := User{
			Surname:   entry.GetAttributeValue("sn"),
			GivenName: entry.GetAttributeValue("givenName"),
			Mail:      entry.GetAttributeValue("mail"),
			Org:       org,
		}
		allPages = append(allPages, data)
	}

	sort.Sort(allPages)
	totalResults := len(allPages)
	totalPages := int(math.Ceil(float64(totalResults) / conf.PAGE_SIZE))

	for _, user := range allPages {
		org := user.Org
		data := User{
			Surname: user.Surname,
			GivenName: user.GivenName,
			Mail: user.Mail,
			Org: org,
		}
		page.Users = append(page.Users, data)
		counter++

		if counter >= conf.PAGE_SIZE {
			result.UserResult = append(result.UserResult, page)
			page = UserCollection{}
			counter = 0
		}
	}

	if counter > 0 {
		result.UserResult = append(result.UserResult, page)
	}

	return conf.SUCCESS, keyword, totalResults, totalPages, result
}

func PutUser(user User, l *ldap.Conn) (string, string) {
	userDN := GetUserDN(user.Mail, user.Org, l)

	modifyRequest := ldap.NewModifyRequest(
		userDN,
		nil)
	modifyRequest.Replace("uid", []string{user.Mail})
	modifyRequest.Replace("cn", []string{user.GivenName + " " + user.Surname})
	modifyRequest.Replace("givenName", []string{user.GivenName})
	modifyRequest.Replace("sn", []string{user.Surname})
	modifyRequest.Replace("displayName", []string{user.GivenName + user.Surname})
	modifyRequest.Replace("mail", []string{user.Mail})
	modifyRequest.Replace("objectClass", []string{"top", "inetOrgPerson"})
	modifyRequest.Replace("ou", []string{user.Org})

	err := l.Modify(modifyRequest)
	if err != nil {
		println(err.Error())
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, user.Mail
}

func PostUser(user User, l *ldap.Conn) (string, string) {
	userDN := "uid="  + user.Mail + ",ou=" + user.Org + ",dc=ibm,dc=com"

	addRequest := ldap.NewAddRequest(
		userDN,
		nil)
	addRequest.Attribute("uid", []string{user.Mail})
	addRequest.Attribute("cn", []string{user.GivenName + " " + user.Surname})
	addRequest.Attribute("givenName", []string{user.GivenName})
	addRequest.Attribute("sn", []string{user.Surname})
	addRequest.Attribute("displayName", []string{user.GivenName + user.Surname})
	addRequest.Attribute("mail", []string{user.Mail})
	addRequest.Attribute("objectClass", []string{"top", "inetOrgPerson"})
	addRequest.Attribute("ou", []string{user.Org})
	addRequest.Attribute("userPassword", []string{user.Password})

	err := l.Add(addRequest)
	if err != nil {
		log.Println(err)
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, user.Mail
}

func DeleteUser(uid string, org string, l *ldap.Conn) (string, string) {
	userDN := GetUserDN(uid, org, l)

	delRequest := ldap.NewDelRequest(userDN, nil)
	err := l.Del(delRequest)

	if err != nil {
		log.Println(err)
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, uid
}

func ModifyUserPassword(uid string, org string, newPassword string, l *ldap.Conn) (string, string) {
	userDN := GetUserDN(uid, org, l)

	passwordModifyRequest := ldap.NewPasswordModifyRequest(userDN, "", newPassword)
	_, err := l.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Println("Password could not be changed: " + err.Error())
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, uid
}

func ModifyUidOrg(uid string, org string, newUid string, newOrg string, l *ldap.Conn) (string, string) {
	userDN := GetUserDN(uid, org, l)
	newRDN := "uid=" + newUid
	newSup := "ou=" + newOrg + ",dc=ibm,dc=com"

	modifyDNRequest := ldap.NewModifyDNRequest(userDN, newRDN, true, newSup)
	err := l.ModifyDN(modifyDNRequest)

	if err != nil {
		errMessage := "DN could not be changed: " + err.Error()
		log.Println(errMessage)
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, uid
}

func LogUserIn(uid string, org string, password string, l *ldap.Conn) (string, string) {
	userDN := GetUserDN(uid, org, l)

	err := l.Bind(userDN, password)
	if err != nil {
		log.Println(err)
		return conf.FAILURE, err.Error()
	}
	return conf.SUCCESS, ""
}

func GetUserDN(uid string, org string, l *ldap.Conn) string {
	searchRequest := ldap.NewSearchRequest(
		"ou=" + org + "," + conf.BASE_DN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", uid),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return conf.FAILURE
	}

	userDN := conf.MISSING

	if len(sr.Entries) != 1 {
		log.Println("User does not exist or too many entries returned")
	} else {
		userDN = sr.Entries[0].DN
	}

	return userDN
}

func AddUserToGroup(uid string, org string, groupDN string, l *ldap.Conn) string {
	userDN := GetUserDN(uid, org, l)

	modifyRequest := ldap.NewModifyRequest(
		groupDN,
		nil)

	modifyRequest.Add("member", []string{userDN})

	err := l.Modify(modifyRequest)
	if err != nil {
		log.Println(err)
		return conf.FAILURE
	}

	return conf.SUCCESS
}

func CheckUserInGroup(uid string, org string, groupDN string, l *ldap.Conn) bool {
	userDN := GetUserDN(uid, org, l)

	searchRequest := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectclass=groupOfNames)(member=%s))", userDN),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	if len(sr.Entries) != 1 {
		return false
	}

	return true
}

func DeleteUserFromGroup(uid string, org string, groupDN string, l *ldap.Conn) string {
	userDN := GetUserDN(uid, org, l)

	modifyRequest := ldap.NewModifyRequest(
		groupDN,
		nil)

	modifyRequest.Delete("member", []string{userDN})

	err := l.Modify(modifyRequest)
	if err != nil {
		log.Println(err)
		return conf.FAILURE
	}

	return conf.SUCCESS
}

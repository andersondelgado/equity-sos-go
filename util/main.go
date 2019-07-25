package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"../config"
	"../model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/timjacobi/go-couchdb"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseJSON struct {
	Data interface{} `json:"data"`
}

type RolesDev struct {
	IDs    string   `json:"_id"`
	Rev    string   `json:"rev"`
	Rol    string   `json:"rol"`
	Module []Module `json:"module"`
}

type Roles struct {
	Rol    string   `json:"rol"`
	Module []Module `json:"module"`
}

type Rol struct {
	Rol    string `json:"rol"`
	Module Module `json:"module"`
	Acl    string `json:"acl"`
}

type Module struct {
	Url          string   `json:"url"`
	IconWeb      string   `json:"iconWeb"`
	Mobile       string   `json:"mobile"`
	Android      string   `json:"android"`
	Ios          string   `json:"ios"`
	Ionic        string   `json:"ionic"`
	Name         string   `json:"name"`
	LangProperty string   `json:"lang_property"`
	Acl          string   `json:"acl"`
	Visible      bool     `json:"visible"`
	Value        []string `json:"value"`
}

type User struct {
	// ID       string `json:"id"`
	// Avatar   string `json:"avatar"`
	// Username string `json:"username"`
	// Email    string `json:"email"`
	IDs      string `json:"_id"`
	Rev      string `json:"_rev"`
	ID       string `json:"id"`
	Avatar   string `json:"avatar"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type ResponseU struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    model.User `json:"data"`
}

type ResponseProfile struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    model.Profile `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type TokenDecode struct {
	ID    string
	Email string
	jwt.StandardClaims
}

type CustomClaims struct {
	// ID    uint
	ID    string
	Email string
	jwt.StandardClaims
}

type Paginate struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// type CustomValidator struct {
// 	Validator *validator.Validate
// }

func SetToken(claims CustomClaims) string {
	//claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokensJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	app_key := config.EnviromentsRaw().AppKey
	key := []byte(app_key)

	tokenStr, er := tokensJWT.SignedString(key)
	if er != nil {
		fmt.Println("error: ", er)
	}

	return tokenStr
}

func DecodeHeaderToken(c *gin.Context) ResponseU {
	var datas ResponseU
	r := c.Request.Header.Get("Authorization")
	if r == "" {
		datas = ResponseU{
			false,
			"Empty token",
			model.User{},
		}
	} else {
		headers := c.Request.Header
		token := headers.Get("Authorization")

		splitToken := strings.Split(token, " ")

		t := DecodeToken(splitToken[1])
		if t.Success == false {
			datas = ResponseU{
				false,
				"token invalid",
				model.User{},
			}
		} else {
			datas = ResponseU{
				true,
				"ok",
				t.Data,
			}
		}
	}

	return datas
}

func DecodeToken(strToken string) ResponseU {

	var datas ResponseU

	if strings.HasPrefix(strToken, "") {
		datas = ResponseU{
			false,
			"token invalid",
			model.User{},
		}
	}
	app_key := config.EnviromentsRaw().AppKey
	tokenObj, er := jwt.ParseWithClaims(strToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(app_key), nil
	})

	var (
		find model.User
	)

	if claims, ok := tokenObj.Claims.(*CustomClaims); ok && tokenObj.Valid {
		//fmt.Printf("%vccc %v", claims.Email)
		// db.Debug().Where(map[string]interface{}{"email": claims.Email}).Find(&find)
		var arrKey = []string{"users", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta":        arrKey[0],
				"users.email": claims.Email,
			},
			Fields: arrKey,
		}

		respText := FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.UserDocumentsArray
		json.Unmarshal(decode, &results)

		var us []model.User
		for i := range results.Doc {
			// fmt.Println("##i: ", results.Doc[i])
			us = append(us, model.User{
				IDs:      results.Doc[i].ID,
				Rev:      results.Doc[i].Rev,
				ID:       results.Doc[i].ID,
				Avatar:   results.Doc[i].User.Avatar,
				Username: results.Doc[i].User.Username,
				Email:    results.Doc[i].User.Email,
				Password: results.Doc[i].User.Password,
			})
		}

		if len(us) == 0 {
			datas = ResponseU{
				false,
				"error_exception",
				model.User{},
			}
			// c.JSON(200, datas)
		} else {
			find = us[0]
			datas = ResponseU{
				true,
				"ok",
				find,
			}
		}

	} else {
		fmt.Printf("#####er", er)
		datas = ResponseU{
			false,
			"token invalid",
			model.User{},
		}
	}

	return datas
}

func Auth(c *gin.Context) User {
	user := c.MustGet("user").(string)
	// fmt.Println("###: ", user)
	str := []byte(user)
	var raw User

	json.Unmarshal(str, &raw)

	return raw
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func IsRead(c *gin.Context, rol Rol) Response {
	var response Response
	var arrayPerm []string
	if rol.Rol != "" {
		arrayPerm = ArrayPermWithRol(c, rol)
	} else {
		arrayPerm = ArrayPerm(c, rol)
	}
	perm := CanPerm()
	index := IndexOf(perm[0], arrayPerm)

	if index != -1 {
		response = Response{
			true,
			"ok",
			"",
		}
	} else {
		response = Response{
			false,
			"error permission " + perm[0],
			"",
		}
	}

	return response
}

func IsCreate(c *gin.Context, rol Rol) Response {

	var response Response
	var arrayPerm []string
	if rol.Rol != "" {
		arrayPerm = ArrayPermWithRol(c, rol)
	} else {
		arrayPerm = ArrayPerm(c, rol)
	}

	perm := CanPerm()
	index := IndexOf(perm[1], arrayPerm)

	if index != -1 {
		response = Response{
			true,
			"ok",
			"",
		}
	} else {
		response = Response{
			false,
			"error permission " + perm[1],
			"",
		}
	}

	return response
}

func IsEdit(c *gin.Context, rol Rol) Response {

	var response Response
	var arrayPerm []string
	if rol.Rol != "" {
		arrayPerm = ArrayPermWithRol(c, rol)
	} else {
		arrayPerm = ArrayPerm(c, rol)
	}

	perm := CanPerm()
	index := IndexOf(perm[2], arrayPerm)
	if index != -1 {
		response = Response{
			true,
			"ok",
			"",
		}
	} else {
		response = Response{
			false,
			"error permission " + perm[2],
			"",
		}
	}

	return response
}

func IsUpdate(c *gin.Context, rol Rol) Response {

	var response Response
	var arrayPerm []string
	if rol.Rol != "" {
		arrayPerm = ArrayPermWithRol(c, rol)
	} else {
		arrayPerm = ArrayPerm(c, rol)
	}

	perm := CanPerm()
	index := IndexOf(perm[3], arrayPerm)

	if index != -1 {
		response = Response{
			true,
			"ok",
			"",
		}
	} else {
		response = Response{
			false,
			"error permission " + perm[3],
			"",
		}
	}

	return response
}

func IsDelete(c *gin.Context, rol Rol) Response {

	var response Response
	var arrayPerm []string
	if rol.Rol != "" {
		arrayPerm = ArrayPermWithRol(c, rol)
	} else {
		arrayPerm = ArrayPerm(c, rol)
	}

	perm := CanPerm()
	index := IndexOf(perm[4], arrayPerm)

	if index != -1 {
		response = Response{
			true,
			"ok",
			"",
		}
	} else {
		response = Response{
			false,
			"error permission " + perm[4],
			"",
		}
	}

	return response
}

func CanPerm() []string {
	perm := []string{"read", "create", "edit", "update", "delete"}

	return perm
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func ArrayPermWithRol(c *gin.Context, rol Rol) []string {

	var arrayPerm []string
	user := c.MustGet("user").(string)
	// fmt.Println("##user: ", user)
	str := []byte(user)
	var raw User
	json.Unmarshal(str, &raw)

	// cloudantUrl := config.StrNoSQLDrive()
	// cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	// if err != nil {
	// 	log.Println("Can not connect to Cloudant database")
	// }

	//ensure db exists
	//if the db exists the db will be returned anyway
	// dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)
	var datas Response

	// if cloudantUrl == "" {
	// 	c.JSON(200, gin.H{})
	// 	//return
	// }

	var arrKey = []string{"profiles", "_id", "_rev"}

	query := model.QuerySelectorAll{
		Selector: map[string]interface{}{
			"meta": arrKey[0],
		},
		Fields: arrKey,
	}

	respText := FindDataAll(query)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.ProfileDocumentsArray
	json.Unmarshal(decode, &results)

	var pr []model.Permission

	for i := range results.Doc {
		if results.Doc[i].Profile.UserID == raw.ID {
			pr = append(pr, model.Permission{
				ID: results.Doc[i].Profile.PermissionID,
			})
		}
	}

	if len(pr) == 0 {
		datas = Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	} else {
		// find := pr[0]

		// fmt.Println("##find: ", (find.ID))

		// var result2 map[string]interface{}

		// errs2 := cloudant.DB(dbName).Get(find.ID, &result2, couchdb.Options{"include_docs": true})
		// if errs2 != nil {
		// 	fmt.Println("##errs2: ", errs2)
		// }

		// jsonToString2, _ := json.Marshal(result2)
		// // fmt.Println("##jsonToString: ", string(jsonToString2))
		// decode2 := []byte(jsonToString2)

		// var resultPerm model.PermissionDoc
		// json.Unmarshal(decode2, &resultPerm)
		// var perms model.Permission

		find := pr[0]

		var arrKey1 = []string{"permissions", "_id", "_rev"}
		query1 := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey1[0],
				"_id":  find.ID,
			},
			Fields: arrKey1,
		}

		respText1 := FindDataAll(query1)

		jsonToString1 := (respText1)
		decode1 := []byte(jsonToString1)
		var results1 model.PermissionDocumentsArray
		json.Unmarshal(decode1, &results1)

		resultPerm := results1.Doc[0]

		var perms model.Permission

		if resultPerm.ID == find.ID {

			// fmt.Println("##resultPerm.Permission.ObjectModulePermission: ",
			// 	(resultPerm.Permission.ObjectModulePermission))

			decodeModules := []byte(resultPerm.Permission.ObjectModulePermission)
			var modules interface{}
			json.Unmarshal(decodeModules, &modules)

			perms = model.Permission{
				Modules: modules,
			}
		}

		jsonToString3, _ := json.Marshal(perms.Modules)
		// fmt.Println("##jsonToString3: ", string(jsonToString3))
		decode3 := []byte(jsonToString3)

		var resultx Roles
		json.Unmarshal(decode3, &resultx)
		if resultx.Rol == rol.Rol {
			for i := range resultx.Module {

				if resultx.Module[i].Acl == rol.Acl {

					arrayPerm = resultx.Module[i].Value
				}
			}
		}
	}

	return arrayPerm
}

func ArrayPerm(c *gin.Context, rol Rol) []string {

	var arrayPerm []string
	user := c.MustGet("user").(string)
	// fmt.Println("##user: ", user)
	str := []byte(user)
	var raw User
	json.Unmarshal(str, &raw)

	// cloudantUrl := config.StrNoSQLDrive()
	// cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	// if err != nil {
	// 	log.Println("Can not connect to Cloudant database")
	// }

	// //ensure db exists
	// //if the db exists the db will be returned anyway
	// dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)
	var datas Response

	// if cloudantUrl == "" {
	// 	c.JSON(200, gin.H{})
	// 	//return
	// }

	// var result model.AlldocsResult
	// var result map[string]interface{}
	// errs := cloudant.DB(dbName).Get(raw.ID, &result, couchdb.Options{"include_docs": true})
	// var result map[string]interface{}

	var arrKey = []string{"profiles", "_id", "_rev"}

	query := model.QuerySelectorAll{
		Selector: map[string]interface{}{
			"meta": arrKey[0],
		},
		Fields: arrKey,
	}

	respText := FindDataAll(query)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.ProfileDocumentsArray
	json.Unmarshal(decode, &results)

	var pr []model.Permission

	for i := range results.Doc {
		if results.Doc[i].Profile.UserID == raw.ID {
			pr = append(pr, model.Permission{
				ID: results.Doc[i].Profile.PermissionID,
			})
		}
	}

	if len(pr) == 0 {
		datas = Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	} else {

		// find := pr[0]

		// fmt.Println("##find: ", (find.ID))

		// var result2 map[string]interface{}

		// errs2 := cloudant.DB(dbName).Get(find.ID, &result2, couchdb.Options{"include_docs": true})
		// if errs2 != nil {
		// 	fmt.Println("##errs2: ", errs2)
		// }

		// jsonToString2, _ := json.Marshal(result2)
		// // fmt.Println("##jsonToString: ", string(jsonToString2))
		// decode2 := []byte(jsonToString2)

		// var resultPerm model.PermissionDoc
		// json.Unmarshal(decode2, &resultPerm)
		// var perms model.Permission
		find := pr[0]

		var arrKey1 = []string{"permissions", "_id", "_rev"}
		query1 := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey1[0],
				"_id":  find.ID,
			},
			Fields: arrKey1,
		}

		respText1 := FindDataAll(query1)

		jsonToString1 := (respText1)
		decode1 := []byte(jsonToString1)
		var results1 model.PermissionDocumentsArray
		json.Unmarshal(decode1, &results1)

		resultPerm := results1.Doc[0]

		var perms model.Permission

		if resultPerm.ID == find.ID {

			// fmt.Println("##resultPerm.Permission.ObjectModulePermission: ",
			// 	(resultPerm.Permission.ObjectModulePermission))

			decodeModules := []byte(resultPerm.Permission.ObjectModulePermission)
			var modules interface{}
			json.Unmarshal(decodeModules, &modules)

			perms = model.Permission{
				Modules: modules,
			}
		}

		jsonToString3, _ := json.Marshal(perms.Modules)
		// fmt.Println("##jsonToString3: ", string(jsonToString3))
		decode3 := []byte(jsonToString3)

		var resultx Roles
		json.Unmarshal(decode3, &resultx)

		for i := range resultx.Module {

			if resultx.Module[i].Acl == rol.Acl {

				arrayPerm = resultx.Module[i].Value
			}
		}
	}

	return arrayPerm
}

func ArrayPermMenu(c *gin.Context) []Module {

	// cloudantUrl := config.StrNoSQLDrive()
	// cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	// if err != nil {
	// 	log.Println("Can not connect to Cloudant database")
	// }

	// //ensure db exists
	// //if the db exists the db will be returned anyway
	// dbName := config.StrNoSQLDBname()
	var datas Response
	var arrayPerms []Module
	user := c.MustGet("user").(string)
	// fmt.Println("##user: ", user)
	str := []byte(user)
	var raw User
	json.Unmarshal(str, &raw)
	var arrKey = []string{"profiles"}

	query := model.QuerySelectorAll{
		Selector: map[string]interface{}{
			"meta": arrKey[0],
		},
		Fields: arrKey,
	}

	respText := FindDataAll(query)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.ProfileDocumentsArray
	json.Unmarshal(decode, &results)

	var pr []model.Permission
	// fmt.Println("##pr: ", (results.Doc))

	for i := range results.Doc {
		// fmt.Println("##pr: ", (results))
		// fmt.Println("##results.Doc[i].Profile.UserID#0: ", (results.Doc[i].Profile.UserID))
		if results.Doc[i].Profile.UserID == raw.ID {
			// fmt.Println("##results.Doc[i].Profile.UserID#1: ", (results.Doc[i].Profile.UserID))
			pr = append(pr, model.Permission{
				ID: results.Doc[i].Profile.PermissionID,
			})
		}
	}

	// fmt.Println("##pr: ", (results0))
	// fmt.Println("##pr: ", (pr[0]))

	if len(pr) == 0 {
		datas = Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	} else {
		find := pr[0]

		var arrKey1 = []string{"permissions", "_id", "_rev"}
		query1 := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey1[0],
				"_id":  find.ID,
			},
			Fields: arrKey1,
		}

		respText1 := FindDataAll(query1)

		// fmt.Println("##respText1: ", (respText1))

		jsonToString1 := (respText1)
		decode1 := []byte(jsonToString1)
		var results1 model.PermissionDocumentsArray
		json.Unmarshal(decode1, &results1)

		// fmt.Println("##results1: ", (results1))
		resultPerm := results1.Doc[0]
		// fmt.Println("##find: ", (find.ID))

		// var result2 map[string]interface{}

		// errs2 := cloudant.DB(dbName).Get(find.ID, &result2, couchdb.Options{"include_docs": true})
		// if errs2 != nil {
		// 	fmt.Println("##errs2: ", errs2)
		// }

		// jsonToString2, _ := json.Marshal(result2)
		// // fmt.Println("##jsonToString: ", string(jsonToString2))
		// decode2 := []byte(jsonToString2)

		// var resultPerm model.PermissionDoc
		// json.Unmarshal(decode2, &resultPerm)
		var perms model.Permission

		if resultPerm.ID == find.ID {

			decodeModules := []byte(resultPerm.Permission.ObjectModulePermission)
			var modules interface{}
			json.Unmarshal(decodeModules, &modules)

			perms = model.Permission{
				Modules: modules,
			}
		}

		jsonToString3, _ := json.Marshal(perms.Modules)
		// fmt.Println("##jsonToString3: ", string(jsonToString3))
		decode3 := []byte(jsonToString3)

		var resultx Roles
		json.Unmarshal(decode3, &resultx)
		for i := range resultx.Module {
			// fmt.Println("##resultx: ", (resultx.Module[i]))
			if resultx.Module[i].Visible == true {
				// fmt.Println("##resultx visible: ", (resultx.Module[i]))
				arrayPerms = append(arrayPerms, Module{
					Url:          resultx.Module[i].Url,
					IconWeb:      resultx.Module[i].IconWeb,
					Mobile:       resultx.Module[i].Mobile,
					Android:      resultx.Module[i].Android,
					Ios:          resultx.Module[i].Ios,
					Ionic:        resultx.Module[i].Ionic,
					Name:         resultx.Module[i].Name,
					LangProperty: resultx.Module[i].LangProperty,
					Acl:          resultx.Module[i].Acl,
					Visible:      resultx.Module[i].Visible,
					Value:        resultx.Module[i].Value,
				})
			}
		}
	}

	return arrayPerms
}

func GetIDsByRol(c *gin.Context, rol Rol) []model.RolesDev {
	//
	var rol_dev []model.RolesDev
	var ts []model.Permission

	var arrKey = []string{"permissions", "_id", "_rev"}

	query := model.QuerySelectorAll{
		Selector: map[string]interface{}{
			"meta": arrKey[0],
		},
		Fields: arrKey,
	}

	respText := FindDataAll(query)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.PermissionDocumentsArray
	json.Unmarshal(decode, &results)

	for i := range results.Doc {
		if results.Doc[i].Permission.ObjectModulePermission != "" {
			ts = append(ts, model.Permission{
				IDs:                    results.Doc[i].ID,
				Rev:                    results.Doc[i].Rev,
				ID:                     results.Doc[i].ID,
				ObjectModulePermission: results.Doc[i].Permission.ObjectModulePermission,
				CreatedAt:              results.Doc[i].Permission.CreatedAt,
				UpdatedAt:              results.Doc[i].Permission.UpdatedAt,
			})
		}
	}

	perm := ts
	// fmt.Println("##perm: ", perm)

	if len(perm) > 0 {
		for i := range perm {
			objPerm := perm[i].ObjectModulePermission
			// fmt.Println("##objPerm: ", objPerm)
			decode := []byte(objPerm)
			var modules interface{}
			json.Unmarshal(decode, &modules)
			//
			var modul RolesDev
			json.Unmarshal(decode, &modul)
			//
			// fmt.Println("##modul: ", modul)

			perms := model.RolesDev{
				IDs:    perm[i].IDs,
				Rev:    perm[i].Rev,
				Rol:    modul.Rol,
				Module: modules,
			}

			if perms.Rol == rol.Rol {
				rol_dev = append(rol_dev, model.RolesDev{
					IDs:    perms.IDs,
					Rev:    perms.Rev,
					Rol:    perms.Rol,
					Module: perms.Module,
				})
			}
		}
	} else {
		fmt.Println("##error al procesar la data ")
	}

	return rol_dev
}

func CloudantDefault() *couchdb.Client {
	cloudantUrl := config.StrNoSQLDrive()
	cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	if err != nil {
		log.Println("Can not connect to Cloudant database")
	}

	return cloudant
}

func CloudantOther(drive string) *couchdb.Client {
	cloudantUrl := drive
	cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	if err != nil {
		log.Println("Can not connect to Cloudant database")
	}

	return cloudant
}

func CurlPost(payload interface{}, uri string) string {
	strToken, _ := json.Marshal(payload)
	jsonstr := []byte(strToken)
	remote_host0 := config.EnviromentsRaw().RemoteHost[0].Name
	//url := "http://localhost:1324/api/instance/create"
	url := remote_host0 + uri
	contentType := "application/json"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", contentType)
	client := http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	//fmt.Println("##response remote: ", resp)

	body, _ := ioutil.ReadAll(resp.Body)

	strBody := string(body)
	return strBody
}

func CurlGet(uri string) string {
	remote_host0 := config.EnviromentsRaw().RemoteHost[0].Name
	//url := "http://localhost:1324/api/instance/create"
	url := remote_host0 + uri
	contentType := "application/json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", contentType)
	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	//fmt.Println("##response remote: ", resp)
	body, _ := ioutil.ReadAll(resp.Body)

	strBody := string(body)
	return strBody
}

func CurlBodyJSON(method string, uri string, payload interface{}) string {

	url := uri
	contentType := "application/json"

	var btes io.Reader
	if payload != nil {
		strToken, _ := json.Marshal(payload)
		jsonstr := []byte(strToken)
		btes = bytes.NewBuffer(jsonstr)
	} else {
		btes = nil
	}

	req, _ := http.NewRequest(method, url, btes)
	req.Header.Set("Content-Type", contentType)
	client := http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	//fmt.Println("##response remote: ", resp)

	body, _ := ioutil.ReadAll(resp.Body)

	strBody := string(body)
	return strBody
}

func FindDataIndex(payload model.QuerySelectorPaginateIndex) string {
	cloudantUrl := config.StrNoSQLDrive()
	dbName := config.StrNoSQLDBname()
	uri := cloudantUrl + "/" + dbName + "/_find"
	resultData := CurlBodyJSON("POST", uri, payload)
	return resultData
}

func FindDataPaginate(payload model.QuerySelectorPaginate) string {
	cloudantUrl := config.StrNoSQLDrive()
	dbName := config.StrNoSQLDBname()
	uri := cloudantUrl + "/" + dbName + "/_find"
	resultData := CurlBodyJSON("POST", uri, payload)
	return resultData
}

func FindDataAll(payload model.QuerySelectorAll) string {
	cloudantUrl := config.StrNoSQLDrive()
	dbName := config.StrNoSQLDBname()
	uri := cloudantUrl + "/" + dbName + "/_find"
	resultData := CurlBodyJSON("POST", uri, payload)
	return resultData
}

func FindDataInterface(payload interface{}) string {
	cloudantUrl := config.StrNoSQLDrive()
	dbName := config.StrNoSQLDBname()
	uri := cloudantUrl + "/" + dbName + "/_find"
	resultData := CurlBodyJSON("POST", uri, payload)
	return resultData
}

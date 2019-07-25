package security

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"../../config"
	"../../model"
	"../../util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/timjacobi/go-couchdb"
)

func Register(c *gin.Context) {

	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	var (
		// result model.AlldocsResult
		// doc    model.Document
		// docs   []model.Document
		datas util.Response
		user  model.User
	)

	if c.BindJSON(&user) == nil {
		var arrKey = []string{"users", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta":        arrKey[0],
				"users.email": user.Email,
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var result model.UserDocumentsArray
		json.Unmarshal(decode, &result)

		var us []model.User
		for i := range result.Doc {

			fmt.Println("##i: ", result.Doc[i])
			us = append(us, model.User{
				IDs:      result.Doc[i].ID,
				Rev:      result.Doc[i].Rev,
				ID:       result.Doc[i].ID,
				Username: result.Doc[i].User.Username,
				Email:    result.Doc[i].User.Email,
				Password: result.Doc[i].User.Password,
			})

		}

		if len(us) != 0 {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {

			byteP := []byte(user.Password)
			p := util.HashAndSalt(byteP)
			tt := model.User{
				Avatar:   user.Avatar,
				Username: user.Username,
				Email:    user.Email,
				Password: p,
			}

			if user.Username == "" && user.Email == "" && user.Password == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"users"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "users": tt})
				c.JSON(200, user)
			}
		}
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}
}

func Login(c *gin.Context) {
	var (
		datas util.Response
		user  model.User
	)
	if c.BindJSON(&user) == nil {

		if user.Username == "" && user.Email == "" && user.Password == "" {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			var arrKey = []string{"users", "_id", "_rev"}

			query := model.QuerySelectorAll{
				Selector: map[string]interface{}{
					"meta":        arrKey[0],
					"users.email": user.Email,
				},
				Fields: arrKey,
			}

			respText := util.FindDataAll(query)

			jsonToString := (respText)
			decode := []byte(jsonToString)
			var results model.UserDocumentsArray
			json.Unmarshal(decode, &results)

			var us []model.User
			for i := range results.Doc {
				fmt.Println("##i: ", results.Doc[i])
				us = append(us, model.User{
					IDs:      results.Doc[i].ID,
					Rev:      results.Doc[i].Rev,
					ID:       results.Doc[i].ID,
					Username: results.Doc[i].User.Username,
					Email:    results.Doc[i].User.Email,
					Password: results.Doc[i].User.Password,
				})
			}

			if len(results.Doc) == 0 {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			}

			if len(us) == 0 {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			}

			find := us[0]
			byteP := []byte(user.Password)

			if util.ComparePasswords(find.Password, byteP) {
				id := find.ID
				claims := util.CustomClaims{
					ID:    id,
					Email: find.Email,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
					},
				}
				tokenStr := util.SetToken(claims)

				strTo := util.Token{Token: tokenStr}

				datas = util.Response{
					true,
					"ok",
					strTo,
				}
			} else {
				datas = util.Response{
					false,
					"error_email_password_mismatch",
					nil,
				}
			}

			c.JSON(200, datas)
		}
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}

}

func EditUser(c *gin.Context) {

	var (
		// result model.AlldocsResult
		// doc    model.Document
		// docs   []model.Document
		datas util.Response
		user  model.User
	)

	if c.BindJSON(&user) == nil {

		var arrKey = []string{"users", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta":        arrKey[0],
				"users.email": user.Email,
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.UserDocumentsArray
		json.Unmarshal(decode, &results)

		var us []model.User
		for i := range results.Doc {

			fmt.Println("##i: ", results.Doc[i])
			us = append(us, model.User{
				IDs:      results.Doc[i].ID,
				Rev:      results.Doc[i].Rev,
				ID:       results.Doc[i].ID,
				Username: results.Doc[i].User.Username,
				Email:    results.Doc[i].User.Email,
				Password: results.Doc[i].User.Password,
			})

		}

		if len(results.Doc) == 0 {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		}

		if len(us) == 0 {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			find := us[0]

			datas = util.Response{
				true,
				"ok",
				find,
			}

			c.JSON(200, datas)
		}
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}
}

func PermissionFaker(c *gin.Context) {
	var datas util.Response
	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	var result model.AlldocsResult
	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	data0, err := ioutil.ReadFile("FileSystem/permissionFaker.json")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("#strData: ", string(data))

	data1, err1 := ioutil.ReadFile("FileSystem/permissionFakerOld.json")
	if err1 != nil {
		fmt.Println(err1)
	}

	var arrKey = []string{"permissions"}
	errs := cloudant.DB(dbName).AllDocs(&result, couchdb.Options{"include_docs": true, "startkey": arrKey})
	if errs != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch docs"})
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}

	jsonToString, _ := json.Marshal(result.Rows)
	// fmt.Println("##jsonToString: ", string(jsonToString))

	decode := []byte(jsonToString)
	var results []model.PermissionDocument
	json.Unmarshal(decode, &results)

	var ts []model.Permission

	for i := range results {
		// fmt.Println("##doc: ", results[i].Doc)
		// fmt.Println("##i: ", results[i].Doc.Test)
		if results[i].Doc.Permission.ID != "" {
			ts = append(ts, model.Permission{
				IDs:                    results[i].Doc.ID,
				Rev:                    results[i].Doc.Rev,
				ID:                     results[i].Doc.ID,
				ObjectModulePermission: results[i].Doc.Permission.ObjectModulePermission,
				CreatedAt:              results[i].Doc.Permission.CreatedAt,
				UpdatedAt:              results[i].Doc.Permission.UpdatedAt,
			})
		}
	}

	// var perm []model.Permission
	// count := "*"
	// db.Debug().Find(&perm).Count(&count)
	perm := ts

	decode0 := []byte(string(data0))
	var result0 []interface{}
	json.Unmarshal(decode0, &result0)

	decode1 := []byte(string(data1))
	var result1 []interface{}
	json.Unmarshal(decode1, &result1)

	if len(perm) == 0 {
		for i := range result0 {
			fmt.Println("##i: ", result0[i])
			js, _ := json.Marshal(result0[i])
			fmt.Println("##js: ", string(js))

			t := model.Permission{
				ObjectModulePermission: string(js),
			}

			// db.Debug().Create(&t)
			var arrKey = []string{"permissions"}
			cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "permissions": t})
		}
	} else {
		for i := range result1 {
			js1, _ := json.Marshal(result1[i])
			// var perms model.Permission
			// db.Debug().Where(map[string]interface{}{"object_module_permission": string(js1)}).First(&perms)
			permissionx := result1
			for j := range permissionx {
				fmt.Println(j)
				fmt.Println("##perms. : ", string(js1))
				// fmt.Println("##perms.id: ", perms.ID)
				// perms.ObjectModulePermission = string(js1)
			}

			// db.Save(&perms)
			arrayData, _ := json.Marshal(permissionx)
			byteArrayData := []byte(arrayData)

			err := ioutil.WriteFile("FileSystem/permissionFakerOld.json", byteArrayData, 0777)
			if err != nil {
				fmt.Println("err file: ", err)
			}
		}

	}

	datas = util.Response{
		true,
		"ok",
		"",
	}

	c.JSON(http.StatusAccepted, datas)
}

func PermissionDeleteFaker(c *gin.Context) {
	var datas util.Response
	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	var result model.AlldocsResult
	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	errs := cloudant.DB(dbName).AllDocs(&result, couchdb.Options{"include_docs": true})
	if errs != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch docs"})
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}

	jsonToString, _ := json.Marshal(result.Rows)
	// fmt.Println("##jsonToString: ", string(jsonToString))

	decode := []byte(jsonToString)
	var results []model.PermissionDocument
	json.Unmarshal(decode, &results)

	var ts []model.Permission

	for i := range results {
		// fmt.Println("##doc: ", results[i].Doc)
		// fmt.Println("##i: ", results[i].Doc.Permission)
		if results[i].Doc.Permission.ObjectModulePermission != "" {
			ts = append(ts, model.Permission{
				IDs:                    results[i].Doc.ID,
				Rev:                    results[i].Doc.Rev,
				ID:                     results[i].Doc.ID,
				ObjectModulePermission: results[i].Doc.Permission.ObjectModulePermission,
				CreatedAt:              results[i].Doc.Permission.CreatedAt,
				UpdatedAt:              results[i].Doc.Permission.UpdatedAt,
			})
		}
	}

	for i := range ts {
		// fmt.Println("##ts[i].IDs: ", ts[i].IDs)
		cloudant.DB(dbName).Delete(ts[i].IDs, ts[i].Rev)
	}

	datas = util.Response{
		true,
		"ok",
		nil,
	}
	c.JSON(200, datas)

}

func PermissionAll(c *gin.Context) {
	var datas util.Response

	var arrKey = []string{"permissions", "_id", "_rev"}

	query := model.QuerySelectorAll{
		Selector: map[string]interface{}{
			"meta": arrKey[0],
		},
		Fields: arrKey,
	}

	respText := util.FindDataAll(query)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.PermissionDocumentsArray
	json.Unmarshal(decode, &results)

	var ts []model.Permission

	for i := range results.Doc {
		fmt.Println("##docs: ", results.Doc[i])
		// fmt.Println("##i: ", results[i].Doc.Permission)
		if results.Doc[i].Permission.ObjectModulePermission != "" {

			// fmt.Println("##ObjectModulePermission: ", results[i].Doc.Permission.ObjectModulePermission)

			decodeModules := []byte(results.Doc[i].Permission.ObjectModulePermission)
			// var modules map[string]interface{}
			var modules interface{}
			json.Unmarshal(decodeModules, &modules)

			ts = append(ts, model.Permission{
				IDs:                    results.Doc[i].ID,
				Rev:                    results.Doc[i].Rev,
				ID:                     results.Doc[i].ID,
				ObjectModulePermission: results.Doc[i].Permission.ObjectModulePermission,
				Modules:                modules,
				CreatedAt:              results.Doc[i].Permission.CreatedAt,
				UpdatedAt:              results.Doc[i].Permission.UpdatedAt,
			})
		}
	}

	if len(results.Doc) == 0 {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	} else {

		datas = util.Response{
			true,
			"ok",
			ts,
		}

		c.JSON(200, datas)
	}

}

func AssignRoles(c *gin.Context) {
	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	var (
		datas util.Response
		t     model.Profile
	)

	if c.BindJSON(&t) == nil {
		// cloudant.DB(dbName).Put(id, map[string]interface{}{"tests": t}, rev

		if t.UserID == "" && t.PermissionID == "" {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			var arrKey = []string{"profiles"}
			cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "profiles": t})
			datas = util.Response{
				true,
				"ok",
				t,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}
}

func EditRoles(c *gin.Context) {
	id := c.Param("id")

	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)
	var datas util.Response

	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	// var result model.AlldocsResult
	var result map[string]interface{}
	errs := cloudant.DB(dbName).Get(id, &result, couchdb.Options{"include_docs": true})
	if errs != nil {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch docs"})
	} else {

		jsonToString, _ := json.Marshal(result)
		// fmt.Println("##jsonToString: ", string(jsonToString))
		decode := []byte(jsonToString)

		var results model.ProfileDoc
		json.Unmarshal(decode, &results)

		var result1 map[string]interface{}
		var result2 map[string]interface{}

		errs1 := cloudant.DB(dbName).Get(results.Profile.UserID, &result1, couchdb.Options{"include_docs": true})
		if errs1 != nil {
			fmt.Println("##errs1: ", errs1)
		}

		errs2 := cloudant.DB(dbName).Get(results.Profile.PermissionID, &result2, couchdb.Options{"include_docs": true})
		if errs2 != nil {
			fmt.Println("##errs2: ", errs2)
		}

		jsonToString1, _ := json.Marshal(result1)
		//  fmt.Println("##jsonToString1: ", string(jsonToString1))
		decode1 := []byte(jsonToString1)

		jsonToString2, _ := json.Marshal(result2)
		//  fmt.Println("##jsonToString2: ", string(jsonToString2))
		decode2 := []byte(jsonToString2)

		var resultUsers model.UserDoc
		json.Unmarshal(decode1, &resultUsers)
		var us model.User

		//for i := range resultUsers {
		// fmt.Println("#--<resultUsers.Doc.ID: ", resultUsers.ID)
		// fmt.Println("#--<results.Profile.UserID: ", results.Profile.UserID)
		if resultUsers.ID == results.Profile.UserID {
			// fmt.Println("#--resultUsers.Doc.User.ID: ", resultUsers.ID)
			us = model.User{
				IDs:      resultUsers.ID,
				Rev:      resultUsers.Rev,
				ID:       resultUsers.ID,
				Username: resultUsers.User.Username,
				Email:    resultUsers.User.Email,
				// Password: resultUsers[i].Doc.User.Password,
			}
		}
		//}

		var resultPerm model.PermissionDoc
		json.Unmarshal(decode2, &resultPerm)
		var perms model.Permission

		//for i := range resultPerm {

		if resultPerm.ID == results.Profile.PermissionID {
			decodeModules := []byte(resultPerm.Permission.ObjectModulePermission)
			// var modules map[string]interface{}
			var modules interface{}
			json.Unmarshal(decodeModules, &modules)

			perms = model.Permission{
				IDs:                    resultPerm.ID,
				Rev:                    resultPerm.Rev,
				ID:                     resultPerm.ID,
				ObjectModulePermission: resultPerm.Permission.ObjectModulePermission,
				Modules:                modules,
			}
		}

		find := us
		find1 := perms

		fmt.Println("##find: ", (find))
		fmt.Println("##find2: ", (find1))

		// var results model.ProfileDoc
		// json.Unmarshal(decode, &results)

		t := model.Profile{
			IDs:        results.ID,
			Rev:        results.Rev,
			ID:         results.ID,
			User:       find,
			Permission: find1,
			CreatedAt:  results.Profile.CreatedAt,
			UpdatedAt:  results.Profile.UpdatedAt,
		}

		datas = util.Response{
			true,
			"ok",
			t,
		}

		c.JSON(200, datas)

		// c.JSON(200, t)
	}
}

func GetRolByUser(c *gin.Context) {
	id := c.Param("user_id")

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"profiles", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta":             arrKey[0],
				"profiles.user_id": id,
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var result model.ProfileDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.Profile{
			IDs:          results.ID,
			Rev:          results.Rev,
			ID:           results.ID,
			UserID:       results.Profile.UserID,
			PermissionID: results.Profile.PermissionID,
		}

		datas = util.Response{
			true,
			"ok",
			t,
		}

		c.JSON(200, datas)

	} else {
		datas = util.IsEdit(c, rol)

		c.JSON(200, datas)
	}
}

func UpdateRoles(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	var (
		datas util.Response
		t     model.Profile
	)

	if c.BindJSON(&t) == nil {

		if t.UserID == "" && t.PermissionID == "" {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			var arrKey = []string{"profiles"}
			cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "profiles": t}, rev)
			datas = util.Response{
				true,
				"ok",
				t,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}
}

func DeleteRoles(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	cloudantUrl := config.StrNoSQLDrive()
	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway
	dbName := config.StrNoSQLDBname()
	// cloudant.CreateDB(dbName)

	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	cloudant.DB(dbName).Delete(id, rev)
	var datas util.Response
	datas = util.Response{
		true,
		"ok",
		nil,
	}
	c.JSON(200, datas)
}

func Menu(c *gin.Context) {
	var datas util.Response

	menu := util.ArrayPermMenu(c)
	fmt.Println("#menu: ", menu)
	if len(menu) != 0 {
		datas = util.Response{
			true,
			"ok",
			menu,
		}
	} else {
		datas = util.Response{
			false,
			"menu empty",
			"",
		}
	}

	c.JSON(http.StatusAccepted, datas)
}

func SetupPermission(c *gin.Context) {
	var datas util.Response
	strResponse0 := util.CurlGet("/api/testdb/all")
	fmt.Println("##strResponse0: ", string(strResponse0))

	strResponse := util.CurlGet("/api/permission/faker")
	fmt.Println("##strResponse: ", string(strResponse))
	strTk, _ := json.Marshal(strResponse)
	decode := []byte(strTk)
	var results interface{}
	json.Unmarshal(decode, &results)

	datas = util.Response{
		true,
		"ok",
		results,
	}

	c.JSON(http.StatusAccepted, datas)
}

func SetupAdmin(c *gin.Context) {

	var (
		rol   util.Rol
		datas util.Response
		t     model.User
	)
	rol.Rol = "admin"

	if c.BindJSON(&t) == nil {

		strResponse0 := util.CurlPost(t, "/api/87780FA5DE684E87CB92B279F0BC07B14F572851E73B8943A097C1770A5F38E6")
		fmt.Println("##strResponse0: ", string(strResponse0))

		objUser := map[string]interface{}{
			"email": t.Email,
		}

		g := util.GetIDsByRol(c, rol)
		p := g[0]
		fmt.Println("##GetIDsByRol: ", (p.IDs))

		strResponse1 := util.CurlPost(objUser, "/api/edit-user")
		fmt.Println("##strResponse1: ", string(strResponse1))

		decode := []byte(strResponse1)
		var results util.ResponseU
		json.Unmarshal(decode, &results)

		// fmt.Println("##strResponse data1: ", (results))
		// fmt.Println("##strResponse   1: ", (results.Data.IDs))

		objUserRoles := map[string]interface{}{
			"user_id":       results.Data.IDs,
			"permission_id": p.IDs,
		}

		strResponse2 := util.CurlPost(objUserRoles, "/api/roles/assign")
		fmt.Println("##strResponse2: ", string(strResponse2))

		datas = util.Response{
			true,
			"ok",
			t,
		}
		c.JSON(200, datas)
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}

}

func SetupUser(c *gin.Context) {
	var (
		rol   util.Rol
		datas util.Response
		t     model.User
	)
	rol.Rol = "users"

	if c.BindJSON(&t) == nil {

		strResponse0 := util.CurlPost(t, "/api/87780FA5DE684E87CB92B279F0BC07B14F572851E73B8943A097C1770A5F38E6")
		fmt.Println("##strResponse0: ", string(strResponse0))

		objUser := map[string]interface{}{
			"email": t.Email,
		}

		g := util.GetIDsByRol(c, rol)
		p := g[0]
		fmt.Println("##GetIDsByRol: ", (p.IDs))

		strResponse1 := util.CurlPost(objUser, "/api/edit-user")
		fmt.Println("##strResponse1: ", string(strResponse1))

		decode := []byte(strResponse1)
		var results util.ResponseU
		json.Unmarshal(decode, &results)

		// fmt.Println("##strResponse data1: ", (results))
		// fmt.Println("##strResponse   1: ", (results.Data.IDs))

		objUserRoles := map[string]interface{}{
			"user_id":       results.Data.IDs,
			"permission_id": p.IDs,
		}

		strResponse2 := util.CurlPost(objUserRoles, "/api/roles/assign")
		fmt.Println("##strResponse2: ", string(strResponse2))

		datas = util.Response{
			true,
			"ok",
			t,
		}
		c.JSON(200, datas)
	} else {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	}
}

func AddUsers(c *gin.Context) {
	var (
		rol   util.Rol
		role  util.Rol
		datas util.Response
		t     model.PayloadUserRoles
		u     model.User
	)
	// rol.Rol = "users"
	///////////////////////
	role.Rol = "admin"
	role.Acl = "register_users"
	if util.IsRead(c, role).Success == true {

		if c.BindJSON(&t) == nil {
			u = model.User{
				Avatar:   t.Avatar,
				Username: t.Username,
				Email:    t.Email,
				Password: t.Password,
			}
			strResponse0 := util.CurlPost(u, "/api/87780FA5DE684E87CB92B279F0BC07B14F572851E73B8943A097C1770A5F38E6")
			fmt.Println("##strResponse0: ", string(strResponse0))

			objUser := map[string]interface{}{
				"email": t.Email,
			}

			rol.Rol = t.Rol
			g := util.GetIDsByRol(c, rol)
			p := g[0]
			fmt.Println("##GetIDsByRol: ", (p.IDs))

			strResponse1 := util.CurlPost(objUser, "/api/edit-user")
			fmt.Println("##strResponse1: ", string(strResponse1))

			decode := []byte(strResponse1)
			var results util.ResponseU
			json.Unmarshal(decode, &results)

			// fmt.Println("##strResponse data1: ", (results))
			// fmt.Println("##strResponse   1: ", (results.Data.IDs))

			objUserRoles := map[string]interface{}{
				"user_id":       results.Data.IDs,
				"permission_id": p.IDs,
			}

			strResponse2 := util.CurlPost(objUserRoles, "/api/roles/assign")
			fmt.Println("##strResponse2: ", string(strResponse2))

			datas = util.Response{
				true,
				"ok",
				t,
			}
			c.JSON(200, datas)
		} else {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.IsRead(c, role)

		c.JSON(200, datas)
	}
}

func SearchPaginateUsers(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	skips, errk := strconv.ParseUint(skip, 10, 32)
	if errk != nil {
		fmt.Println(errk)
	}

	limits, errl := strconv.ParseUint(limit, 10, 32)
	if errl != nil {
		fmt.Println(errl)
	}

	sk := uint(skips)
	lm := uint(limits)

	fmt.Println("##sk: ", sk)
	fmt.Println("##lm: ", lm)

	var rol util.Rol
	rol.Rol = "admin"
	rol.Acl = "register_users"
	var datas util.Response
	var t model.User

	if util.IsRead(c, rol).Success == true {
		if c.BindJSON(&t) == nil {
			var arrKey = []string{"users", "_id", "_rev", "meta", "tag"}
			var arrTag = []string{"users"}
			input := []byte(`[{"_id": "desc"}]`)
			var sort []interface{}
			errj := json.Unmarshal(input, &sort)
			if errj != nil {
				log.Fatal(errj)
			}

			query := map[string]interface{}{
				"selector": map[string]interface{}{
					"meta": arrKey[0],
					"users.username": map[string]interface{}{
						"$gt": t.Username,
					},
					"tag": arrTag,
				},
				"limit":  lm,
				"skip":   sk,
				"fields": arrKey,
				"sort":   sort,
			}

			fmt.Println(query)

			str, _ := json.Marshal(query)
			decodex := []byte(str)

			var resultsx interface{}
			json.Unmarshal(decodex, &resultsx)

			fmt.Println("##str: ", resultsx)

			respText := util.FindDataInterface(resultsx)

			jsonToString := (respText)
			decode := []byte(jsonToString)
			var results model.UserDocumentsArray
			json.Unmarshal(decode, &results)

			var ts []model.User

			for i := range results.Doc {
				ts = append(ts, model.User{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Username:  results.Doc[i].User.Username,
					Email:     results.Doc[i].User.Email,
					CreatedAt: results.Doc[i].User.CreatedAt,
					UpdatedAt: results.Doc[i].User.UpdatedAt,
				})
			}

			if len(results.Doc) == 0 {
				datas = util.Response{
					true,
					"empty_data",
					results.Doc,
				}
				c.JSON(200, datas)
				return
			} else {

				datas = util.Response{
					true,
					"ok",
					ts,
				}
				c.JSON(200, datas)
				return
			}
		} else {
			datas = util.Response{
				false,
				"empty_data",
				nil,
			}
			c.JSON(200, datas)
			return
		}

	} else {
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
		// return
	}
}

func PaginateUsers(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	skips, errk := strconv.ParseUint(skip, 10, 32)
	if errk != nil {
		fmt.Println(errk)
	}

	limits, errl := strconv.ParseUint(limit, 10, 32)
	if errl != nil {
		fmt.Println(errl)
	}

	sk := uint(skips)
	lm := uint(limits)

	fmt.Println("##sk: ", sk)
	fmt.Println("##lm: ", lm)

	var rol util.Rol
	rol.Rol = "admin"
	rol.Acl = "register_users"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"users", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"users"}
		input := []byte(`[{"_id": "desc"}]`)
		var sort []interface{}
		errj := json.Unmarshal(input, &sort)
		if errj != nil {
			log.Fatal(errj)
		}

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"meta": arrKey[0],
				"tag":  arrTag,
			},
			"limit":  lm,
			"skip":   sk,
			"fields": arrKey,
			"sort":   sort,
		}

		fmt.Println(query)

		str, _ := json.Marshal(query)
		decodex := []byte(str)

		var resultsx interface{}
		json.Unmarshal(decodex, &resultsx)

		fmt.Println("##str: ", resultsx)

		respText := util.FindDataInterface(resultsx)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.UserDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.User

		for i := range results.Doc {
			ts = append(ts, model.User{
				IDs:       results.Doc[i].ID,
				Rev:       results.Doc[i].Rev,
				ID:        results.Doc[i].ID,
				Username:  results.Doc[i].User.Username,
				Email:     results.Doc[i].User.Email,
				CreatedAt: results.Doc[i].User.CreatedAt,
				UpdatedAt: results.Doc[i].User.UpdatedAt,
			})
		}

		if len(results.Doc) == 0 {
			datas = util.Response{
				true,
				"empty_data",
				results.Doc,
			}
			c.JSON(200, datas)
			// return
		} else {

			datas = util.Response{
				true,
				"ok",
				ts,
			}
			c.JSON(200, datas)
		}

	} else {
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func EditUsers(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Rol = "admin"
	rol.Acl = "register_users"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"users", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
				"_id":  id,
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var result model.UserDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.User{
			IDs:       results.ID,
			Rev:       results.Rev,
			ID:        results.ID,
			Avatar:    results.User.Avatar,
			Username:  results.User.Username,
			Email:     results.User.Email,
			CreatedAt: results.User.CreatedAt,
			UpdatedAt: results.User.UpdatedAt,
		}

		datas = util.Response{
			true,
			"ok",
			t,
		}

		c.JSON(200, datas)

	} else {
		datas = util.IsEdit(c, rol)

		c.JSON(200, datas)
	}
}

func PutUsers(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	var roles util.Rol
	rol.Rol = "admin"
	rol.Acl = "register_users"
	var datas util.Response

	if util.IsUpdate(c, rol).Success == true {

		cloudantUrl := config.StrNoSQLDrive()
		cloudant := util.CloudantDefault()

		//ensure db exists
		//if the db exists the db will be returned anyway
		dbName := config.StrNoSQLDBname()
		// cloudant.CreateDB(dbName)

		if cloudantUrl == "" {
			c.JSON(200, gin.H{})
			return
		}

		var (
			// datas util.Response
			// t1 model.User
			t model.PayloadUserRoles
			// u  model.User
		)

		if c.BindJSON(&t) == nil {

			if t.Username == "" && t.Email == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {

				roles.Rol = t.Rol
				objUser := map[string]interface{}{
					"email": t.Email,
				}

				byteP := []byte(t.Password)
				pass := util.HashAndSalt(byteP)
				payloadUser := map[string]interface{}{
					"avatar":   t.Avatar,
					"username": t.Username,
					"email":    t.Email,
					"password": pass,
				}

				g := util.GetIDsByRol(c, roles)
				p := g[0]
				fmt.Println("##GetIDsByRol: ", (p.IDs))

				strResponse1 := util.CurlPost(objUser, "/api/edit-user")
				fmt.Println("##strResponse1: ", string(strResponse1))

				decode := []byte(strResponse1)
				var results util.ResponseU
				json.Unmarshal(decode, &results)

				objUserRoles := map[string]interface{}{
					"user_id":       results.Data.IDs,
					"permission_id": p.IDs,
				}

				strResponse2 := util.CurlGet("/api/roles/by-user/" + results.Data.IDs)
				fmt.Println("##strResponse2: ", string(strResponse2))

				decode1 := []byte(strResponse2)
				var results1 util.ResponseProfile
				json.Unmarshal(decode1, &results1)
				strResponse3 := util.CurlPost(objUserRoles, "/api/roles/update/"+results1.Data.IDs+"/"+results1.Data.Rev)
				// roles/update/:id/:rev
				fmt.Println("##strResponse3: ", string(strResponse3))

				var arrKey = []string{"users"}

				cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "users": payloadUser}, rev)
				datas = util.Response{
					true,
					"ok",
					t,
				}
				c.JSON(200, datas)
			}
		} else {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.IsUpdate(c, rol)

		c.JSON(200, datas)
	}
}

func DeleteUsers(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	var rol util.Rol
	rol.Rol = "admin"
	rol.Acl = "register_users"
	var datas util.Response

	if util.IsDelete(c, rol).Success == true {

		cloudantUrl := config.StrNoSQLDrive()
		cloudant := util.CloudantDefault()

		//ensure db exists
		//if the db exists the db will be returned anyway
		dbName := config.StrNoSQLDBname()
		// cloudant.CreateDB(dbName)

		if cloudantUrl == "" {
			c.JSON(200, gin.H{})
			return
		}

		cloudant.DB(dbName).Delete(id, rev)
		var datas util.Response
		datas = util.Response{
			true,
			"ok",
			nil,
		}
		c.JSON(200, datas)
	} else {
		datas = util.IsDelete(c, rol)

		c.JSON(200, datas)
	}
}

func ListPermissions(c *gin.Context) {
	var datas util.Response
	var arrKey = []string{"permissions", "_id", "_rev", "meta", "tag"}
	var arrTag = []string{"permissions"}
	input := []byte(`[{"_id": "desc"}]`)
	var sort []interface{}
	errj := json.Unmarshal(input, &sort)
	if errj != nil {
		log.Fatal(errj)
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"meta": arrKey[0],
			"tag":  arrTag,
		},
		"fields": arrKey,
		"sort":   sort,
	}

	fmt.Println(query)

	str, _ := json.Marshal(query)
	decodex := []byte(str)

	var resultsx interface{}
	json.Unmarshal(decodex, &resultsx)

	fmt.Println("##str: ", resultsx)

	respText := util.FindDataInterface(resultsx)

	jsonToString := (respText)
	decode := []byte(jsonToString)
	var results model.PermissionDocumentsArray
	json.Unmarshal(decode, &results)

	var ts []model.ObjectModulePermission
	for i := range results.Doc {
		// perm:=results.Doc[i].Permission.ObjectModulePermission
		decodeModules := []byte(results.Doc[i].Permission.ObjectModulePermission)
		var modules model.ObjectModulePermission
		json.Unmarshal(decodeModules, &modules)
		ts = append(ts, model.ObjectModulePermission{
			IDs:    results.Doc[i].ID,
			Rev:    results.Doc[i].Rev,
			Rol:    modules.Rol,
			Module: modules.Module,
		})
	}

	if len(results.Doc) == 0 {
		datas = util.Response{
			false,
			"error_exception",
			nil,
		}
		c.JSON(200, datas)
	} else {

		datas = util.Response{
			true,
			"ok",
			ts,
		}

		c.JSON(200, datas)
	}
}

func InfoUser(c *gin.Context) {
	var datas util.Response
	raw := util.Auth(c)
	if raw.Email == "" {
		datas = util.Response{
			false,
			"error",
			nil,
		}
	} else {

		datas = util.Response{
			true,
			"ok",
			raw,
		}
	}
	c.JSON(200, datas)
}

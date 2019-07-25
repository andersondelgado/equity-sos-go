package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../../config"
	"../../model"
	"../../util"
	"github.com/gin-gonic/gin"
	"github.com/timjacobi/go-couchdb"
)

func SelectDBTest(c *gin.Context) {
	cloudantUrl := config.StrNoSQLDrive()
	// cloudant, err := couchdb.NewClient(cloudantUrl, nil)
	// if err != nil {
	// 	log.Println("Can not connect to Cloudant database")
	// }

	cloudant := util.CloudantDefault()

	//ensure db exists
	//if the db exists the db will be returned anyway //id := c.Param("id")
	dbName := config.StrNoSQLDBname()
	cloudant.CreateDB(dbName)

	var result model.AlldocsResult
	// var result model.Test
	if cloudantUrl == "" {
		c.JSON(200, gin.H{})
		return
	}

	errs := cloudant.DB(dbName).AllDocs(&result, couchdb.Options{"include_docs": true, "skip": 1, "limit": 100})
	// errs := cloudant.DB(dbName).AllDocs(&result, couchdb.Options{"include_docs": false})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch docs"})
	} else {

		c.JSON(200, result.Rows)

	}
}

func SelectTest(c *gin.Context) {

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"tests", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.TestDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Test

		for i := range results.Doc {
			if results.Doc[i].Test.Name != "" {
				ts = append(ts, model.Test{
					IDs:         results.Doc[i].ID,
					Rev:         results.Doc[i].Rev,
					ID:          results.Doc[i].ID,
					Name:        results.Doc[i].Test.Name,
					Description: results.Doc[i].Test.Description,
					CreatedAt:   results.Doc[i].Test.CreatedAt,
					UpdatedAt:   results.Doc[i].Test.UpdatedAt,
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
	} else {
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func PaginateTest(c *gin.Context) {
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
	rol.Acl = "test"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"tests", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"tests"}
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
		var results model.TestDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Test

		for i := range results.Doc {
			ts = append(ts, model.Test{
				IDs:         results.Doc[i].ID,
				Rev:         results.Doc[i].Rev,
				ID:          results.Doc[i].ID,
				Name:        results.Doc[i].Test.Name,
				Description: results.Doc[i].Test.Description,
				CreatedAt:   results.Doc[i].Test.CreatedAt,
				UpdatedAt:   results.Doc[i].Test.UpdatedAt,
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

func SearchPaginateTest(c *gin.Context) {
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
	rol.Acl = "test"
	var datas util.Response
	var t model.Test

	if util.IsRead(c, rol).Success == true {
		if c.BindJSON(&t) == nil {
			var arrKey = []string{"tests", "_id", "_rev", "meta", "tag"}
			var arrTag = []string{"tests"}
			input := []byte(`[{"_id": "desc"}]`)
			var sort []interface{}
			errj := json.Unmarshal(input, &sort)
			if errj != nil {
				log.Fatal(errj)
			}

			query := map[string]interface{}{
				"selector": map[string]interface{}{
					"meta": arrKey[0],
					"tests.name": map[string]interface{}{
						"$gt": t.Name,
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
			var results model.TestDocumentsArray
			json.Unmarshal(decode, &results)

			var ts []model.Test

			for i := range results.Doc {
				ts = append(ts, model.Test{
					IDs:         results.Doc[i].ID,
					Rev:         results.Doc[i].Rev,
					ID:          results.Doc[i].ID,
					Name:        results.Doc[i].Test.Name,
					Description: results.Doc[i].Test.Description,
					CreatedAt:   results.Doc[i].Test.CreatedAt,
					UpdatedAt:   results.Doc[i].Test.UpdatedAt,
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

func EditTest(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"tests", "_id", "_rev"}

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
		var result model.TestDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.Test{
			IDs:         results.ID,
			Rev:         results.Rev,
			ID:          results.ID,
			Name:        results.Test.Name,
			Description: results.Test.Description,
			CreatedAt:   results.Test.CreatedAt,
			UpdatedAt:   results.Test.UpdatedAt,
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

func DeleteTest(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	var rol util.Rol
	rol.Acl = "test"
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

func AddTest(c *gin.Context) {
	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsCreate(c, rol).Success == true {

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
			t     model.Test
		)

		if c.BindJSON(&t) == nil {
			// cloudant.DB(dbName).Post(t)
			if t.Name == "" && t.Description == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"tests"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "tests": t})

				datas = util.Response{
					true,
					"ok",
					t,
				}
				c.JSON(200, datas)
			}
			// c.JSON(200, t)
		} else {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.IsCreate(c, rol)

		c.JSON(200, datas)
	}
}

func PutTest(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	rol.Acl = "test"
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
			datas util.Response
			t     model.Test
		)

		if c.BindJSON(&t) == nil {

			if t.Name == "" && t.Description == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"tests"}

				cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "tests": t}, rev)
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

func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Jokes handler not implemented yet",
	})
}

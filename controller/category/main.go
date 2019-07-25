package category

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"../../config"
	"../../model"
	"../../util"
	"github.com/gin-gonic/gin"
)

func SelectCategory(c *gin.Context) {

	var rol util.Rol
	rol.Acl = "category"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"categories", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.CategoryDocumentsArray
		json.Unmarshal(decode, &results)
		var ts []model.Category

		for i := range results.Doc {
			if results.Doc[i].Category.Name != "" {
				ts = append(ts, model.Category{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Name:      results.Doc[i].Category.Name,
					CreatedAt: results.Doc[i].Category.CreatedAt,
					UpdatedAt: results.Doc[i].Category.UpdatedAt,
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

func PaginateCategory(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")
	var rol util.Rol
	rol.Acl = "category"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

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
		var arrKey = []string{"categories", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"categories"}
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
		var results model.CategoryDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Category

		for i := range results.Doc {
			if results.Doc[i].Category.Name != "" {
				ts = append(ts, model.Category{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Name:      results.Doc[i].Category.Name,
					CreatedAt: results.Doc[i].Category.CreatedAt,
					UpdatedAt: results.Doc[i].Category.UpdatedAt,
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

func SearchPaginateCategory(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	var rol util.Rol
	rol.Acl = "category"
	var datas util.Response
	var t model.Category

	if util.IsRead(c, rol).Success == true {

		// errs := cloudant.DB(dbName).AllDocs(&result, couchdb.Options{"include_docs": true, "startkey": arrKey, "skip": skip, "limit": limit})

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

		var arrKey = []string{"categories", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"categories"}
		input := []byte(`[{"_id": "desc"}]`)
		var sort []interface{}
		errj := json.Unmarshal(input, &sort)
		if errj != nil {
			log.Fatal(errj)
		}

		if c.BindJSON(&t) == nil {
			query := map[string]interface{}{
				"selector": map[string]interface{}{
					"meta": arrKey[0],
					"categories.name": map[string]interface{}{
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
			var results model.CategoryDocumentsArray
			json.Unmarshal(decode, &results)

			// jsonToString, _ := json.Marshal(result.Rows)
			// // fmt.Println("##jsonToString: ", string(jsonToString))

			// decode := []byte(jsonToString)
			// var results []model.ArticleDocument
			// json.Unmarshal(decode, &results)

			var ts []model.Category

			for i := range results.Doc {
				ts = append(ts, model.Category{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Name:      results.Doc[i].Category.Name,
					CreatedAt: results.Doc[i].Category.CreatedAt,
					UpdatedAt: results.Doc[i].Category.UpdatedAt,
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
	}
}

func EditCategory(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "category"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"categories", "_id", "_rev"}

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
		var result model.CategoryDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.Category{
			IDs:       results.ID,
			Rev:       results.Rev,
			ID:        results.ID,
			Name:      results.Category.Name,
			CreatedAt: results.Category.CreatedAt,
			UpdatedAt: results.Category.UpdatedAt,
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

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	var rol util.Rol
	rol.Acl = "category"
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

func AddCategory(c *gin.Context) {
	var rol util.Rol
	rol.Acl = "category"
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
			t     model.Category
		)

		if c.BindJSON(&t) == nil {
			// cloudant.DB(dbName).Post(t)
			if t.Name == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"categories"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "categories": t})

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

func PutCategory(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	rol.Acl = "category"
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

			if t.Name == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"categories"}
				cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "categories": t}, rev)
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

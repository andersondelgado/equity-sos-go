package article

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"../../config"
	"../../model"
	"../../util"
	"github.com/gin-gonic/gin"
)

func SelectArticle(c *gin.Context) {

	var rol util.Rol
	rol.Acl = "article"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"articles", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.ArticleDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Article

		for i := range results.Doc {
			if results.Doc[i].Article.Name != "" {
				ts = append(ts, model.Article{
					IDs:          results.Doc[i].ID,
					Rev:          results.Doc[i].Rev,
					ID:           results.Doc[i].ID,
					Name:         results.Doc[i].Article.Name,
					TypeArticle:  results.Doc[i].Article.TypeArticle,
					Composition:  results.Doc[i].Article.Composition,
					Presentation: results.Doc[i].Article.Presentation,
					CreatedAt:    results.Doc[i].Article.CreatedAt,
					UpdatedAt:    results.Doc[i].Article.UpdatedAt,
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

func PaginateArticle(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	var rol util.Rol
	rol.Acl = "article"
	var datas util.Response

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

		var arrKey = []string{"articles", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"articles"}
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
		var results model.ArticleDocumentsArray
		json.Unmarshal(decode, &results)

		// jsonToString, _ := json.Marshal(result.Rows)
		// // fmt.Println("##jsonToString: ", string(jsonToString))

		// decode := []byte(jsonToString)
		// var results []model.ArticleDocument
		// json.Unmarshal(decode, &results)

		var ts []model.Article

		for i := range results.Doc {
			ts = append(ts, model.Article{
				IDs:          results.Doc[i].ID,
				Rev:          results.Doc[i].Rev,
				ID:           results.Doc[i].ID,
				Name:         results.Doc[i].Article.Name,
				TypeArticle:  results.Doc[i].Article.TypeArticle,
				Composition:  results.Doc[i].Article.Composition,
				Presentation: results.Doc[i].Article.Presentation,
				CreatedAt:    results.Doc[i].Article.CreatedAt,
				UpdatedAt:    results.Doc[i].Article.UpdatedAt,
			})
		}

		if len(results.Doc) == 0 {
			datas = util.Response{
				false,
				"error_exception",
				nil,
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
		}

		// c.JSON(200, ts)

	} else {
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func SearchPaginateArticle(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	var rol util.Rol
	rol.Acl = "article"
	var datas util.Response
	var t model.Article

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

		var arrKey = []string{"articles", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"articles"}
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
					"articles.name": map[string]interface{}{
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
			var results model.ArticleDocumentsArray
			json.Unmarshal(decode, &results)

			// jsonToString, _ := json.Marshal(result.Rows)
			// // fmt.Println("##jsonToString: ", string(jsonToString))

			// decode := []byte(jsonToString)
			// var results []model.ArticleDocument
			// json.Unmarshal(decode, &results)

			var ts []model.Article

			for i := range results.Doc {
				ts = append(ts, model.Article{
					IDs:          results.Doc[i].ID,
					Rev:          results.Doc[i].Rev,
					ID:           results.Doc[i].ID,
					Name:         results.Doc[i].Article.Name,
					TypeArticle:  results.Doc[i].Article.TypeArticle,
					Composition:  results.Doc[i].Article.Composition,
					Presentation: results.Doc[i].Article.Presentation,
					CreatedAt:    results.Doc[i].Article.CreatedAt,
					UpdatedAt:    results.Doc[i].Article.UpdatedAt,
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

func EditArticle(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "article"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"articles", "_id", "_rev"}

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
		var result model.ArticleDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.Article{
			IDs:          results.ID,
			Rev:          results.Rev,
			ID:           results.ID,
			Name:         results.Article.Name,
			TypeArticle:  results.Article.TypeArticle,
			Composition:  results.Article.Composition,
			Presentation: results.Article.Presentation,
			CreatedAt:    results.Article.CreatedAt,
			UpdatedAt:    results.Article.UpdatedAt,
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

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	var rol util.Rol
	rol.Acl = "article"
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

func AddArticle(c *gin.Context) {
	var rol util.Rol
	rol.Acl = "article"
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
			t     model.Article
		)

		if c.BindJSON(&t) == nil {
			// cloudant.DB(dbName).Post(t)
			if t.Name == "" && t.TypeArticle == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"articles"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "articles": t})

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

func PutArticle(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	rol.Acl = "article"
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
			t     model.Article
		)

		if c.BindJSON(&t) == nil {

			if t.Name == "" && t.TypeArticle == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"articles"}
				cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "articles": t}, rev)
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

func BulkArticle(c *gin.Context) {
	// var rol util.Rol
	// rol.Acl = "test"
	// var datas util.Response

	// if util.IsCreate(c, rol).Success == true {

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
		t     model.Article
	)

	if c.BindJSON(&t) == nil {
		// cloudant.DB(dbName).Post(t)
		if t.Name == "" && t.TypeArticle == "" {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			var arrKey = []string{"articles"}
			cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "articles": t})

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
	// } else {
	// 	datas = util.IsCreate(c, rol)

	// 	c.JSON(200, datas)
	// }
}

func ArticleFaker(c *gin.Context) {
	data0, err := ioutil.ReadFile("FileSystem/articles.json")
	if err != nil {
		fmt.Println(err)
	}

	decode0 := []byte(string(data0))
	// var result0 []interface{}
	var result0 []model.Article
	json.Unmarshal(decode0, &result0)

	for i := range result0 {

		t := model.Article{
			Name:         result0[i].Name,
			TypeArticle:  result0[i].TypeArticle,
			Composition:  result0[i].Composition,
			Presentation: result0[i].Presentation,
		}
		strResponse0 := util.CurlPost(t, "/api/article/bulk")
		fmt.Println("##strResponse0: ", string(strResponse0))
	}

	datas := util.Response{
		true,
		"ok",
		result0,
	}
	c.JSON(200, datas)

}

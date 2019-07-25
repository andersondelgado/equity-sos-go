package countrys

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

func SelectCountrys(c *gin.Context) {

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"countrys", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.CountrysDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Countrys

		for i := range results.Doc {
			// fmt.Println("##doc: ", results[i].Doc)
			// fmt.Println("##i: ", results[i].Doc.Test)
			if results.Doc[i].Countrys.Short != "" {
				ts = append(ts, model.Countrys{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Short:     results.Doc[i].Countrys.Short,
					Country:   results.Doc[i].Countrys.Country,
					CreatedAt: results.Doc[i].Countrys.CreatedAt,
					UpdatedAt: results.Doc[i].Countrys.UpdatedAt,
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

func PaginateCountrys(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")
	var rol util.Rol
	rol.Acl = "test"
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
		var arrKey = []string{"countrys", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"countrys"}
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
		var results model.CountrysDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Countrys

		for i := range results.Doc {
			if results.Doc[i].Countrys.Short != "" {
				ts = append(ts, model.Countrys{
					IDs:       results.Doc[i].ID,
					Rev:       results.Doc[i].Rev,
					ID:        results.Doc[i].ID,
					Short:     results.Doc[i].Countrys.Short,
					Country:   results.Doc[i].Countrys.Country,
					CreatedAt: results.Doc[i].Countrys.CreatedAt,
					UpdatedAt: results.Doc[i].Countrys.UpdatedAt,
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

func EditCountrys(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"countrys", "_id", "_rev"}

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
		var result model.CountrysDocumentsArray
		json.Unmarshal(decode, &result)

		results := result.Doc[0]

		t := model.Countrys{
			IDs:       results.ID,
			Rev:       results.Rev,
			ID:        results.ID,
			Short:     results.Countrys.Short,
			Country:   results.Countrys.Country,
			CreatedAt: results.Countrys.CreatedAt,
			UpdatedAt: results.Countrys.UpdatedAt,
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

func DeleteCountrys(c *gin.Context) {
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

func AddCountrys(c *gin.Context) {
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
			t     model.Countrys
		)

		if c.BindJSON(&t) == nil {
			// cloudant.DB(dbName).Post(t)
			if t.Short == "" && t.Country == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"countrys"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "countrys": t})

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

func BulkCountrys(c *gin.Context) {
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
		t     model.Countrys
	)

	if c.BindJSON(&t) == nil {
		// cloudant.DB(dbName).Post(t)
		if t.Short == "" && t.Country == "" {
			datas = util.Response{
				false,
				"error_exception",
				nil,
			}
			c.JSON(200, datas)
		} else {
			var arrKey = []string{"countrys"}
			cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "countrys": t})

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

func PutCountrys(c *gin.Context) {
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
			t     model.Countrys
		)

		if c.BindJSON(&t) == nil {

			if t.Short == "" {
				datas = util.Response{
					false,
					"error_exception",
					nil,
				}
				c.JSON(200, datas)
			} else {
				var arrKey = []string{"countrys"}
				cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "countrys": t}, rev)
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

func CountryFaker(c *gin.Context) {
	data0, err := ioutil.ReadFile("FileSystem/country.json")
	if err != nil {
		fmt.Println(err)
	}

	decode0 := []byte(string(data0))
	// var result0 []interface{}
	var result0 []model.Countrys
	json.Unmarshal(decode0, &result0)

	for i := range result0 {

		t := model.Countrys{
			Short:   result0[i].Short,
			Country: result0[i].Country,
		}
		strResponse0 := util.CurlPost(t, "/api/country/bulk")
		fmt.Println("##strResponse0: ", string(strResponse0))
	}

	datas := util.Response{
		true,
		"ok",
		result0,
	}
	c.JSON(200, datas)

}

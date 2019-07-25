package post

import (
	"encoding/json"
	"fmt"
	"log"
	rando "math/rand"
	"strconv"
	"time"

	"../../config"
	"../../model"

	"../../util"
	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func SelectPost(c *gin.Context) {

	var rol util.Rol
	rol.Acl = "test"
	var datas util.Response

	if util.IsRead(c, rol).Success == true {

		var arrKey = []string{"posts", "_id", "_rev"}

		query := model.QuerySelectorAll{
			Selector: map[string]interface{}{
				"meta": arrKey[0],
			},
			Fields: arrKey,
		}

		respText := util.FindDataAll(query)

		jsonToString := (respText)
		decode := []byte(jsonToString)
		var results model.PostsDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Posts

		for i := range results.Doc {

			ts = append(ts, model.Posts{
				IDs:          results.Doc[i].ID,
				Rev:          results.Doc[i].Rev,
				ID:           results.Doc[i].ID,
				PostsData:    results.Doc[i].Posts.PostsData,
				ArticlesData: results.Doc[i].Posts.ArticlesData,
				StatusPost:   results.Doc[i].Posts.StatusPost,
				CreatedAt:    results.Doc[i].Posts.CreatedAt,
				UpdatedAt:    results.Doc[i].Posts.UpdatedAt,
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
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func PaginatePost(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")
	var rol util.Rol
	rol.Acl = "post"
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
		var arrKey = []string{"posts", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"posts"}
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
		var results model.PostsDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Posts

		for i := range results.Doc {

			ts = append(ts, model.Posts{
				IDs:          results.Doc[i].ID,
				Rev:          results.Doc[i].Rev,
				ID:           results.Doc[i].ID,
				PostsData:    results.Doc[i].Posts.PostsData,
				ArticlesData: results.Doc[i].Posts.ArticlesData,
				StatusPost:   results.Doc[i].Posts.StatusPost,
				CreatedAt:    results.Doc[i].Posts.CreatedAt,
				UpdatedAt:    results.Doc[i].Posts.UpdatedAt,
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
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func SearchPaginatePost(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	var rol util.Rol
	rol.Acl = "post"
	var datas util.Response
	var t model.Posts

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

		var arrKey = []string{"posts", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"posts"}
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
					"posts.post_data.full_text": map[string]interface{}{
						"$gt": t.PostsData.FullText,
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
			var results model.PostsDocumentsArray
			json.Unmarshal(decode, &results)

			// jsonToString, _ := json.Marshal(result.Rows)
			// // fmt.Println("##jsonToString: ", string(jsonToString))

			// decode := []byte(jsonToString)
			// var results []model.ArticleDocument
			// json.Unmarshal(decode, &results)

			var ts []model.Posts

			for i := range results.Doc {
				ts = append(ts, model.Posts{
					IDs:          results.Doc[i].ID,
					Rev:          results.Doc[i].Rev,
					ID:           results.Doc[i].ID,
					PostsData:    results.Doc[i].Posts.PostsData,
					ArticlesData: results.Doc[i].Posts.ArticlesData,
					StatusPost:   results.Doc[i].Posts.StatusPost,
					CreatedAt:    results.Doc[i].Posts.CreatedAt,
					UpdatedAt:    results.Doc[i].Posts.UpdatedAt,
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

func PaginatePostRolBusiness(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")
	var rol util.Rol
	rol.Acl = "business_post"
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
		var arrKey = []string{"posts", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"posts"}
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
		var results model.PostsDocumentsArray
		json.Unmarshal(decode, &results)

		var ts []model.Posts

		for i := range results.Doc {

			ts = append(ts, model.Posts{
				IDs:          results.Doc[i].ID,
				Rev:          results.Doc[i].Rev,
				ID:           results.Doc[i].ID,
				PostsData:    results.Doc[i].Posts.PostsData,
				ArticlesData: results.Doc[i].Posts.ArticlesData,
				StatusPost:   results.Doc[i].Posts.StatusPost,
				CreatedAt:    results.Doc[i].Posts.CreatedAt,
				UpdatedAt:    results.Doc[i].Posts.UpdatedAt,
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
		datas = util.IsRead(c, rol)

		c.JSON(200, datas)
	}
}

func SearchPaginatePostBusiness(c *gin.Context) {
	skip := c.Param("skip")
	limit := c.Param("limit")

	var rol util.Rol
	rol.Acl = "business_post"
	var datas util.Response
	var t model.Posts

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

		var arrKey = []string{"posts", "_id", "_rev", "meta", "tag"}
		var arrTag = []string{"posts"}
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
					"posts.post_data.full_text": map[string]interface{}{
						"$gt": t.PostsData.FullText,
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
			var results model.PostsDocumentsArray
			json.Unmarshal(decode, &results)

			// jsonToString, _ := json.Marshal(result.Rows)
			// // fmt.Println("##jsonToString: ", string(jsonToString))

			// decode := []byte(jsonToString)
			// var results []model.ArticleDocument
			// json.Unmarshal(decode, &results)

			var ts []model.Posts

			for i := range results.Doc {
				ts = append(ts, model.Posts{
					IDs:          results.Doc[i].ID,
					Rev:          results.Doc[i].Rev,
					ID:           results.Doc[i].ID,
					PostsData:    results.Doc[i].Posts.PostsData,
					ArticlesData: results.Doc[i].Posts.ArticlesData,
					StatusPost:   results.Doc[i].Posts.StatusPost,
					CreatedAt:    results.Doc[i].Posts.CreatedAt,
					UpdatedAt:    results.Doc[i].Posts.UpdatedAt,
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

func AddPost(c *gin.Context) {
	var rol util.Rol
	rol.Acl = "post"
	rol.Acl = "add_post"
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
			datas   util.Response
			payload model.Posts
			pd      model.PostsData
			// t1      interface{}
			t   model.Posts
			st  model.StatusPost
			st1 []model.StatusPost
		)

		if c.BindJSON(&t) == nil {

			fmt.Println("##i ..t: ", t)
			// return
			if t.PostsData.FullText == "" {
				datas = util.Response{
					false,
					"error_exception_null",
					nil,
				}
				c.JSON(200, datas)
			} else {
				user := util.Auth(c)

				pd = model.PostsData{
					UserID:         user.ID,
					RangeDamage:    t.PostsData.RangeDamage,
					LatLong:        t.PostsData.LatLong,
					Atachments:     t.PostsData.Atachments,
					CategoryPostID: t.PostsData.CategoryPostID,
					Category:       t.PostsData.Category,
					PromedyPeoples: t.PostsData.PromedyPeoples,
					FullText:       t.PostsData.FullText,
				}

				date := time.Now()

				st = model.StatusPost{
					StatusID:  "pending",
					CreatedAt: date,
				}

				st1 = append(st1, st)

				payload = model.Posts{
					PostsData:    pd,
					ArticlesData: t.ArticlesData,
					StatusPost:   st1,
					CreatedAt:    date,
				}
				var arrKey = []string{"posts"}
				cloudant.DB(dbName).Post(map[string]interface{}{"meta": arrKey[0], "tag": arrKey, "posts": payload})

				// smartContract.Init(payload)

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
				"error_exception_data",
				nil,
			}
			c.JSON(200, datas)
		}
	} else {
		datas = util.IsCreate(c, rol)

		c.JSON(200, datas)
	}
}

func EditPost(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "post"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"posts", "_id", "_rev"}

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
		var result model.PostsDocumentsArray
		json.Unmarshal(decode, &result)
		results := result.Doc[0]

		// artData := results.Posts.ArticlesData
		var artData []model.ArticlesData
		for i := range results.Posts.ArticlesData {
			// var artData0 []model.ArticlesData
			// var artData0 []model.ArticlesData
			for j := range results.Posts.ArticlesData[i].ArticlesActors {

				var arrKey1 = []string{"users", "_id", "_rev"}

				id1 := results.Posts.ArticlesData[i].ArticlesActors[j].UserID

				query1 := model.QuerySelectorAll{
					Selector: map[string]interface{}{
						"meta": arrKey1[0],
						"_id":  id1,
					},
					Fields: arrKey1,
				}

				respText1 := util.FindDataAll(query1)

				jsonToString1 := (respText1)
				decode1 := []byte(jsonToString1)
				var result1 model.UserDocumentsArray
				json.Unmarshal(decode1, &result1)
				results1 := result1.Doc[0]

				results.Posts.ArticlesData[i].ArticlesActors[j].User = model.User{
					Avatar:   results1.User.Avatar,
					Username: results1.User.Username,
					Email:    results1.User.Email,
				}

				// artData0 = append(artData0, results.Posts.ArticlesData[i])

			}

			artData = append(artData, results.Posts.ArticlesData[i])
			// artData = append(artData, artData0)
		}

		t := model.Posts{
			IDs:       results.ID,
			Rev:       results.Rev,
			ID:        results.ID,
			PostsData: results.Posts.PostsData,
			// ArticlesData: results.Posts.ArticlesData,
			ArticlesData: artData,
			StatusPost:   results.Posts.StatusPost,
			CreatedAt:    results.Posts.CreatedAt,
			UpdatedAt:    results.Posts.UpdatedAt,
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

func PutPostBusiness(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	rol.Acl = "business_post"
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
			// t     []model.ArticlesData
			// t2 []model.ArticlesData
			t1 interface{}
			p  model.Posts
		)

		if c.BindJSON(&t1) == nil {

			fmt.Println("##i ..t1: ", t1)
			str, _ := json.Marshal(t1)
			decodex := []byte(str)
			var rst model.Posts
			json.Unmarshal(decodex, &rst)

			fmt.Println("### ..t: ", rst.ArticlesData)

			var arrKey = []string{"posts", "_id", "_rev"}

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
			var result model.PostsDocumentsArray
			json.Unmarshal(decode, &result)
			results := result.Doc[0]
			// var arAct []model.ArticlesActors
			date := time.Now()
			user := util.Auth(c)

			for i := range results.Posts.ArticlesData {

				var (
					SequencesID string
				)

				SequencesID = (RandStringBytes(10))
				var arAct []model.ArticlesActors

				fmt.Println("\nafter results.Posts.ArticlesData[i].ArticlesActors: ", results.Posts.ArticlesData[i].ArticlesActors)
				fmt.Println("\n")
				if len(results.Posts.ArticlesData[i].ArticlesActors) > 0 && results.Posts.ArticlesData[i].ArticlesActors != nil {
					for j := range rst.ArticlesData {
						for k := range rst.ArticlesData[j].ArticlesActors {
							if results.Posts.ArticlesData[i].ArticleID == rst.ArticlesData[j].ArticleID {
								if results.Posts.ArticlesData[i].QuantityLeft != 0 {
									arAct = append(results.Posts.ArticlesData[i].ArticlesActors, model.ArticlesActors{
										SequencesID:            SequencesID,
										UserID:                 user.ID,
										QuantityDelivery:       rst.ArticlesData[j].ArticlesActors[k].QuantityDelivery,
										QuantityLeft:           results.Posts.ArticlesData[i].QuantityLeft,
										StatusDeliverySender:   true,
										StatusDeliveryReciever: false,
										StatusOrder:            "pending_delivery_receiver",
										CreatedAt:              date,
										DateDeliverySender:     date,
									})
								} else {
									for L := range results.Posts.ArticlesData[i].ArticlesActors {
										arAct = append(results.Posts.ArticlesData[i].ArticlesActors, model.ArticlesActors{
											SequencesID:            results.Posts.ArticlesData[i].ArticlesActors[L].SequencesID,
											UserID:                 results.Posts.ArticlesData[i].ArticlesActors[L].UserID,
											QuantityDelivery:       results.Posts.ArticlesData[i].ArticlesActors[L].QuantityDelivery,
											QuantityLeft:           results.Posts.ArticlesData[i].ArticlesActors[L].QuantityLeft,
											StatusDeliverySender:   results.Posts.ArticlesData[i].ArticlesActors[L].StatusDeliverySender,
											StatusDeliveryReciever: results.Posts.ArticlesData[i].ArticlesActors[L].StatusDeliveryReciever,
											StatusOrder:            results.Posts.ArticlesData[i].ArticlesActors[L].StatusOrder,
											CreatedAt:              date,
											DateDeliverySender:     date,
										})
									}
								}
							}
						}
					}

				} else {
					for j := range rst.ArticlesData {
						for k := range rst.ArticlesData[j].ArticlesActors {
							fmt.Println("\n null ", rst.ArticlesData[j])
							if results.Posts.ArticlesData[i].ArticleID == rst.ArticlesData[j].ArticleID {
								fmt.Println("\n null rst.ArticlesData[j].ArticleID ", rst.ArticlesData[j].ArticleID)
								arAct = append(arAct, model.ArticlesActors{
									SequencesID:            SequencesID,
									UserID:                 user.ID,
									QuantityDelivery:       rst.ArticlesData[j].ArticlesActors[k].QuantityDelivery,
									QuantityLeft:           results.Posts.ArticlesData[i].QuantityLeft,
									StatusDeliverySender:   true,
									StatusDeliveryReciever: false,
									StatusOrder:            "pending_delivery_receiver",
									CreatedAt:              date,
									DateDeliverySender:     date,
								})
							}
						}
					}

				}

				fmt.Println("\nartActors: ", arAct)
				results.Posts.ArticlesData[i].ArticlesActors = arAct

				fmt.Println("\nbefore results.Posts.ArticlesData[i].ArticlesActors: ", results.Posts.ArticlesData[i].ArticlesActors)

				fmt.Println("\n af results.Posts.ArticlesData: ", results.Posts.ArticlesData)

			}

			fmt.Println("\naf results.Posts.ArticlesData: ", results.Posts.ArticlesData)

			status_pos := results.Posts.StatusPost
			for i := range status_pos {
				if status_pos[i].StatusID != "partial_completed" {
					status_pos = append(status_pos, model.StatusPost{
						StatusID:  "in_progress",
						CreatedAt: date,
					})
				}
			}

			// p = model.Posts{
			// 	// IDs:          results.ID,
			// 	// Rev:          results.Rev,
			// 	// ID:           results.ID,
			// 	PostsData:    results.Posts.PostsData,
			// 	ArticlesData: t2,
			// 	StatusPost:   status_pos,
			// 	CreatedAt:    results.Posts.CreatedAt,
			// 	UpdatedAt:    results.Posts.UpdatedAt,
			// }

			p = model.Posts{
				// IDs:          results.ID,
				// Rev:          results.Rev,
				// ID:           results.ID,
				PostsData:    results.Posts.PostsData,
				ArticlesData: results.Posts.ArticlesData,
				StatusPost:   status_pos,
				CreatedAt:    results.Posts.CreatedAt,
				UpdatedAt:    results.Posts.UpdatedAt,
			}

			var arrKeyPost = []string{"posts"}

			cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKeyPost[0], "tag": arrKeyPost, "posts": p}, rev)
			datas = util.Response{
				true,
				"ok",
				p,
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
		datas = util.IsUpdate(c, rol)

		c.JSON(200, datas)
	}
}

func EditPostBusiness(c *gin.Context) {
	id := c.Param("id")

	var rol util.Rol
	rol.Acl = "business_post"
	var datas util.Response

	if util.IsEdit(c, rol).Success == true {

		var arrKey = []string{"posts", "_id", "_rev"}

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
		var result model.PostsDocumentsArray
		json.Unmarshal(decode, &result)
		results := result.Doc[0]
		t := model.Posts{
			IDs:          results.ID,
			Rev:          results.Rev,
			ID:           results.ID,
			PostsData:    results.Posts.PostsData,
			ArticlesData: results.Posts.ArticlesData,
			StatusPost:   results.Posts.StatusPost,
			CreatedAt:    results.Posts.CreatedAt,
			UpdatedAt:    results.Posts.UpdatedAt,
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

func PutPost(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")
	var rol util.Rol
	rol.Acl = "add_post"
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
			t1    interface{}
		)

		if c.BindJSON(&t1) == nil {

			fmt.Println("##i ..t1: ", t1)
			str, _ := json.Marshal(t1)
			decodex := []byte(str)
			var rst model.Posts
			json.Unmarshal(decodex, &rst)

			fmt.Println("### ..t: ", rst.ArticlesData)

			var arrKey = []string{"posts", "_id", "_rev"}

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
			var result model.PostsDocumentsArray
			json.Unmarshal(decode, &result)
			results := result.Doc[0]
			var StatusID string
			date := time.Now()
			countArtCompleted := 0
			for i := range results.Posts.ArticlesData {

				arData := results.Posts.ArticlesData[i]
				artId := arData.ArticleID
				arAct := arData.ArticlesActors
				if len(arAct) > 0 && arAct != nil {

					for z := range arAct {
						artActUserId := arAct[z].UserID

						if len(rst.ArticlesData) > 0 && rst.ArticlesData != nil {
							for j := range rst.ArticlesData {

								rstArt := rst.ArticlesData[j]
								rstArtId := rstArt.ArticleID
								rstArAct := rstArt.ArticlesActors

								if artId == rstArtId {
									if len(rstArAct) > 0 && rstArAct != nil {
										for k := range rstArAct {
											rstartAct := rstArAct[k]
											rstartActUserId := rstartAct.UserID

											if rstartActUserId == artActUserId {

												quantity_ask := arData.QuantityAsk
												quantity_left := arData.QuantityLeft
												quantity_delivery := rstartAct.QuantityDelivery

												if rstartAct.StatusDeliveryReciever == true {
													if rstartAct.StatusOrder != "completed" {
														if quantity_left > 0 {
															if quantity_left >= quantity_delivery {
																quantity_left = (quantity_left - quantity_delivery)
															}

															fmt.Println("\n")
															fmt.Println("Print quantity_ask --- ", quantity_ask)
															fmt.Println("Print quantity_delivery --- ", quantity_delivery)
															fmt.Println("Print quantity_left 1 --- ", quantity_left)
															fmt.Println("Print quantity_left 2 After **--- ", quantity_left)
															fmt.Println("////////////////////////")

															results.Posts.ArticlesData[i].QuantityLeft = quantity_left
															results.Posts.ArticlesData[i].ArticlesActors[z].QuantityLeft = quantity_left
															results.Posts.ArticlesData[i].ArticlesActors[z].StatusDeliveryReciever = rstartAct.StatusDeliveryReciever
															results.Posts.ArticlesData[i].ArticlesActors[z].DateDeliveryReciever = date
															results.Posts.ArticlesData[i].ArticlesActors[z].StatusOrder = "completed"

															// if results.Posts.ArticlesData[i].ArticlesActors[z].StatusOrder == "completed"{
															// 	results.Posts.ArticlesData[i].ArticlesActors[z].StatusOrder = "finalize"
															// }

														}
													} else if results.Posts.ArticlesData[i].ArticlesActors[z].StatusDeliveryReciever == false && quantity_left == 0 {
														results.Posts.ArticlesData[i].ArticlesActors[z].StatusOrder = "finalize"
														// }
														// results.Posts.ArticlesData[i].ArticlesActors[z].StatusOrder = "finalize"
													}
												} else {
													if rstartAct.StatusOrder == "reverse" {
														if quantity_left <= quantity_ask {
															quantity_left = (quantity_left + quantity_delivery)
														}
													}
												}
											}

											if results.Posts.ArticlesData[i].QuantityLeft != 0 {
												StatusID = "partial_completed"
												countArtCompleted++
											} else {
												StatusID = "completed"
											}
										}
									}
								}
							}
						}

					}

				}
			}

			status_pos := results.Posts.StatusPost
			status_pos = append(status_pos, model.StatusPost{
				StatusID:  StatusID,
				CreatedAt: date,
			})

			p := model.Posts{
				// IDs:          results.ID,
				// Rev:          results.Rev,
				// ID:           results.ID,
				PostsData:    results.Posts.PostsData,
				ArticlesData: results.Posts.ArticlesData,
				StatusPost:   status_pos,
				CreatedAt:    results.Posts.CreatedAt,
				UpdatedAt:    results.Posts.UpdatedAt,
			}

			var arrKeyPost = []string{"posts"}

			cloudant.DB(dbName).Put(id, map[string]interface{}{"meta": arrKeyPost[0], "tag": arrKeyPost, "posts": p}, rev)

			datas = util.Response{
				true,
				"ok",
				p,
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
		datas = util.IsUpdate(c, rol)

		c.JSON(200, datas)
	}
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	rev := c.Param("rev")

	var rol util.Rol
	rol.Acl = "post"
	rol.Rol = "admin"
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

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rando.Intn(len(letterBytes))]
	}
	return string(b)
}

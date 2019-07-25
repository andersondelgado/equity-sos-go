package main

import (
	"fmt"
	"net/http"
	"os"
	"./controller/article"
	"./controller/category"
	"./controller/countrys"
	"./controller/post"
	"./controller/security"
	"./controller/test"
	"./middlewares"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/size"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

const (
	DEFAULT_PORT = "8000"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	var size int64
	size=(200 * 1024 * 1024)
	router.Use(limits.RequestSizeLimiter(size))
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{"HEAD", "POST", "GET", "PUT", "DELETE", "PATCH"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	// AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return origin == "https://github.com"
	// 	// },
	// 	// MaxAge: 12 * time.Hour,
	// }))
	// router.Use(cors.Default())
	router.Use(CORSMiddleware())

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	api.GET("/jokes", test.JokeHandler)
	//
	api.GET("/testdb/all", test.SelectDBTest)
	//
	api.GET("/setup-permission", security.SetupPermission)

	api.POST("/setup-admin", security.SetupAdmin)
	api.POST("/setup-user", security.SetupUser)

	api.POST("/87780FA5DE684E87CB92B279F0BC07B14F572851E73B8943A097C1770A5F38E6", security.Register)
	// User Actions
	api.POST("/register", security.SetupUser)
	// api.POST("/register-user", security.AddUser)
	api.POST("/edit-user", security.EditUser)
	api.POST("/login", security.Login)
	api.GET("/permission/faker", security.PermissionFaker)
	api.GET("/permission/delete-faker", security.PermissionDeleteFaker)
	api.GET("/permission/all", security.PermissionAll)
	// roles
	api.POST("/roles/assign", security.AssignRoles)
	api.GET("/roles/edit/:id", security.EditRoles)
	api.PUT("/roles/update/:id/:rev", security.UpdateRoles)
	api.GET("/roles/delete/:id/:rev", security.DeleteRoles)
	api.GET("/roles/by-user/:user_id", security.GetRolByUser)
	//

	r := api.Group("/")
	r.Use(middlewares.AuthJWT())

	r.GET("/menu", security.Menu)
	r.GET("/user-info", security.InfoUser)
	// test
	r.POST("/test/add", test.AddTest)
	r.GET("/test/all", test.SelectTest)
	r.GET("/test/paginate/:skip/:limit", test.PaginateTest)
	r.POST("/test/search-paginate/:skip/:limit", test.SearchPaginateTest)
	r.GET("/test/edit/:id", test.EditTest)
	r.PUT("/test/update/:id/:rev", test.PutTest)
	r.GET("/test/delete/:id/:rev", test.DeleteTest)
	// article
	api.GET("/article/faker", article.ArticleFaker)
	api.POST("/article/bulk", article.BulkArticle)
	r.POST("/article/add", article.AddArticle)
	r.GET("/article/all", article.SelectArticle)
	r.GET("/article/paginate/:skip/:limit", article.PaginateArticle)
	r.POST("/article/search-paginate/:skip/:limit", article.SearchPaginateArticle)
	r.GET("/article/edit/:id", article.EditArticle)
	r.PUT("/article/update/:id/:rev", article.PutArticle)
	r.GET("/article/delete/:id/:rev", article.DeleteArticle)
	// category
	r.POST("/category/add", category.AddCategory)
	r.GET("/category/all", category.SelectCategory)
	r.GET("/category/paginate/:skip/:limit", category.PaginateCategory)
	r.POST("/category/search-paginate/:skip/:limit", category.SearchPaginateCategory)
	r.GET("/category/edit/:id", category.EditCategory)
	r.PUT("/category/update/:id/:rev", category.PutCategory)
	r.GET("/category/delete/:id/:rev", category.DeleteCategory)
	// country
	api.GET("/country/faker", countrys.CountryFaker)
	api.POST("/country/bulk", countrys.BulkCountrys)
	r.POST("/country/add", countrys.AddCountrys)
	r.GET("/country/all", countrys.SelectCountrys)
	r.GET("/country/paginate/:skip/:limit", countrys.PaginateCountrys)
	r.GET("/country/edit/:id", countrys.EditCountrys)
	r.PUT("/country/update/:id/:rev", countrys.PutCountrys)
	r.GET("/country/delete/:id/:rev", countrys.DeleteCountrys)
	// post
	r.POST("/post/add", post.AddPost)
	r.GET("/post/all", post.SelectPost)
	r.GET("/post/paginate/:skip/:limit", post.PaginatePost)
	r.GET("/post/search-paginate/:skip/:limit", post.SearchPaginatePost)
	r.GET("/post/business/paginate/:skip/:limit", post.PaginatePostRolBusiness)
	r.GET("/post/business/search-paginate/:skip/:limit", post.SearchPaginatePostBusiness)
	r.GET("/post/business/edit/:id", post.EditPostBusiness)
	r.GET("/post/edit/:id", post.EditPost)
	r.PUT("/post/business/update/:id/:rev", post.PutPostBusiness)
	r.PUT("/post/update/:id/:rev", post.PutPost)
	r.GET("/post/delete/:id/:rev", post.DeletePost)
	// register users
	r.POST("/users/add", security.AddUsers)
	r.GET("/users/paginate/:skip/:limit", security.PaginateUsers)
	r.POST("/users/search-paginate/:skip/:limit", security.SearchPaginateUsers)
	r.PUT("/users/update/:id/:rev", security.PutUsers)
	r.GET("/users/edit/:id", security.EditUsers)
	r.GET("/users/delete/:id/:rev", security.DeleteUsers)
	// list roles
	api.GET("/permission/list", security.ListPermissions)
	// Start and run the server
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	router.Run(":"+port)
}

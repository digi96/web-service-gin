package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"example/web-service-gin/controllers"
	dbCon "example/web-service-gin/db/sqlc"
	"example/web-service-gin/routes"
	"example/web-service-gin/token"
	"example/web-service-gin/util"

	docs "example/web-service-gin/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	server *gin.Engine
	db     *dbCon.Queries
	ctx    context.Context

	ContactController   controllers.ContactController
	ContactRoutes       routes.ContactRoutes
	RideOrderController controllers.RideOrderController
	RideOrderRoutes     routes.RideOrderRoutes
	AuthController      controllers.AuthController
	AuthRoutes          routes.AuthRoutes
	tokenMaker          token.Maker
)

func init() {
	ctx = context.TODO()
	config, err := util.LoadConfig(".")

	tokenMakerTemp, err := token.NewPasetoMaker(config.SymmetricKey)

	tokenMaker = tokenMakerTemp

	if err != nil {
		log.Fatalf("Could not create token maker %v", err)
	}

	//ginMiddlewareHandler := middleware.AuthMiddleware(tokenMaker)

	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db = dbCon.New(conn)

	fmt.Println("PostgreSql connected successfully...")

	ContactController = *controllers.NewContactController(db, ctx, tokenMaker)
	ContactRoutes = routes.NewRouteContact(ContactController)

	RideOrderController = *controllers.NewRideOrderController(db, ctx)
	RideOrderRoutes = routes.NewRouteRideOrder(RideOrderController)

	AuthController = *controllers.NewAuthController(db, ctx, tokenMaker)
	AuthRoutes = routes.NewRouteAuth(AuthController)

	server = gin.Default()
}

// // album represents data about a record album.
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

func main() {
	//router := gin.Default()
	// router.GET("/albums", getAlbums)
	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	// router.Run("localhost:8080")

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	docs.SwaggerInfo.BasePath = "/api"

	routerGroup := server.Group("/api")

	server.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "The contact APi is working fine"})
	})

	ContactRoutes.ContactRoute(routerGroup)
	RideOrderRoutes.RideOrderRoutes(routerGroup)
	AuthRoutes.AuthRoute(routerGroup)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	// Web-Service-gin godoc
	// @Summary Login
	// @Schemes
	// @Description Login to get token.
	// @Tags Login
	// @Accept json
	// @Produce json
	// @Success 200
	// @Router /login [post]
	server.POST("/login/", func(c *gin.Context) {

		generatedToken, err := tokenMaker.CreateToken("testuser", time.Minute*30)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"staus": "failed", "message": fmt.Sprintf("Failed to get token maker, %s", err.Error())})
		}

		c.JSON(200, gin.H{"token": generatedToken})

	})

	url := ginSwagger.URL("doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Fatal(server.Run(":" + config.ServerAddress))
}

// // albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// // getAlbums responds with the list of all albums as JSON.
// func getAlbums(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, albums)
// }

// // postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// 	// Call BindJSON to bind the received JSON to
// 	// newAlbum.
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop over the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }

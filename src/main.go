package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Product type
type Product struct {
	gorm.Model
	Code  string `json:”code”`
	Price uint   `json:”price”`
}

//DB var
var db *gorm.DB

func main() {
	initDb()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	//index
	r.GET("/", func(context *gin.Context) {
		context.HTML(200, "static/index.tmpl", gin.H{
			"routes": r.Routes(),
		})
	})

	//routes as JSON
	r.GET("/routes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"routes": r.Routes(),
		})
	})

	//hello
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"greet": "hello, world!",
		})
	})

	//echo
	r.GET("/echo/:echo", func(c *gin.Context) {
		echo := c.Param("echo")
		c.JSON(http.StatusOK, gin.H{
			"echo": echo,
		})
	})

	//fake upload
	r.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"uploaded": len(files),
		})
	})

	//REST
	r.GET("/products", GetProducts)
	r.POST("/products", AddProduct)

	//go-go-go
	r.Run() // listen and serve on 0.0.0.0:8080
}

func initDb() {
	database, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	db = database

	defer database.Close()
	database.AutoMigrate(&Product{})
}

func GetProducts(c *gin.Context) {
	if db == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("Failed to connect to DB")
	}
	var products []Product

	if err := db.Find(&products).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, products)
		log.Println(products)
	}
}

func AddProduct(c *gin.Context) {
	if db == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("Failed to connect to DB")
	}

	product := Product{}
	err := c.BindJSON(&product)

	if err != nil {
		exception := err.Error()
		log.Println(err)
		c.JSON(400, gin.H{
			"exception": exception,
			"data":      product,
		})
	} else {
		status := db.Create(&product)
		if status.Error == nil && status.RowsAffected > 0 {
			c.JSON(201, gin.H{
				"data": product,
			})
		} else {
			c.AbortWithStatusJSON(500, gin.H{
				"exception": status.Error,
				"data":      product,
			})
		}
	}
}

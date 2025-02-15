package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book 图书模型
type Book struct {
	gorm.Model
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

// Order 订单模型
type Order struct {
	gorm.Model
	BookID   uint `json:"book_id"`
	Quantity int  `json:"quantity"`
}

var db *gorm.DB

func main() {
	// 初始化数据库
	var err error
	db, err = gorm.Open(sqlite.Open("bookshop.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据库结构
	db.AutoMigrate(&Book{}, &Order{})

	// 创建Gin路由
	r := gin.Default()
	// gin.Default() 是 Gin 框架提供的一个默认配置，它包含了一些常用的中间件，如日志记录、请求解析、响应渲染等。
	// 打印日志
	r.Use(gin.Logger())
	// 恢复，用于在发生 panic 时恢复应用程序的运行。 panic 是 Go 语言中的一种错误处理机制，当发生不可恢复的错误时，程序会抛出 panic 并终止运行。
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return "请求信息：" + param.Method + " " + param.Path + "\n"
	})) // 打印请求信息到控制台，可以根据需要自定义日志格式

	// 配置CORS
	r.Use(cors.Default())

	// 管理员API
	admin := r.Group("/admin")
	{
		// 增加库存
		admin.POST("/stock/add", func(c *gin.Context) {
			var book Book
			if err := c.ShouldBindJSON(&book); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// 检查图书是否存在
			var existingBook Book
			result := db.First(&existingBook, book.ID)
			if result.Error != nil {
				// 如果图书不存在，创建新图书
				db.Create(&book)
			} else {
				// 如果图书存在，更新库存
				existingBook.Stock += book.Stock
				db.Save(&existingBook)
			}

			c.JSON(200, gin.H{"message": "库存更新成功"})
		})

		// 减少库存
		admin.POST("/stock/reduce", func(c *gin.Context) {
			var book Book
			if err := c.ShouldBindJSON(&book); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			var existingBook Book
			result := db.First(&existingBook, book.ID)
			if result.Error != nil {
				c.JSON(404, gin.H{"error": "图书不存在"})
				return
			}

			if existingBook.Stock < book.Stock {
				c.JSON(400, gin.H{"error": "库存不足"})
				// gin.H 用于构建 JSON 响应，它是一个 map[string]any 类型的别名。
				// 400 意思是请求错误，通常是客户端发送的请求不符合服务器的要求。
				return
			}

			existingBook.Stock -= book.Stock
			db.Save(&existingBook)

			c.JSON(200, gin.H{"message": "库存更新成功"})
		})
	}

	// 客户API
	// 获取所有图书
	r.GET("/books", func(c *gin.Context) {
		var books []Book
		db.Find(&books)
		c.JSON(200, books)
	})

	// 创建订单
	r.POST("/order", func(c *gin.Context) {
		var order Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 检查库存
		var book Book
		result := db.First(&book, order.BookID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "图书不存在"})
			return
		}

		if book.Stock < order.Quantity {
			c.JSON(400, gin.H{"error": "库存不足"})
			return
		}

		// 更新库存
		book.Stock -= order.Quantity
		db.Save(&book)

		// 创建订单
		db.Create(&order)

		c.JSON(200, gin.H{"message": "订单创建成功"})
	})

	// 启动服务器
	r.Run(":8080")
}

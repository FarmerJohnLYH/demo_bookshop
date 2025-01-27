package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// User 用户模型
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

var db *gorm.DB
var jwtKey = []byte("your-secret-key") // 在实际应用中应该使用环境变量

func main() {
	// 初始化数据库
	var err error
	db, err = gorm.Open(sqlite.Open("bookshop.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据库结构
	db.AutoMigrate(&Book{}, &Order{}, &User{})

	// 创建日志文件
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("无法创建日志文件: %v", err))
	}

	// 创建Gin路由
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(logFile))
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s | %s | %s | %d | %v\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	// 配置CORS
	r.Use(cors.Default())

	// 初始化root账号
	var rootUser User
	result := db.Where("username = ?", "root").First(&rootUser)
	if result.Error != nil {
		// 如果root用户不存在，创建一个新的root用户
		rootUser = User{
			Username: "root",
			Password: "root",
		}
		db.Create(&rootUser)
	}

	// 注册API
	r.POST("/api/register", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(400, gin.H{"error": "无效的请求数据"})
			return
		}
		// 检查用户名是否已存在
		var existingUser User
		result := db.Where("username =?", newUser.Username).First(&existingUser)
		if result.Error == nil {
			c.JSON(400, gin.H{"error": "用户名已存在"})
			return
		}
		// 创建新用户
		db.Create(&newUser)
		c.JSON(201, gin.H{"message": "注册成功"})
	})

	// 登录API
	r.POST("/api/login", func(c *gin.Context) {
		var loginUser User
		if err := c.ShouldBindJSON(&loginUser); err != nil {
			// ShouldBindJSON 用于解析请求体中的JSON数据到 loginUser 结构体中
			c.JSON(400, gin.H{"error": "无效的请求数据"})
			return
		}

		var user User
		result := db.Where("username = ?", loginUser.Username).First(&user)
		if result.Error != nil || user.Password != loginUser.Password {
			c.JSON(401, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 创建JWT令牌
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "生成令牌失败"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	})

	// 管理员API
	admin := r.Group("/admin")
	{
		// 重置图书数据
		admin.POST("/reset-books", func(c *gin.Context) {
			// 删除所有现有图书
			db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&Book{})

			// 预设的图书数据
			initialBooks := []Book{
				{Title: "百年孤独", Author: "加西亚·马尔克斯", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
				{Title: "鲁迅全集", Author: "鲁迅", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
				{Title: "毛泽东选集", Author: "毛泽东", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
				{Title: "白夜行", Author: "东野圭吾", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
				{Title: "Norwegian Wood", Author: "村上春树", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
				{Title: "The Old Man and the Sea", Author: "海明威", Price: float64(30+rand.Intn(121)) + float64(rand.Intn(100))/100, Stock: 50 + rand.Intn(151)},
			}

			// 创建新的图书记录
			for _, book := range initialBooks {
				db.Create(&book)
			}

			c.JSON(200, gin.H{"message": "图书数据已重置"})
		})

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

		// 使用事务来确保数据一致性
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 在事务中使用悲观锁检查库存
		var book Book
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, order.BookID)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(404, gin.H{"error": "图书不存在"})
			return
		}

		// 确保库存充足
		if book.Stock < order.Quantity {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "库存不足"})
			return
		}

		// 更新库存
		book.Stock -= order.Quantity
		if err := tx.Save(&book).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "更新库存失败"})
			return
		}

		// 创建订单
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "创建订单失败"})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "提交事务失败"})
			return
		}

		c.JSON(200, gin.H{"message": "订单创建成功"})
	})

	// 启动服务器
	r.Run(":8080")
}

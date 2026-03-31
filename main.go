package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type UserAuth struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

var db *gorm.DB

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/todo_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	db.AutoMigrate(&User{}, &UserAuth{})

	r := gin.Default()

	r.POST("/register", register)
	r.POST("/login", login)

	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.POST("/user", createUser)
		auth.GET("/users", getUsers)
		auth.GET("/user/:id", getUserByID)
		auth.PUT("/user/:id", updateUser)
		auth.DELETE("/user/:id", deleteUser)
	}

	r.Run(":8080")
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Name == "" {
		c.JSON(400, gin.H{"error": "姓名不能为空"})
		return
	}
	if user.Age < 0 || user.Age > 150 {
		c.JSON(400, gin.H{"error": "年龄必须在0-150之间"})
		return
	}

	db.Create(&user)
	c.JSON(200, user)
}

func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(200, users)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(200, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}

	var updates User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Model(&user).Updates(updates)
	c.JSON(200, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}

	db.Delete(&user)
	c.JSON(200, gin.H{"message": "删除成功"})
}

func register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var exist UserAuth
	if db.Where("username = ?", input.Username).First(&exist).Error == nil {
		c.JSON(400, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	user := UserAuth{
		Username: input.Username,
		Password: string(hashedPassword),
	}
	db.Create(&user)

	c.JSON(200, gin.H{"message": "注册成功", "user_id": user.ID})
}

func login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user UserAuth
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("your-secret-key"))

	c.JSON(200, gin.H{"token": tokenString, "user_id": user.ID})
}

// JWT认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "请提供认证token"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}

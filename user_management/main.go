package main

import (
	"errors"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func getConnectionString() *string {
	connectionString :=
		"host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
            " port=" + os.Getenv("DB_PORT") +
			" dbname=" + os.Getenv("DB_NAME") +
			" sslmode=disable" +
			" password=" + os.Getenv("DB_PASSWORD")

	return &connectionString
}

func initDB() {
	var err error

	// connect to db
	db, err = gorm.Open("postgres", *getConnectionString())
	if err != nil {
		panic("Failed to connect to database")
	}

	// make migration user model
	db.AutoMigrate(&User{})
}

func main() {
    godotenv.Load()

	// make db
	initDB()

	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)

	router.Run(":8080")
}

func validateUserInput(user *User) error {
	if len(user.Username) < 3 || len(user.Username) > 25 {
		return errors.New("Username must be between 3 and 25 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(user.Email) == false {
		return errors.New("Invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("Password must be equal to or longer than 6 characters")
	}

	return nil
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validateUserInput(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash user password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully registered user"})
}

func validateLoginInput(username *string, password *string) error {
	if len(*username) == 0 || len(*password) == 0 {
		return errors.New("username and password are required")
	}
	return nil
}

func login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	if err := validateLoginInput(&input.Username, &input.Password); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully logged in"})
}

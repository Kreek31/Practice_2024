package handlers

import (
	"context"
	"myproject/databases"
	"myproject/middlewares"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	database *databases.MongoDB
}

func CreateUserHandler(database *databases.MongoDB) *UserHandler {
	return &UserHandler{database: database}
}

// GetUserByUsername - вспомогательная функция для поиска пользователя по имени
func (handler *UserHandler) GetUserByUsername(username string) (*models.User, error) {
	collection := handler.database.Collection("users")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login Выполняет вход в аккаунт пользоваетля по username и password
// @Summary Выполняет вход в аккаунт пользоваетля
// @Description Выполняет вход в аккаунт пользоваетля по username и password
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param credentials body models.User true "username и password пользователя"
// @Success 200 {object} map[string]string "token"
// @Failure 401 {object} map[string]string "error"
// @Router /login [post]
func (handler *UserHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := handler.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Генерация JWT
	token, err := middlewares.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Регистрирует пользователя
// @Description Регистрирует нового пользователя
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param user body models.User true "New user data"
// @Success 200 {object} map[string]string "status"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /register [post]
func (handler *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := handler.database.Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user registered"})
}

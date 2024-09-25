package handlers

import (
	"context"
	"myproject/databases"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PetHandler struct {
	db *databases.MongoDB
}

func NewPetHandler(db *databases.MongoDB) *PetHandler {
	return &PetHandler{db: db}
}

// GetPet получает информацию о домашнем животном по ID
// @Summary Получение домашнего животного
// @Description Возвращает информацию о домашнем животном по ID
// @Tags Домашние животные
// @Accept json
// @Produce json
// @Param id path string true "ID домашнего животного"
// @Success 200 {object} models.Pet
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /pets/{id} [get]
func (h *PetHandler) GetPet(c *gin.Context) {
	id := c.Param("id")
	collection := h.db.Collection("pets")

	var pet models.Pet
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&pet)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pet"})
		return
	}

	c.JSON(http.StatusOK, pet)
}

// @Summary Создать новое домажнее животное
// @Description создает новое домашнее животное в системе
// @Tags Домашние животные
// @Accept  json
// @Produce  json
// @Param pet body models.Pet true "Информация о питомце"
// @Success 200 {object} map[string]string "status"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /pets [post]
func (h *PetHandler) CreatePet(c *gin.Context) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.Collection("pets")
	_, err := collection.InsertOne(context.TODO(), pet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create pet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "pet created"})
}

// GetPets получает список домашних животных по заданным параметрам
// @Summary Получение списка домашних животных
// @Description Возвращает список домашних животных по заданным параметрам фильтрации
// @Tags Домашние животные
// @Accept json
// @Produce json
// @Param petid query int false "ID домашнего животного"
// @Param name query string false "Имя домашнего животного"
// @Param age query int false "Возраст"
// @Param gender query string false "Пол"
// @Param species query string false "Вид домашнего животного"
// @Param breed query string false "Порода"
// @Success 200 {array} models.Pet
// @Failure 500 {object} map[string]string "error"
// @Router /pets [get]
func (h *PetHandler) GetPets(c *gin.Context) {
	collection := h.db.Collection("pets")

	// Построение фильтра на основе параметров запроса
	filter := bson.M{}

	// Получаем параметры запроса (они могут быть пустыми)
	petid := c.Query("id")
	name := c.Query("name")
	age := c.Query("age")
	gender := c.Query("gender")
	species := c.Query("species")
	breed := c.Query("breed")

	if petid != "" {
		filter["id"] = petid
	}

	if name != "" {
		filter["name"] = name
	}

	if age != "" {
		filter["age"] = age
	}

	if gender != "" {
		filter["gender"] = gender
	}

	if species != "" {
		filter["species"] = species
	}

	if breed != "" {
		filter["breed"] = breed
	}

	// Выполняем поиск в базе данных
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pets"})
		return
	}
	defer cursor.Close(context.TODO())

	// Создаем слайс для хранения пользователей
	var pets []models.Pet
	for cursor.Next(context.TODO()) {
		var pet models.Pet
		if err := cursor.Decode(&pet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode pet"})
			return
		}
		pets = append(pets, pet)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	// Возвращаем список пользователей
	c.JSON(http.StatusOK, pets)
}

// UpdatePet обновляет данные домашнего животного
// @Summary Обновление данных домашнего животного
// @Description Обновляет данные домашнего животного по ID
// @Tags Домашние животные
// @Accept json
// @Produce json
// @Param id path string true "ID домашнего животного"
// @Param pet body models.Pet true "Новые данные домашнего животного"
// @Success 200 {object} map[string]string "status"
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /pets/{id} [put]
func (h *PetHandler) UpdatePet(c *gin.Context) {
	id := c.Param("id")

	// Преобразование id из строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Фильтр для поиска домашнего животного по ID
	filter := bson.M{"_id": objectID}

	// Данные для обновления
	update := bson.M{
		"$set": bson.M{
			"name":    pet.Name,
			"age":     pet.Age,
			"gender":  pet.Gender,
			"species": pet.Species,
			"breed":   pet.Breed,
		},
	}

	collection := h.db.Collection("pets")
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Проверяем, был ли найден и обновлен пользователь
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "pet updated"})
}

// DeletePet удаляет домашнее животное по ID
// @Summary Удаление домашнего животного
// @Description Удаляет домашнее животное по ID
// @Tags Домашние животные
// @Accept json
// @Produce json
// @Param id path string true "ID домашнего животного"
// @Success 200 {object} map[string]string "status"
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /pets/{id} [delete]
func (h *PetHandler) DeletePet(c *gin.Context) {
	id := c.Param("id")

	// Преобразование строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	// Фильтр для удаления по _id
	filter := bson.M{"_id": objectID}

	collection := h.db.Collection("pets")
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pet"})
		return
	}

	// Если домашнее животное не найдено
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "pet deleted"})
}

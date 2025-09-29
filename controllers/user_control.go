package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"ginapi/models"
)

type UserController struct {
	Collection *mongo.Collection
}

// CreateUser godoc
// @Summary      Create user
// @Description  Buat user baru
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Data"
// @Security     BearerAuth
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /users [post]
// create a new user controller
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := uc.Collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Ambil semua user dari database
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
// get all users
func (uc *UserController) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := uc.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Ambil profil user berdasarkan token JWT
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      500  {object}  map[string]string
// @Router       /users/profile [get]
// get user profile from JWT token
func (uc *UserController) GetProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := uc.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

//	GetUserAddresses godoc
//
// @Summary Get addresses for a user
// @Description Retrieve all addresses associated with a user
// @Tags users
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} models.Address
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{userID}/addresses [get]
func (uc *UserController) GetUserAddresses(c *gin.Context) {
	email := c.Param("email")
	addressCollection := uc.Collection.Database().Collection("address")
	var addresses []models.Address

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := addressCollection.Find(ctx, bson.M{"user_id.$email": email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var address models.Address
		if err := cursor.Decode(&address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode address"})
			return
		}
		addresses = append(addresses, address)
	}

	c.JSON(http.StatusOK, addresses)
}

// AddUserAddresses godoc
// @Summary Create a new address
// @Description Create a new address for a user
// @Tags addresses
// @Accept json
// @Produce json
// @Param address body models.Address true "Address to create"
// @Success 201 {object} models.Address
// @Failure 400 {object} string "Internal Server Error"
// @Failure 500 {object} string "Internal Server Error"
// @Router /addresses [post]
func (uc *UserController) AddUserAddresses(c *gin.Context) {
	var body struct {
		Street      string  `json:"street" binding:"required"`
		ShipToName  string  `json:"ship_to_name" binding:"required"`
		PhoneNumber string  `json:"phone_number" binding:"required"`
		Longitude   float64 `json:"longitude" binding:"required"`
		Latitude    float64 `json:"latitude" binding:"required"`
		Notes       string  `json:"notes"`
		IsDefault   bool    `json:"is_default"`
		RegionID    string  `json:"region_id" binding:"required"`
		UserID      string  `json:"user_id" binding:"required"`
		Email       string  `json:"email" binding:"required"`
	}

	email := c.Param("email")

	regionObjID, err := primitive.ObjectIDFromHex(body.RegionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	address := models.Address{
		ID:          primitive.NewObjectID(),
		Sreet:       body.Street,
		ShipToName:  body.ShipToName,
		PhoneNumber: body.PhoneNumber,
		Longitude:   body.Longitude,
		Latitude:    body.Latitude,
		Notes:       body.Notes,
		IsDefault:   body.IsDefault,
		RegionID:    regionObjID,
		UserID:      oid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := uc.Collection.InsertOne(ctx, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": result})
}

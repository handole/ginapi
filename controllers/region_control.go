package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ginapi/models"
)

type RegionController struct {
	Collection *mongo.Collection
}

// CreateRegion godoc
// @Summary      Create region
// @Description  Buat region baru
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        region  body      models.Region  true  "Region Data"
// @Security     BearerAuth
// @Success      201     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Router       /regions [post]
// create a new region controller
func (rc *RegionController) CreateRegion(c *gin.Context) {
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.Collection.InsertOne(ctx, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create region"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Region created successfully"})
}

// GetRegions godoc
// @Summary      Get all regions
// @Description  Ambil semua region dari database
// @Tags         regions
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Param        state       query     string  false  "Filter by state"
// @Param        city        query     string  false  "Filter by city"
// @Param        district    query     string  false  "Filter by district"
// @Param        code        query     string  false  "Filter by code"
// @Param        zipcode     query     string  false  "Filter by zipcode"
// @Param        sub_district query    string  false  "Filter by sub_district"
// @Security     BearerAuth
// @Success      200  {array}   models.Region
// @Failure      500  {object}  map[string]string
// @Router       /regions [get]
func (rc *RegionController) GetRegions(c *gin.Context) {
	// default values
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// filter query
	state := c.Query("state")
	city := c.Query("city")
	district := c.Query("district")
	code := c.Query("code")
	zipcode := c.Query("zipcode")
	subDistrict := c.Query("sub_district")

	filter := bson.M{}
	if state != "" {
		filter["state"] = bson.M{"$regex": state, "$options": "i"}
	}
	if city != "" {
		filter["city"] = bson.M{"$regex": city, "$options": "i"}
	}
	if district != "" {
		filter["district"] = bson.M{"$regex": district, "$options": "i"}
	}
	if code != "" {
		filter["code"] = code
	}
	if zipcode != "" {
		filter["zipcode"] = zipcode
	}
	if subDistrict != "" {
		filter["sub_district"] = bson.M{"$regex": subDistrict, "$options": "i"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// take total count
	total, err := rc.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count regions"})
		return
	}

	// set pagination headers
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}) // sort by created_at desc

	cursor, err := rc.Collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch regions"})
		return
	}
	defer cursor.Close(ctx)

	var regions []models.Region
	if err = cursor.All(ctx, &regions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse regions"})
		return
	}

	// response with pagination info
	c.JSON(http.StatusOK, gin.H{
		"data":       regions,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetRegionByID godoc
// @Summary      Get region by ID
// @Description  Ambil region berdasarkan ID
// @Tags         regions
// @Produce      json
// @Param        id   path      string  true  "Region ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Region
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /regions/{id} [get]
// get a region by ID
func (rc *RegionController) GetRegionByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var region models.Region
	err = rc.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&region)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
		return
	}

	c.JSON(http.StatusOK, region)
}

// UpdateRegion godoc
// @Summary      Update region
// @Description  Update region berdasarkan ID
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        id      path      string        true  "Region ID"
// @Param        region  body      models.Region  true  "Region Data"
// @Security     BearerAuth
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /regions/{id} [put]
// update a region by ID
func (rc *RegionController) UpdateRegion(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"state":        region.State,
			"city":         region.City,
			"district":     region.District,
			"code":         region.Code,
			"zipcode":      region.Zipcode,
			"sub_district": region.SubDisctrict,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rc.Collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update region"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Region updated successfully"})
}

// DeleteRegion godoc
// @Summary      Delete region
// @Description  Hapus region berdasarkan ID
// @Tags         regions
// @Produce      json
// @Param        id   path      string  true  "Region ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /regions/{id} [delete]
// delete a region by ID
func (rc *RegionController) DeleteRegion(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := rc.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete region"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Region deleted successfully", "deletedCount": result.DeletedCount})
}

package controllers

import (
	"context"
	"fmt"
	"goauth/database"
	"goauth/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, "book")

var validateBook = validator.New()

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var book models.Book

		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validateBook.Struct(book)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()

		uid := c.GetString("uid")
		book.Id = primitive.NewObjectID()
		book.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		book.Created_by = uid
		resultInsectionNumber, insertErr := bookCollection.InsertOne(ctx, book)

		if insertErr != nil {
			msg := fmt.Sprintf("book item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsectionNumber)
	}
}

func GetBooks() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		recoredPerPage, err := strconv.Atoi(c.Query("recoredPerPage"))
		if err != nil || recoredPerPage < 1 {
			recoredPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recoredPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		uid := c.GetString("uid")            //{"created_by", uid},
		userType := c.GetString("user_type") //userType, "ADMIN"

		fmt.Println(userType == "ADMIN")

		var matchStage = bson.D{}

		if userType == "ADMIN" {
			matchStage = bson.D{{"$match", bson.D{{}}}}
		} else {
			matchStage = bson.D{{"$match", bson.D{{"created_by", uid}}}}
		}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recoredPerPage}}}}}}}
		result, err := bookCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})

		}

		var books []bson.M
		if err = result.All(ctx, &books); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, books[0])
	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		//bookId := c.Param("bookId")

		bookId, err := primitive.ObjectIDFromHex(c.Param("bookId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		var book models.Book

		err1 := bookCollection.FindOne(ctx, bson.M{"_id": bookId}).Decode((&book))

		defer cancel()

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		uid := c.GetString("uid")
		userType := c.GetString("user_type")
		if book.Created_by == uid || userType == "ADMIN" {
			c.JSON(http.StatusOK, book)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorize to access this resource"})
			return
		}

	}
}

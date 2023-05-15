package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	// configs to connect to gmail smtp server
	smtpHost = "smtp.gmail.com"
	smtpPort = 587

	origin = "localhost:3000"
)

var appkey string
var hostEmail string

var client *mongo.Client

func initDatabase() error {
	uri := os.Getenv("MONGODB_URI")
	fmt.Println(uri)
	if uri == "" {
		log.Println("URI not found in environmental variables")
		return errors.New("URI not found")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Error connecting to database")
		return errors.New("cannot connect to db")
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("failed to load .env file")
		return
	}

	appkey = os.Getenv("APPKEY")
	hostEmail = os.Getenv("HOSTEMAIL")

	err = initDatabase()
	if err != nil {
		return
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{fmt.Sprintf("http://%s", origin)},
		AllowMethods:  []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Set-Cookie"},
		ExposeHeaders: []string{"Content-Length", "Content-Type", "Set-Cookie", "Access-Control-Allow-Credentials", "Access-Control-Expose-Headers", "Access-Control-Allow-Origin", "set-cookie"},

		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/verify", verifyCode)
	r.GET("/verify/link", verifyLink)
	r.POST("/register", register)
	r.POST("/login", login)

	r.Run(":3001")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func verifyCode(ctx *gin.Context) {
	// read from query params
	email := ctx.Query("email")
	hash := ctx.Query("hash")
	code := ctx.Query("code")

	log.Println(email, hash, code)

	// verify if hash || code is the same as database
	var user User
	if hash == "" && code == "" {

	}

	coll := client.Database("auth").Collection("users")
	filter := bson.D{{"email", email}}
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(fmt.Sprintf("searching for users with email: %s yield no result found", email))
		}
		log.Println("failed to search for user")
	}

	// set user to verified
	if user.VerifyCode == code && user.ValidCode.After(time.Now()) {
		id, _ := primitive.ObjectIDFromHex(user.ID)
		filter = bson.D{{"_id", id}}
		update := bson.D{{"$set", bson.D{{"isVerified", true}}}}
		result, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println("failed to update:", err)
		}
		log.Println(result)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": false,
		"error":   "incorrect code",
	})
}

func verifyLink(ctx *gin.Context) {
	// read from query params
	email := ctx.Query("email")
	hash := ctx.Query("hash")
	code := ctx.Query("code")

	log.Println(email, hash, code)

	// verify if hash || code is the same as database
	var user User
	if hash == "" && code == "" {

	}

	coll := client.Database("auth").Collection("users")
	filter := bson.D{{"email", email}}
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(fmt.Sprintf("searching for users with email: %s yield no result found", email))
		}
		log.Println("failed to search for user")
	}

	// set user to verified
	if user.VerifyHash == hash && user.ValidHash.After(time.Now()) {
		id, _ := primitive.ObjectIDFromHex(user.ID)
		filter = bson.D{{"_id", id}}
		update := bson.D{{"$set", bson.D{{"isVerified", true}}}}
		result, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println("failed to update:", err)
		}
		log.Println(result)

		ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/verify")
		return
	}

	log.Println(fmt.Sprintf("%s has done incorrect verification", email))
	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/notfound")
}

func register(ctx *gin.Context) {
	var newUser User

	err := ctx.BindJSON(&newUser)
	if err != nil {
		log.Println("Error binding JSON")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	// ensure email does not exist in database
	// ensure email and password is valid
	// left as an exercise

	// generate hash from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 6)
	if err != nil {
		log.Println("error hashing password")
	}
	newUser.Password = string(hash)

	// generate verification link
	newUser.VerifyHash = generateVerifyHash()
	// generate verification code
	newUser.VerifyCode = generateVerifyCode()

	// add deadline for hash and code expiry
	// in the example here, we use 15mins for hash and 3 mins for code
	now := time.Now()
	newUser.DateCreated = now
	newUser.ValidHash = now.Add(15 * time.Minute)
	newUser.ValidCode = now.Add(3 * time.Minute)

	// add user to database
	coll := client.Database("auth").Collection("users")
	log.Println("add: ", newUser)
	result, err := coll.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Println("failed to add new user:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
		})
		return
	}

	fmt.Println(result)

	// send the verification email
	sendVerificationEmail(newUser)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func login(ctx *gin.Context) {
	var user User

	err := ctx.BindJSON(&user)
	if err != nil {
		log.Println("Error binding JSON")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	var dbUser User
	coll := client.Database("auth").Collection("users")
	filter := bson.D{{"email", user.Email}}
	err = coll.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(fmt.Sprintf("searching for users with email: %s yield no result found", dbUser.Email))
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Either email or password is incorrect",
			})
			return
		}
		log.Println("failed to search for user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Println("password is different")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Either email or password is incorrect",
		})
		return
	}

	if !user.IsVerified {
		ctx.JSON(http.StatusOK, gin.H{
			"success":    false,
			"isVerified": false,
			"error":      "user not verified",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"isVerified": false,
		"error":      "",
	})
}

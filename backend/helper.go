package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

type User struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	Password    string    `json:"password" bson:"passwordHash"`
	VerifyHash  string    `bson:"verifyhash"`
	VerifyCode  string    `bson:"verifycode"`
	IsVerified  bool      `json:"isVerified" bson:"isVerified"`
	DateCreated time.Time `bson:"dateCreated"`
	ValidHash   time.Time `bson:"validHash"`
	ValidCode   time.Time `bson:"validCode"`
}

func generateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(10000)

	paddedNumber := fmt.Sprintf("%04d", randomNumber)

	return paddedNumber
}

func generateVerifyHash() string {
	charLen := 50
	randomBytes := make([]byte, charLen)

	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("failed to generate random bytes")
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return randomString[:charLen]
}

func getVerifyLink(email string, verifyHash string, verifyCode string) string {
	return fmt.Sprintf("http://localhost:3001/verify/link?email=%s&hash=%s&code=%s", email, verifyHash, verifyCode)
}

func sendVerificationEmail(user User) {
	verifyLink := getVerifyLink(user.Email, user.VerifyHash, user.VerifyCode)
	to := []string{user.Email}
	subject := "Email verification"
	body := fmt.Sprintf("Hello, \r\n\r\nThis email contains the verification link and verification code, please the link below to verify your account.\r\n\r\nVerification link: %s\r\n\r\nAlternatively, u can also key in the verification code to verify.\r\n\r\nVerification code: %s", verifyLink, user.VerifyCode)

	message := []byte("To: " + user.Email + "\r\n" +
		"From: " + hostEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", hostEmail, appkey, smtpHost)

	err := smtp.SendMail("smtp.gmail.com:587", auth, hostEmail, to, message)
	if err != nil {
		log.Println("Failed to send email")
		log.Println(err)
		return
	}
	log.Println("Email sent successfully")
}

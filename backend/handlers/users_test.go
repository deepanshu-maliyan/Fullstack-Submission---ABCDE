package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ecommerce-backend/database"
	"ecommerce-backend/handlers"
	"ecommerce-backend/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gin-gonic/gin"
)

var _ = Describe("User Handlers", func() {
	var router *gin.Engine

	BeforeEach(func() {
		database.Connect()
		router = gin.Default()
		router.POST("/users", handlers.CreateUser)
		router.POST("/users/login", handlers.LoginUser)
	})

	Describe("Creating a user", func() {
		Context("with valid data", func() {
			It("should create a user successfully", func() {
				userData := map[string]string{
					"username": "testuser",
					"password": "testpassword",
				}
				jsonData, _ := json.Marshal(userData)

				req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
				
				var user models.User
				json.Unmarshal(w.Body.Bytes(), &user)
				Expect(user.Username).To(Equal("testuser"))
				Expect(user.Password).To(BeEmpty()) // Password should be hidden
			})
		})

		Context("with missing data", func() {
			It("should return an error", func() {
				userData := map[string]string{
					"username": "testuser",
				}
				jsonData, _ := json.Marshal(userData)

				req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("User login", func() {
		BeforeEach(func() {
			// Create a test user first
			userData := map[string]string{
				"username": "logintest",
				"password": "testpassword",
			}
			jsonData, _ := json.Marshal(userData)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		})

		Context("with valid credentials", func() {
			It("should login successfully and return a token", func() {
				loginData := map[string]string{
					"username": "logintest",
					"password": "testpassword",
				}
				jsonData, _ := json.Marshal(loginData)

				req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var response handlers.LoginResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response.Token).ToNot(BeEmpty())
				Expect(response.User.Username).To(Equal("logintest"))
			})
		})

		Context("with invalid credentials", func() {
			It("should return unauthorized", func() {
				loginData := map[string]string{
					"username": "logintest",
					"password": "wrongpassword",
				}
				jsonData, _ := json.Marshal(loginData)

				req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

type App struct {
	db  *gorm.DB
	mux *mux.Router
}

type Contact struct {
	Name  string
	Email string
	Phone string
}

type Seller struct {
	SellerId int     `gorm:"primaryKey"`
	Contact  Contact `gorm:"embedded"`
	Rating   int     `gorm:"default:0"`
}

type Category struct {
	CategoryId int `gorm:"primaryKey;autoIncrement"`
	Name       string
}

type Subcategory struct {
	SubcategoryId int `gorm:"primaryKey"`
	CategoryId    int
	Category      Category
	Name          string
}

type Post struct {
	gorm.Model
	SellerId      int
	Seller        Seller `json:"-"`
	CategoryId    int
	SubcategoryId int
	Category      Category    `json:"-"`
	Subcategory   Subcategory `json:"-"`
	Title         string      `gorm:"not null"`
	Description   string      `gorm:"not null"`
	IsHidden      bool        `gorm:"default:false"`
	IsFlagged     bool        `gorm:"default:false"`
	IsDeleted     bool        `gorm:"default:false"`
	HasImage      bool        `gorm:"default:false"`
}

type User struct {
	ID       int  `gorm:"primaryKey;autoIncrement"`
	IsSeller bool `gorm:"default:false"`
	SellerId int
	Seller   Seller
	Contact  Contact `gorm:"embedded"`
}

func (app *App) start() {
	db = app.db
	err := db.AutoMigrate(&Contact{}, &Category{}, &Subcategory{}, &User{}, &Seller{}, &Post{})
	if err != nil {
		return
	}
	result := db.Exec("PRAGMA foreign_keys = ON", nil)
	if result.Error != nil {
		print(result.Error)
		return
	}
	setupEndpoints(app.mux)
	if isCategoryNotExists(db) {
		createDefaultValues(db)
	}
	handler := cors.AllowAll().Handler(app.mux)
	log.Fatal(http.ListenAndServe(":8081", handler))
}

func setupEndpoints(request *mux.Router) {
	request.HandleFunc("/post", createPost).Methods("POST", "OPTIONS")
	request.HandleFunc("/post", getAllPosts).Methods("GET", "OPTIONS")
	request.HandleFunc("/post/{id}", getPostById).Methods("GET")
	request.HandleFunc("/post/{id}", updatePost).Methods("PUT")
	request.HandleFunc("/post/{id}", deletePost).Methods("DELETE")
	request.HandleFunc("/categories", getAllCategories).Methods("GET")
	request.HandleFunc("/subcategories/{id}", getSubcategories).Methods("GET")
	//request.HandleFunc("/", getAllPosts).Methods("GET")

	//request.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist/Tradelist/")))
}

func isCategoryNotExists(db *gorm.DB) bool {
	var categories []Category

	result := db.Find(&categories)
	fmt.Println("Categories rows:", result.RowsAffected)
	if result.RowsAffected > 0 {
		return false
	} else {
		return true
	}
}

func createDefaultValues(db *gorm.DB) {
	var categories = []Category{{Name: "Jobs"}, {Name: "Property"}, {Name: "For Sale"}}
	result := db.Create(&categories)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	var subcategories = []Subcategory{
		{CategoryId: 1, Name: "Accounting"},
		{CategoryId: 1, Name: "HR"},
		{CategoryId: 1, Name: "Legal"},
		{CategoryId: 1, Name: "Customer Service"},
		{CategoryId: 1, Name: "Healthcare"},
		{CategoryId: 1, Name: "Hospitality"},
		{CategoryId: 1, Name: "Housekeeping"},
		{CategoryId: 1, Name: "Software"},
		{CategoryId: 1, Name: "Accounting"},
		{CategoryId: 2, Name: "For Sale"},
		{CategoryId: 2, Name: "To Rent"},
		{CategoryId: 2, Name: "To Share"},
		{CategoryId: 2, Name: "Sublet"},
		{CategoryId: 2, Name: "Storage"},
		{CategoryId: 3, Name: "Appliances"},
		{CategoryId: 3, Name: "Audio equipment"},
		{CategoryId: 3, Name: "Books"},
		{CategoryId: 3, Name: "Clothes"},
		{CategoryId: 3, Name: "Computers"},
		{CategoryId: 3, Name: "Furniture"},
		{CategoryId: 3, Name: "Gym equipment"},
		{CategoryId: 3, Name: "Sports equipment"},
	}
	db.Create(&subcategories)
}

func getAllPosts(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var posts []Post
	err := db.Find(&posts).Error
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(response).Encode(posts)
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
	}
}

func createPost(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Createpost")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var post Post
	err := json.NewDecoder(request.Body).Decode(&post)
	if err != nil {
		sendErr(writer, http.StatusBadRequest, err.Error())
		return
	}
	err = db.Save(&post).Error
	if err != nil {
		sendErr(writer, http.StatusInternalServerError, err.Error())
	} else {
		writer.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(writer).Encode("Post created")
		if err != nil {
			return
		}
	}
}

func getAllCategories(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var categories []Category
	err := db.Find(&categories).Error
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(response).Encode(categories)
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
	}
}

func getSubcategories(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var subcategories []Subcategory

	categoryId := mux.Vars(request)["id"]
	fmt.Println("CategoryId:", categoryId)
	err := db.Find(&subcategories, categoryId).Error
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(response).Encode(subcategories)
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
	}
}

func getPostById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var post Post
	param := mux.Vars(request)
	err := db.First(&post, param["id"]).Error
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(response).Encode(post)
	if err != nil {
		sendErr(response, http.StatusInternalServerError, err.Error())
	}
}

func updatePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var post Post
	param := mux.Vars(request)
	db.First(&post, param["id"])
	err := json.NewDecoder(request.Body).Decode(&post)
	if err != nil {
		return
	}
	db.Save(&post)
	err = json.NewEncoder(writer).Encode(post)
	if err != nil {
		return
	}
}

func deletePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var post Post
	param := mux.Vars(request)
	db.Delete(&post, param["id"])
	err := json.NewEncoder(writer).Encode("Record deleted successfully")
	if err != nil {
		return
	}
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

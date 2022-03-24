package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"tradelist/pkg/api"
	"tradelist/pkg/apihelpers"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims        //go get github.com/golang-jwt/jwt

}

func (server *Server) Signup(writer http.ResponseWriter, request *http.Request) {

	var data map[string]string // string as a key and string as a value
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	contact := api.Contact{
		FirstName: data["first name"],
		LastName:  data["last name"],
		Email:     data["email"],
		Password:  password,
	}
	user := api.User{
		Contact: contact,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(user)

}

func (server *Server) Login(writer http.ResponseWriter, request *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func (server *Server) Home(writer http.ResponseWriter, request *http.Request) {
	c, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	writer.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

}

func (server *Server) Refresh(writer http.ResponseWriter, request *http.Request) {
	c, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//commenting out for testing purposes
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	writer.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `refresh_token` cookie
	http.SetCookie(writer, &http.Cookie{
		Name:    "refresh_token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func (server *Server) Logout(writer http.ResponseWriter, request *http.Request) {
	cookie := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(writer, &cookie)

	writer.Write([]byte("Old cookie deleted. Logged out!\n"))

}

func (server *Server) CreatePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var post api.Post
	err := json.NewDecoder(request.Body).Decode(&post)

	if err != nil {
		sendErr(writer, http.StatusBadRequest, err.Error())
		return
	}
	response := server.PostService.Create(post)
	apihelpers.Respond(writer, response)
}

func (server *Server) GetAllPosts(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	response := server.PostService.GetAllPosts()
	apihelpers.Respond(writer, response)
}

func (server *Server) GetAllCategories(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	response := server.PostService.GetAllCategories()
	apihelpers.Respond(writer, response)
}

func (server *Server) GetSubcategories(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	categoryId := mux.Vars(request)["id"]
	fmt.Println("CategoryId:", categoryId)
	response := server.PostService.GetSubcategories(categoryId)
	apihelpers.Respond(writer, response)
}

func (server *Server) GetPostById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	postId := mux.Vars(request)["id"]
	fmt.Println("PostId:", postId)
	response := server.PostService.GetPostById(postId)
	apihelpers.Respond(writer, response)
}

func (server *Server) UpdatePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	postId := mux.Vars(request)["id"]
	fmt.Println("PostId:", postId)
	var post api.Post
	err := json.NewDecoder(request.Body).Decode(&post)

	if err != nil {
		sendErr(writer, http.StatusBadRequest, err.Error())
		return
	}
	response := server.PostService.UpdatePost(post, postId)
	apihelpers.Respond(writer, response)
}

func (server *Server) DeletePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	postId := mux.Vars(request)["id"]
	fmt.Println("PostId:", postId)
	response := server.PostService.DeletePost(postId)
	apihelpers.Respond(writer, response)
}

func (server *Server) CreateJobPost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var jobPost api.JobPost

	body, _ := ioutil.ReadAll(request.Body)
	err := json.Unmarshal(body, &jobPost)

	fmt.Println(jobPost)

	if err != nil {
		sendErr(writer, http.StatusBadRequest, "Malformed Post data")
		return
	}
	response := server.jobService.CreateJobPost(jobPost)
	apihelpers.Respond(writer, response)
}

func (server *Server) GetPostByCategoryId(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	categoryId := mux.Vars(request)["id"]
	fmt.Println("CategoryId:", categoryId)
	response := server.jobService.GetPostByCategoryId(categoryId)
	apihelpers.Respond(writer, response)
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

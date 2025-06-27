package main 

import (
	"fmt"
	"net/http"
    "encoding/json"
	"github.com/go-playground/validator/v10"
)

type websiteUser struct {
    Username  string `json:"username" validate:"required,min=5,max=20"`
    Email     string `json:"email" validate:"required,email"`
    Age       int    `json:"age" validate:"gte=18,lte=120"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user websiteUser

    
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    validate := validator.New() //validator instance 

    
    err = validate.Struct(user)
    if err != nil {
        errors := err.(validator.ValidationErrors)
        http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "User created successfully!")
}

func main_validator() {
    http.HandleFunc("/create-user", createUserHandler)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
//to create a test request for this api open postman and go to headers select content-type json body -> raw choose json and take this example:
// {
//   "username": "myuser",
//   "email": "me@example.com",
//   "age": 25
// }
//address : localhost:8080/create-user
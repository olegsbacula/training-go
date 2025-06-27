package main 
//validator + fasthttp + CORS + zap logger
import
(
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"os"
	"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
)
var (
	logger   *zap.Logger
	validate *validator.Validate
	users [] UserResponse
)
type PingResponse struct {
	Status string `json:"status"`
}

type UserResponse struct {
	Username string `json:"uname" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func handler(ctx *fasthttp.RequestCtx){
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	path := string(ctx.Path())

	switch path {
	case "/":
		ctx.SetContentType("text/plain")
		ctx.SetBodyString("Hello World")
		ctx.SetStatusCode(fasthttp.StatusOK)
		logger.Info("Hello world has been sent")

	case "/ping":
		resp := PingResponse{Status: "ok"}
		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			logger.Info("Error while pinging")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetBody(jsonBytes)
		ctx.SetStatusCode(fasthttp.StatusOK)
		logger.Info("Status has been sent")

	case "/user":
		resp := UserResponse{Username:"checking", Email:"abc@abc.de"}
		jsonuser,err := json.Marshal(resp)
		if err != nil {
			ctx.Error("Internal Server Error",fasthttp.StatusInternalServerError)
			logger.Info("Error while getting a user")
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetBody(jsonuser)
		ctx.SetStatusCode(fasthttp.StatusOK)
		logger.Info("user info has been sent")

	case "/create":
		if !ctx.IsPost() {
        	ctx.Error("Method Not Allowed", fasthttp.StatusMethodNotAllowed)
			logger.Info("Method that has been used is not correct")
        	return
    	}
		var user UserResponse	
    	err := json.Unmarshal(ctx.PostBody(), &user)
    	if err != nil {
        	ctx.Error("Bad Request", fasthttp.StatusBadRequest)
        	return
    	}
		err = validate.Struct(user)
        if err != nil {
            ctx.SetStatusCode(fasthttp.StatusBadRequest)
            ctx.SetContentType("application/json")
            validationErrors := err.(validator.ValidationErrors)
            errorMsg, _ := json.Marshal(map[string]interface{}{
                "error": validationErrors.Error(),
            })
			logger.Info("Validation error happened")
            ctx.SetBody(errorMsg)
            return
        }
		users = append(users, user)
		response := map[string]string{
        	"message": "User created",
        	"username": user.Username,
        	"email": user.Email,
    	}
		jsonResp, err := json.Marshal(response)
    	if err != nil {
        	ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			logger.Info("Error while sending the response back")
        	return
    	}
    	ctx.SetContentType("application/json")
    	ctx.SetBody(jsonResp)
    	ctx.SetStatusCode(fasthttp.StatusOK)

	case "/users":
		users = append(users, UserResponse{
			Username: "olegs",
			Email:    "email@gmail.com",
		})
		jsonsusers, err := json.Marshal(users)
		if err != nil {
			ctx.Error("Didn't manage to marshal a json file", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetBody(jsonsusers)
		ctx.SetStatusCode(fasthttp.StatusOK)
		logger.Info("Status has been sent")

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

func main(){
// to run use : go run . and change env What
	_ = godotenv.Load()
	if os.Getenv("WhatToRun") == "zap"{
		main_zap()
	}
	if os.Getenv("WhatToRun") == "valid_custom"{
		main_valid_custom()
	}
	if os.Getenv("WhatToRun") == "validator"{
		main_validator()
	}
	if os.Getenv("WhatToRun") == "fasthttp" {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		validate = validator.New()
		if err := fasthttp.ListenAndServe(":8080", handler); err != nil {
			panic(err)
		}
	}
	if os.Getenv("WhatToRun") == "socket" {
		main_socket()
	}
}

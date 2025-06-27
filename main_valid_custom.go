package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)
type User struct {
	Name string `validate:"min=2,max=32"`
	Email string `validate:"required,email"`

}

func validate_values(val interface{}) error { 
	v:=reflect.ValueOf(val)
	fmt.Println(v)
	for i:= 0; i<v.NumField();i++{
		field:=v.Field(i)
		fmt.Println(field)
		tag:=v.Type().Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}
		rules:=strings.Split(tag,",")//reuired,email
		for _,rule := range rules{
			fieldName:=v.Type().Field(i).Name
			switch {
				//"min=2" 
			case strings.HasPrefix(rule,"min="):
				min,_ := strconv.Atoi(strings.TrimPrefix(rule,"min="))
				if len(field.String())< min {
					err := fmt.Errorf("Field %s is too short to be more than a minimum %d",fieldName,min)
					zap.L().Error("Validation error", zap.Error(err))
					return err 
				}
			case strings.HasPrefix(rule,"max="):
				max,_ := strconv.Atoi(strings.TrimPrefix(rule,"max="))
				if len(field.String()) > max {
					err := fmt.Errorf("Field %s is too long to be more than a maximum %d",fieldName,max)
					zap.L().Error("Validation error", zap.Error(err))
					return err 	
			}
		case rule == "required":
			if field.String() == ""{
					err := fmt.Errorf("Field %s is required to be more than empty",fieldName)
					zap.L().Error("Validation error", zap.Error(err))
					return err 
			}
		case rule == "email":
			emailRegex:=regexp.MustCompile((`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`))
			if !emailRegex.MatchString(field.String()){
					err := fmt.Errorf("Field %s must be an email",fieldName)
					zap.L().Error("Validation error", zap.Error(err))
					return err 
			}
		}
	}
	}
	return nil

}

func main_valid_custom() {
	user := User{
		Name:"Oleg",
		Email:"abc@abc.lv",
		}
		fmt.Println(validate_values(user))
	invalidUser := User{
		Name:"o",
		Email:"abcabc.lv",
		}
		fmt.Println(validate_values(invalidUser))
	// t := reflect.TypeOf(user)
	// fmt.Println("Name:", t.Name())
	// fmt.Println("Kind:", t.Kind())
	// for i := 0; i<t.NumField();i++{
	// 	field := t.Field(i)
    // 	tag := field.Tag.Get("validate")
    // 	fmt.Printf("%d. %v (%v), tag: `%v`\n", i+1, field.Name, field.Type.Name(), tag)
	// }


}

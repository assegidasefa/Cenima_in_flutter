package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil{
		_, _ = w.Write([]byte(err.Error()))
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error){
	w.Header().Set("Content-Type", "application/json")
	if err != nil{
		log.Println(err)
		JSON(w, statusCode, struct {
			Error string `json:"message"`
		}{
			Error: err.Error(),
		})
		return
	}
}

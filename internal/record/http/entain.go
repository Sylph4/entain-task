package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sylph4/entain-task/internal/record/model"
	"github.com/sylph4/entain-task/internal/record/service"
)

type ProcessRecordHandler struct {
	processRecordService service.IProcessRecordService
	userService          *service.UserService
	validate             *validator.Validate
}

func NewProcessRecordHandler(processRecordService service.IProcessRecordService, userService *service.UserService, validate *validator.Validate) *ProcessRecordHandler {
	return &ProcessRecordHandler{
		processRecordService: processRecordService,
		userService:          userService,
		validate:             validate,
	}
}

func (h *ProcessRecordHandler) ProcessRecord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	processRecordRequest := &model.ProcessRecordRequest{}
	err := decoder.Decode(&processRecordRequest)
	if err != nil {
		log.Println("DECODE FAIL", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.validate.Struct(processRecordRequest)
	if err != nil {
		log.Println("Request validation error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.processRecordService.ProcessRecord(processRecordRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.StatusText(200)
}

func (h *ProcessRecordHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	//nolint
	w.Write(response)
}

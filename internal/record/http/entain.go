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
	validate             *validator.Validate
}

func NewProcessRecordHandler(
	processRecordService service.IProcessRecordService,
	validate *validator.Validate) *ProcessRecordHandler {
	return &ProcessRecordHandler{
		processRecordService: processRecordService,
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

package handlers

import (
	"GoPass/internal/server/records"
	mw "GoPass/internal/server/transport/middleware"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RecordHandler struct {
	OrderUseCase records.UseCase
}

func NewRecordHandler(recordUseCase records.UseCase) *RecordHandler {
	return &RecordHandler{
		OrderUseCase: recordUseCase,
	}
}

func (rec *RecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var newRecord records.Record
	if err := json.Unmarshal(body, &newRecord); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}
	newRecord.UserID = currentUserID
	record, err := rec.OrderUseCase.Create(r.Context(), &newRecord)

	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "error creating record", http.StatusInternalServerError)
		return
	}

	// Отправка ответа с созданной записью
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)

}

func (rec *RecordHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ur records.Record
	if err := json.Unmarshal(body, &ur); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	_, err = rec.OrderUseCase.Update(r.Context(), ur.ID, currentUserID, ur.Name, ur.Site, ur.Login, ur.Password, ur.Info)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "error deleting record", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rec *RecordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	recordID := chi.URLParam(r, "id")
	fmt.Println(recordID)
	if recordID == "" {
		http.Error(w, "Record ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(recordID)
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	err = rec.OrderUseCase.Delete(r.Context(), id, currentUserID)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "error deleting record", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rec *RecordHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	recs, err := rec.OrderUseCase.List(r.Context(), currentUserID)
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(recs); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (rec *RecordHandler) GetById(w http.ResponseWriter, r *http.Request) {

}

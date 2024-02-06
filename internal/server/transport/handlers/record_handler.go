package handlers

import (
	"GoPass/internal/server/records"
	mw "GoPass/internal/server/transport/middleware"
	"encoding/json"
	"net/http"
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

}

func (rec *RecordHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (rec *RecordHandler) Delete(w http.ResponseWriter, r *http.Request) {

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

//Create(ctx context.Context, record *Record) (*Record, error)
//Update(ctx context.Context, id int, name, site, login, password, info string) (*Record, error)
//Delete(ctx context.Context, id int) error
//List(ctx context.Context, userID int) ([]*Record, error)
//GetById(ctx context.Context, id int) (*Record, error)

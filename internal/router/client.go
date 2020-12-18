package router

import (
	"encoding/json"
	"net/http"

	"github.com/1gkx/salary/internal/store"
)

// SetClientHendler ...
func SetClientHendler() {
	r.HandleFunc("/client", set).Methods("POST")
	r.Handle("/approve", authRequireHandlerWrap(approve)).Methods("POST")
}

// set ...
func set(w http.ResponseWriter, r *http.Request) {

	client := new(store.Client)

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := store.SetClient(client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
	return
}

// approve ...
func approve(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	client := new(store.Client)

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":   501,
			"status": err,
		})
		return
	}

	// Approved with IBSO
	// res, err := utils.Post(
	// 	client,
	// 	conf.Cfg.Methods["APPROVE"],
	// )
	// if err != nil {
	// 	w.WriteHeader(501)
	// 	json.NewEncoder(w).Encode(err.Error())
	// 	return
	// }

	// Если утверждение ок, то удалить из базы!

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   200,
		"status": "success",
	})
	return
}

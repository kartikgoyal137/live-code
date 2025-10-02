package main 
import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status" : "ok"}

	json.NewEncoder(w).Encode(response);
}


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")

	port := ":8080"
	log.Println("Server starting on PORT : ", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("ListenAndServe : ", err)
	}

}
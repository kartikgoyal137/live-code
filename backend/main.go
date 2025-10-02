package main 
import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r , nil)
	if err!=nil {
		log.Println("upgrade error : ", err)
	}
	defer ws.Close()

	log.Println("Client successfully connected...")

	for {
		messageType, p, err := ws.ReadMessage()
		if err!=nil {
			log.Println("read error : ", err)
			break
		}

		log.Printf("Recieved message : %s of type %v", p, messageType)

		if err:=ws.WriteMessage(messageType, p); err!=nil {
			log.Println("write error : ", err)
			break
		}
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status" : "ok"}

	json.NewEncoder(w).Encode(response);
}


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	r.HandleFunc("/ws", handleSocket)

	port := ":8080"
	log.Println("Server starting on PORT : ", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("ListenAndServe : ", err)
	}

}
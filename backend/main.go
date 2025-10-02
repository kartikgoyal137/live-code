package main 
import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"live-code/backend/ws"
	"live-code/backend/docker"
)

func handleSocket(hub *ws.Hub, manager *docker.Manager ,w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r , nil)
	if err!=nil {
		log.Println("upgrade error : ", err)
		return
	}

	log.Println("Client successfully connected...")

	client := &ws.Client{Hub: hub, Manager: manager, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status" : "ok"}

	json.NewEncoder(w).Encode(response);
}


func main() {
	r := mux.NewRouter()
	dockerManager, err := docker.NewManager()
	if err!=nil {
		log.Println("Failed to connect to Docker Daemon : ", err)
		return
	}

	hub := ws.NewHub()
	go hub.Run()

	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleSocket(hub , dockerManager, w , r)
	})

	port := ":8080"
	log.Println("Server starting on PORT : ", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("ListenAndServe : ", err)
	}

}
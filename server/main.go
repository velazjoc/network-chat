package main
import "encoding/json"
import "log"
import "math/rand"
import "net/http"
import "sync"
import "time"
import "github.com/gorilla/websocket"

type Hub struct {
	mu    sync.Mutex
	rooms map[string]map[*websocket.Conn]string 
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*websocket.Conn]string)}
}

func (h *Hub) join(roomId string, conn *websocket.Conn, name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomId]; !ok {
		h.rooms[roomId] = make(map[*websocket.Conn]string)
	}
	h.rooms[roomId][conn] = name
}

func (h *Hub) leave(roomId string, conn *websocket.Conn) (name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomId]; !ok {
		return ""
	}
	name = h.rooms[roomId][conn]
	delete(h.rooms[roomId], conn)
	if len(h.rooms[roomId]) == 0 {
		delete(h.rooms, roomId)
	}
	return name
}

func (h *Hub) broadcast(roomId string, msg any) {
	h.mu.Lock()
	peers := h.rooms[roomId]
	h.mu.Unlock()

	if peers == nil {
		return
	}

	data, _ := json.Marshal(msg)
	for c := range peers {
		_ = c.WriteMessage(websocket.TextMessage, data)
	}
}

func (h *Hub) presence(roomId string) []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	peers := h.rooms[roomId]
	out := make([]string, 0, len(peers))
	for _, name := range peers {
		out = append(out, name)
	}
	return out
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev-only
}

func randRoom() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	hub := NewHub()

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		roomId := randRoom()
		_ = json.NewEncoder(w).Encode(map[string]string{"roomId": roomId})
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("room")
		name := r.URL.Query().Get("name")
		if roomId == "" || name == "" {
			http.Error(w, "room and name required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		hub.join(roomId, conn, name)
		hub.broadcast(roomId, map[string]any{"type": "system", "message": name + " joined"})
		hub.broadcast(roomId, map[string]any{"type": "presence", "users": hub.presence(roomId)})

		defer func() {
			leftName := hub.leave(roomId, conn)
			_ = conn.Close()
			if leftName != "" {
				hub.broadcast(roomId, map[string]any{"type": "system", "message": leftName + " left"})
				hub.broadcast(roomId, map[string]any{"type": "presence", "users": hub.presence(roomId)})
			}
		}()

		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var incoming map[string]any
			if json.Unmarshal(raw, &incoming) != nil {
				continue
			}
			if incoming["type"] == "chat" {
				text, _ := incoming["text"].(string)
				if text == "" {
					continue
				}
				hub.broadcast(roomId, map[string]any{
					"type": "chat",
					"name": name,
					"text": text,
					"ts":   time.Now().UnixMilli(),
				})
			}
		}
	})

	log.Println("Go server: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
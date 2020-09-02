package routes

import (
	"log"
	"net/http"
	"strconv"
	"time"

	okta "../controllers/okta"
	auth "../utils/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type msg struct {
	Num int
}

// wsEndpoint used for websocket connection test
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//go startPolling(ws, err)

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	//reader(ws)
	go echo(ws)

	//select {}
}

// startPolling used to get data
func startPolling(ws *websocket.Conn, err error) {
	counter := 1
	for _ = range time.Tick(2 * time.Second) {
		err = ws.WriteMessage(1, []byte(strconv.Itoa(counter)))
		if err != nil {
			log.Println(err)
		}
		log.Println(counter)
		// doSomething("awesome")
		counter++
	}
}

// echo used for ws testing
func echo(conn *websocket.Conn) {
	counter := 1
	for _ = range time.Tick(2 * time.Second) {
		err := conn.WriteMessage(1, []byte(strconv.Itoa(counter)))
		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
		log.Println(counter)
		// doSomething("awesome")
		counter++
	}
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2024,
	WriteBufferSize: 2024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

// Handlers will route to necessary controllers
func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	r.HandleFunc("/okta/introspect", okta.CheckToken).Methods("POST")
	r.HandleFunc("/ws", wsEndpoint)

	// Auth Routes that REQUIRE JWT Verification
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/okta/introspect", okta.CheckToken).Methods("POST")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, x-access-token, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, X-Auth-Token")
		next.ServeHTTP(w, r)
	})
}

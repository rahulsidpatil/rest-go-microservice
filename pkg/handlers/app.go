package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/rahulsidpatil/rest-go-microservice/pkg/dal"
	"github.com/rahulsidpatil/rest-go-microservice/pkg/util"
	_ "github.com/rahulsidpatil/rest-go-microservice/pkg/util"

	"github.com/gorilla/mux"
	"github.com/rahulsidpatil/rest-go-microservice/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

var svcPort string

type App struct {
	router *mux.Router
	db     dal.Interface
}

func (a *App) Initialize() {
	svcVersion := os.Getenv("SVC_VERSION")
	if svcVersion == "" {
		svcVersion = "/v1"
	}
	svcPort = os.Getenv("SVC_PORT")
	if svcPort == "" {
		svcPort = "8080"
	}
	svcPathPrefix := os.Getenv("SVC_PATH_PREFIX")
	if svcPathPrefix == "" {
		svcPathPrefix = "messages"
	}

	svcPathPrefix = svcVersion + "/" + svcPathPrefix
	swaggerAddr := "http://localhost:" + svcPort + "/swagger/doc.json"

	a.router = mux.NewRouter()
	a.db = dal.GetMySQLDriver()
	a.setSwaggerInfo(swaggerAddr, svcPort, svcVersion)
	a.initializeRoutes(swaggerAddr, svcPathPrefix)
}

func (a *App) Run() {
	svcAddr := ":" + svcPort
	log.Printf("starting message server at:%s", svcAddr)
	log.Println(http.ListenAndServe(svcAddr, a.router))
}

func (a *App) initializeRoutes(swaggerAddr, svcPathPrefix string) {
	a.router.HandleFunc(os.Getenv("SVC_VERSION")+"/hello", WithStats(a.hello)).Methods("GET")
	a.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	a.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerAddr), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	a.router.HandleFunc(svcPathPrefix, WithStats(a.getAll)).Methods("GET")
	a.router.HandleFunc(svcPathPrefix, WithStats(a.addMessage)).Methods("POST")
	a.router.HandleFunc(svcPathPrefix+"/{id:[0-9]+}", WithStats(a.getMessage)).Methods("GET")
	a.router.HandleFunc(svcPathPrefix+"/{id:[0-9]+}", WithStats(a.updateMessage)).Methods("PUT")
	a.router.HandleFunc(svcPathPrefix+"/{id:[0-9]+}", WithStats(a.deleteMessage)).Methods("DELETE")
	a.router.HandleFunc(svcPathPrefix+"/palindromeChk/{id:[0-9]+}", WithStats(a.palindromeChk)).Methods("GET")
}

func (a *App) setSwaggerInfo(swaggerAddr, port, version string) {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample rest-go-microservice server."
	docs.SwaggerInfo.Version = "1.0"
	//TODO: remove hard-coding of Host address
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = version
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

/*
WithStats ... wraps handlers with stats reporting.
It tracks metrics such as the request recieve time and latency
*/
func WithStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		tags := util.GetStatsTags(r)
		util.RequestFrom(tags, start)

		h(w, r)

		duration := time.Since(start)
		util.RecordLatency(tags, duration)
	}
}

// @Summary Say hello to user
// @Description Say hello to user
// @Success 200 {object} string
// @Failure 404
// @Router /hello [get]
func (a *App) hello(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Hello rest-go-microservice..!!!")
}

// @Summary Get all messages
// @Description Get all messages
// @Success 200 {array} dal.Message
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages [get]
func (a *App) getAll(w http.ResponseWriter, r *http.Request) {
	messages, err := a.db.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, messages)
}

// @Summary Add new messages
// @Description Add new messages
// @Accept  json
// @Produce  json
// @Param message body dal.Message true "Add message"
// @Success 200 {object} dal.Message
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages [post]
func (a *App) addMessage(w http.ResponseWriter, r *http.Request) {
	var msg dal.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := a.db.AddMessage(&msg); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, msg)
}

// @Summary Fetch message by ID
// @Description Fetch message by ID
// @Accept  json
// @Produce  json
// @Param id path int true "Message ID"
// @Success 200 {object} dal.Message
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages/{id} [get]
func (a *App) getMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message id")
		return
	}

	msg := dal.Message{ID: id}
	if err := a.db.GetMessage(&msg); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Message not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, msg)
}

// @Summary Check if the message specified by ID is a palindrome or not
// @Description Check if the message specified by ID is a palindrome or not
// @Accept  json
// @Produce  json
// @Param id path int true "Message ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages/palindromeChk/{id} [get]
func (a *App) palindromeChk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message id")
		return
	}

	msg := dal.Message{ID: id}
	if err := a.db.GetMessage(&msg); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Message not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	palindromeChk := util.Palindrome(msg.Message)
	response := map[string]interface{}{
		"Message":    msg,
		"Palindrome": palindromeChk,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// @Summary Update message by ID
// @Description Update message by ID
// @Accept  json
// @Produce  json
// @Param id path int true "Message ID"
// @Param message body dal.Message true "Update message"
// @Success 200 {object} dal.Message
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages/{id} [put]
func (a *App) updateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var msg dal.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	msg.ID = id

	if err := a.db.UpdateMessage(&msg); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, msg)
}

// @Summary Delete message by ID
// @Description Delete message by ID
// @Accept  json
// @Produce  json
// @Param id path int true "Message ID"
// @Success 200 {object} dal.Message
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Router /messages/{id} [delete]
func (a *App) deleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	msg := dal.Message{ID: id}
	if err := a.db.DeleteMessage(&msg); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (server *Server) dashboardEndpoint(w http.ResponseWriter, r *http.Request) {
	var filter, option interface{}
	filter = bson.D{}

	opts := options.Find()
	opts.SetSort(bson.D{{}})

	context := r.Context()
	cursor, err := query(server.client, context, "Slackbot", "Standups", filter, option)
	if err != nil {
		log.Println(err)
	}

	var results []Standup
	if err := cursor.All(context, &results); err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}

func (server *Server) ResponsesEndpoint(w http.ResponseWriter, r *http.Request) {
	//read json payload from request
	var response Response
	response.ReadToResponses(w, r)
	//store into a struct
	//store in the response collection in slackbot database
	context := r.Context()
	insertOneResult, err := UpsertOne(server.client, context, "Slackbot",
		"Responses", response, response.Sender)
	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusServiceUnavailable)
	}
	fmt.Println(insertOneResult)
}

func (server *Server) configureStandupEndpoint(w http.ResponseWriter, r *http.Request) {
	var s Standup
	s.ReadToStandup(w, r)
	context := r.Context()
	Result, err := UpsertOne(server.client, context, "Slackbot",
		"Standups", s, s.Name)
	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusServiceUnavailable)
	}
	fmt.Println(Result)
}

func (server *Server) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	var s Standup
	s.ReadToStandup(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
func main() {

	token := "xoxb-3472770742805-3478067021124-wDsAFtX4Sok6z369uMFYzm2a"
	channelID := "C03E2133QR2"
	appToken := "xapp-1-A03DT21MT7G-3506826440112-7d05e9c046e62ce14e937feb2266d86e27764008f9fc91c5e0195a9582d759cb"
	fmt.Println("Server is UP!!!!")
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		log.Println(err)
	}

	MessageBot(token, channelID, "hiii", "scrum", "Are you fine")
	Getmessage(token, appToken)
	Server := NewServer(client, cancel)
	defer close(Server.client, ctx, Server.cancel)
	newRouter := mux.NewRouter().StrictSlash(true)
	newRouter.HandleFunc("/dashboard", Server.dashboardEndpoint).Methods("GET")
	newRouter.HandleFunc("/dashboard/{name}", Server.configureStandupEndpoint).Methods("GET")
	newRouter.HandleFunc("/dashboard/{name}/configure", Server.configureStandupEndpoint).Methods("POST")
	log.Println(http.ListenAndServe(":8080", newRouter))
}

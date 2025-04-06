package main

import (
	"blogService/data"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JSONPayload struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (app *Config) AddBlog(w http.ResponseWriter, r *http.Request) {

	var requestPayload JSONPayload
	_ = app.readJson(w, r, &requestPayload)

	event := data.BlogData{
		Name:        requestPayload.Name,
		Author:      requestPayload.Author,
		Description: requestPayload.Description,
	}

	documentId, err := app.Models.BlogData.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	id := documentId.InsertedID.(primitive.ObjectID).Hex()
	err = app.logRequest("Blog Service", fmt.Sprintf("Blog created with ID: %s", id))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Blog created successfully",
		Data: data.BlogData{
			ID:          id,
			Name:        event.Name,
			Author:      event.Author,
			Description: event.Description,
		},
	}

	app.writeJson(w, http.StatusAccepted, resp)
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) GetAllBlog(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.Models.BlogData.All()

	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error: false,
		Data:  blogs,
	}
	app.writeJson(w, http.StatusAccepted, resp)

}

type DeleteBlogRequest struct {
	ID string `json:"id"`
}

func (app *Config) DeleteBlog(w http.ResponseWriter, r *http.Request) {

	var requestPayload DeleteBlogRequest
	_ = app.readJson(w, r, &requestPayload)

	event := data.BlogData{
		ID: requestPayload.ID,
	}

	err := app.Models.BlogData.Delete(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.logRequest("Blog Service", fmt.Sprintf("Blog deleted with ID: %s", event.ID))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "Blog deleted successfully",
	}
	app.writeJson(w, http.StatusAccepted, resp)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service:8082/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

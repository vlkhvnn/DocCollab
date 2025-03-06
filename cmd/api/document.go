package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CreateDocumentPayload struct {
	Content string `json:"content"`
}

func (app *application) createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.unauthorizedErrorResponse(w, r, errors.New("missing authorization header"))
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)
	_, err := app.authenticator.ValidateToken(tokenStr)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, errors.New("invalid token"))
		return
	}

	// Decode payload (only content is provided by the client)
	var payload CreateDocumentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newDocID := uuid.New().String()

	ctx := r.Context()
	if err := app.store.Document.CreateDocument(ctx, newDocID, payload.Content); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Retrieve the document (to get timestamps, etc.)
	doc, err := app.store.Document.GetDocumentByDocID(ctx, newDocID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Return the created document as JSON.
	app.jsonResponse(w, http.StatusCreated, doc)
}

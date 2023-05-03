package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cf *apiConfig) handlePostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string
		Email    string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlePostUser: " + err.Error())
		return
	}

	newUser, err := cf.db.CreateUser(params.Password, params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create User")
		log.Println("handlePostUser: " + err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{newUser.Id, newUser.Email})
}

func (cf *apiConfig) handlePutUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string
		Email    string
	}

	signedToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		log.Println("handlePutUser: " + err.Error())
		return
	}

	userId, err := cf.validateJWT(signedToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "")
		log.Println("handlePutUser: " + err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	e := decoder.Decode(&params)

	if e != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlePutUser: " + err.Error())
		return
	}

	updatedUser, err := cf.db.UpdateUser(userId, params.Password, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user")
	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}

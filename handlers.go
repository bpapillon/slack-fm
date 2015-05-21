package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello Umbel.fm!")
}

func RecommendationListHandler(w http.ResponseWriter, r *http.Request) {
    response, err := GetRecommendationMessages()
    if err != nil {
        fmt.Fprintln(w, err)
        return
    }
    responseBody, err := json.Marshal(response)
    if err != nil {
        fmt.Fprintln(w, err)        
    } else {
        fmt.Fprintln(w, string(responseBody))
    }
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    response, err := GetUser(mux.Vars(r)["userId"])
    if err != nil {
        fmt.Fprintln(w, err)
        return
    }
    responseBody, err := json.Marshal(response)
    if err != nil {
        fmt.Fprintln(w, err)
    } else {
        fmt.Fprintln(w, string(responseBody))
    }
}

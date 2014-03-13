package main

import (
	"encoding/json"
	"github.com/300brand/subscription/config"
	"net/http"
)

func init() {
	http.HandleFunc("/config/save", saveConfig)
	http.HandleFunc("/config/view", viewConfig)
}

func saveConfig(w http.ResponseWriter, r *http.Request) {
	if err := config.Save(*ConfigFile); err != nil {
		http.Error(w, "Could not save config file: "+err.Error(), http.StatusInternalServerError)
	}
}

func viewConfig(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(config.Config); err != nil {
		http.Error(w, "Error displaying config: "+err.Error(), http.StatusInternalServerError)
	}
}

package handlers

import (
	"log"
	"net/http"
)

// Home
// Renders home.jet page
func Home(w http.ResponseWriter, r *http.Request) {
  err := RenderPage(w, "home.jet", nil)
  if err != nil {
    log.Println("Unable to render home.jet page", err)
  }
}

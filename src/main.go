package main

import (
	"net/http"
	"sec_2/controllers"
)

func main() {
	controllers.RegisterControllers()
	http.ListenAndServe(":8080", nil)
}

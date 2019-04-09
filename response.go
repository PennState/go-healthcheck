package healthcheck

import "net/http"

func GetHealthHandler(checkers ...Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: iterate over checks and determine overall status
		//      return handler function
	}
}

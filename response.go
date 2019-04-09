package health

import "net/http"

func GetHealthHandler(checks ...Check) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: iterate over checks and determine overall status
		//      return handler function
	}
}

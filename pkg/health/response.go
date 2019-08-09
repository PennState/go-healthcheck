package healthcheck

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetHealthHandler(checkers ...Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checks Checks
		status := Pass
		for _, checker := range checkers {
			c, s := checker.Check()
			checks.AddChecks(c...)
			status = status.Max(s)
		}
		log.WithField("Status", status).Debug("Checks: ", checks)
	}
}

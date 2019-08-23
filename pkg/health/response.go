package health

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Data []byte
	Code int
}

// TODO - Remove and replace with check functions
type Checker interface {
	Check() ([]ComponentDetail, Status)
}

func GetHealthHandler(checkers ...Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checks := Checks{}
		status := Pass
		for _, checker := range checkers {
			c, s := checker.Check()
			checks.Add(c[0].Key, c...)
			status = status.Max(s)
		}

		resp, err := json.Marshal(checks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errMsg := err.Error()
			_, writeError := w.Write([]byte(errMsg))
			if writeError != nil {
				log.WithError(writeError).WithContext(r.Context()).WithField("Status", status).Errorf("Unable to write healthcheck error %v", err)
			} else {
				log.WithError(err).WithContext(r.Context()).WithField("Status", status).Error("Unable to marshal checks")
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.StatusCode())
		_, err = w.Write(resp)
		if err != nil {
			log.WithError(err).WithContext(r.Context()).WithField("Status", status).Error("Unable to write healthcheck response")
		}
		log.WithContext(r.Context()).WithField("Status", status).Debug("Checks: ", checks)
	}
}

package request

import (
	"encoding/json"
	"net/http"
)

// BindJSON binds the request into an interface
func BindJSON(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

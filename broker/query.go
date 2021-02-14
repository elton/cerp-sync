package broker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/elton/cerp-sync/utils/logger"
	"github.com/elton/cerp-sync/utils/signatures"
	"github.com/joho/godotenv"
)

// An InvalidPtrError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidPtrError struct {
	Type reflect.Type
}

func (e *InvalidPtrError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func query(request map[string]interface{}, responseObject interface{}) error {
	rv := reflect.ValueOf(responseObject)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidPtrError{reflect.TypeOf(responseObject)}
	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	apiURL := os.Getenv("apiURL")

	reqJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	sign := signatures.Sign(string(reqJSON), os.Getenv("secret"))
	request["sign"] = sign

	if reqJSON, err = json.Marshal(request); err != nil {
		return err
	}

	logger.Info.Printf("Request JSON:%s \n", string(reqJSON))

	response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqJSON))

	if err != nil {
		return err
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	json.Unmarshal(responseData, &responseObject)

	return nil
}

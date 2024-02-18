package wasmhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type Request struct {
	Program string `json:"program"`
}
type SuccessResponse struct {
	WASM string `json:"wasm"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func init() {
	dir, err := WorkingDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", dir + "/keys/storage-admin.json")
	bucketWriter, err = NewBucketWriter()
	if err != nil {
		log.Fatal(err)
	}
	functions.HTTP("HandleWasm", ServeHTTP)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Closure for http response.
	handleResponse := func(statusCode int, payload any, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(payload)
	}

	if r.Method != http.MethodPost {
		handleResponse(http.StatusMethodNotAllowed, ErrorResponse{Error: "Method not allowed."}, nil)
		return
	}

	var req Request
	if b, err := io.ReadAll(r.Body); err != nil {
		handleResponse(http.StatusUnprocessableEntity, ErrorResponse{Error: "Unprocessable entity."}, err)
		return
	} else if err := json.Unmarshal(b, &req); err != nil {
		handleResponse(http.StatusBadRequest, ErrorResponse{Error: "Invalid request."}, err)
		return
	}

	p, err := NewProc(req.Program)
	if err != nil {
		handleResponse(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error."}, err)
		return
	}
	url, fault, err := p.DoProcess()
	if err != nil && fault == ServerErr {
		handleResponse(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error."}, err)
		return
	} else if err != nil {
		handleResponse(http.StatusBadRequest, ErrorResponse{Error: "Invalid request."}, err)
		return
	}
	// Ok.
	handleResponse(http.StatusCreated, SuccessResponse{WASM: url}, nil)
}

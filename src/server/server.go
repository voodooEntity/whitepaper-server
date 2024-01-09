package server

import (
	"encoding/json"
	"errors"
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/query"
	"github.com/voodooEntity/whitepaper-server/src/config"
	"github.com/voodooEntity/whitepaper-server/src/whitepaper"
	"io/ioutil"
	"net/http"
	"strings"
)

var ServeMux = http.NewServeMux()

func Start() {
	archivist.Info("> Booting HTTP API")

	ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The "/" matches anything not handled elsewhere. If it's not the root
		// then report not found.
		archivist.Debug(r.Method)
		archivist.Debug(r.URL.Path)
		archivist.DebugF("%+v", r.URL.Query())
		archivist.Debug(r.URL.Path)
		archivist.Debug(r.URL.Scheme)
		archivist.Debug(r.URL.RequestURI())
		respond("pong", 200, w)
	})

	// Route: /v1/ping
	//ServeMux.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
	//	respond("pong", 200, w)
	//})

	// Route: /v1/mapJson
	//ServeMux.HandleFunc("/api/Whitepaper/", func(w http.ResponseWriter, r *http.Request) {
	//	archivist.DebugF("Incoming request: %+v", r.Method)
	//	if "" != config.GetValue("CORS_ORIGIN") || "" != config.GetValue("CORS_HEADER") {
	//		if "OPTIONS" == r.Method {
	//			respond("", 200, w)
	//			return
	//		}
	//	}

	//	switch r.Method {
	//	case http.MethodPost:
	//		CreateOrUpdateWhitePaper(w, r)
	//	case http.MethodGet:
	//		GetWhitePaper(w, r)
	//	default:
	//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	}

	//})

	// building server listen string by
	// config values and print it - than listen
	connectString := buildListenConfigString()
	archivist.Info("> Server listening settings by config (" + connectString + ")")
	//http.ListenAndServe(connectString, ServeMux)
	http.ListenAndServeTLS(connectString, config.GetValue("SSL_CERT_FILE"), config.GetValue("SSL_KEY_FILE"), ServeMux)
}

func CreateOrUpdateWhitePaper(w http.ResponseWriter, r *http.Request) {
	// retrieve data from request
	body, err := getRequestBody(r)
	if nil != err {
		archivist.Error("Could not read http request body", err.Error())
		http.Error(w, "Malformed or no body. ", 422)
		return
	}

	// unpack the json
	var whitepaperData whitepaper.WhitePaper
	if err := json.Unmarshal(body, &whitepaperData); err != nil {
		archivist.Error("Invalid json query object", err.Error())
		http.Error(w, "Invalid json query object "+err.Error(), 422)
		return
	}

	archivist.DebugF("Unpacked payload: %+v", whitepaperData)

	// create whitepaper and respond
	whitepaperData.StoreOrUpdate()
	respond("", 200, w)
}

func GetWhitePaper(w http.ResponseWriter, r *http.Request) {
	instanceId, err := getIdFromUrl(w, r)
	if err != nil {
		respond(err.Error(), http.StatusBadRequest, w)
		return
	}

	// first we get the params
	requiredUrlParams := map[string]string{"hash": ""}
	urlParams, err := getRequiredUrlParams(requiredUrlParams, r)
	if nil != err {
		respond(err.Error(), 500, w)
		return
	}

	qry := query.New().Read("Whitepaper").Match("Value", "==", instanceId).Match("Properties.Hash", "==", urlParams["hash"])
	res := query.Execute(qry)
	if 1 == res.Amount {
		respond("", 200, w)
		return
	}
	qry = query.New().Read("Whitepaper").Match("Value", "==", instanceId)
	res = query.Execute(qry)
	archivist.DebugF("Loaded whiptepaper entity %+v", res)
	whitePaper := whitepaper.Load(instanceId)
	if nil == whitePaper {
		respond("unknown instance id given", 500, w)
		return
	}

	respondOk(whitePaper, w)
}

func getOptionalUrlParams(optionalUrlParams map[string]string, urlParams map[string]string, r *http.Request) map[string]string {
	tmpParams := r.URL.Query()
	for paramName := range optionalUrlParams {
		val, ok := tmpParams[paramName]
		if ok {
			urlParams[paramName] = val[0]
		}
	}
	return urlParams
}

func getIdFromUrl(w http.ResponseWriter, r *http.Request) (string, error) {
	// Extract the ID from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 0 {
		return "", errors.New("Url param missing")
	}

	instanceId := parts[len(parts)-1]

	return instanceId, nil
}

func getRequiredUrlParams(requiredUrlParams map[string]string, r *http.Request) (map[string]string, error) {
	urlParams := r.URL.Query()
	for paramName := range requiredUrlParams {
		val, ok := urlParams[paramName]
		if !ok {
			return nil, errors.New("Missing required url param")
		}
		requiredUrlParams[paramName] = val[0]
	}
	return requiredUrlParams, nil
}

func respond(message string, responseCode int, w http.ResponseWriter) {

	corsAllowHeaders := config.GetValue("CORS_HEADER")
	if "" != corsAllowHeaders {
		w.Header().Add("Access-Control-Allow-Headers", corsAllowHeaders)
	}
	corsAllowOrigin := config.GetValue("CORS_ORIGIN")
	if "" != corsAllowOrigin {
		w.Header().Add("Access-Control-Allow-Origin", corsAllowOrigin)
	}

	w.WriteHeader(responseCode)
	messageBytes := []byte(message)

	_, err := w.Write(messageBytes)
	if nil != err {
		archivist.Error("Could not write http response body ", err, message)
	}
}

func respondOk(data *whitepaper.WhitePaper, w http.ResponseWriter) {
	// than we gonne json encode it
	// build the json
	responseData, err := json.Marshal(&data)
	if nil != err {
		http.Error(w, "Error building response data json", 500)
		return
	}

	// finally we gonne send our response
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	_, err = w.Write(responseData)
	if nil != err {
		archivist.Error("Could not write http response body ", err, data)
	}
}

func getRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func buildListenConfigString() string {
	var connectString string
	connectString += config.GetValue("HOST")
	connectString += ":"
	connectString += config.GetValue("PORT")
	return connectString
}

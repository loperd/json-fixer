package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func main() {
	http.HandleFunc("/", handler)

	address := flag.String("address", "", "address to run the web server")
	port := flag.Int("port", 9900, "port to run the web server")
	flag.Parse()

	sprintf := fmt.Sprintf("%s:%d", *address, *port)

	fmt.Printf("Starting HTTP server on %s:%d\n", *address, *port)
	if err := http.ListenAndServe(sprintf, nil); err != nil {
		_ = fmt.Errorf(
			"Can not start http server on address: %s:%d\nError: %s",
			*address, *port, err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()

	if !query.Has("target") {
		response := &response{
			Error:   true,
			Message: "Invalid request",
			Detail:  "Doesn't passed target url.",
		}

		writeResponse(w, response)
		return
	}

	target, err := url.QueryUnescape(query.Get("target"))

	if err != nil || !isUrl(target) {
		response := &response{
			Error:   true,
			Message: "Invalid request",
			Detail:  "Invalid target passed.",
		}

		writeResponse(w, response)
		return
	}

	body, err := getContent(target, "GET")

	if err != nil {
		response := &response{
			Error:   true,
			Message: "Cannot get content.",
			Detail:  err.Error(),
		}

		writeResponse(w, response)
		return
	}

	content := clearContent(body)
	if !isJSON([]byte(content)) {
		response := &response{
			Error:   true,
			Message: "Invalid response.",
			Detail:  "Response is not json value.",
		}

		writeResponse(w, response)
		return
	}

	w.Write([]byte(content))
}

func clearContent(content string) string {
	r := regexp.MustCompile(`(,)([\s\n]*[\]\}])`)

	return r.ReplaceAllString(content, "$2")
}

func getContent(url string, method string) (string, error) {
	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(b), err
}

func isJSON(input []byte) bool {
	var decoded []interface{}
	return json.Unmarshal(input, &decoded) == nil

}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func writeResponse(w http.ResponseWriter, r *response) {
	jsonData, _ := json.Marshal(r)
	w.Write(jsonData)
}

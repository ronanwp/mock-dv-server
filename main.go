package main
// Mock DoubleVerify server. Based on the outline provided by the DoubleVerify IQ Connect guide.
// You can call the mock /dv-iqc endpoint with the params in the DoubleVerify preferred order. The hv parameter
// can be generated using the /hash endpoint with the same salt constant.
//
// Here is a Postman export of sample calls to either endpoint: https://www.postman.com/collections/b6ea82038b6249368c6f
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const salt = "0123456789"

func main() {
	http.HandleFunc("/dv-iqc", dvResponse)
	http.HandleFunc("/hash", computeHash)

	http.ListenAndServe(":8090", nil)
}

func dvResponse(w http.ResponseWriter, req *http.Request)  {
	rawParams := strings.Split(req.URL.RawQuery, "&")
	var lastParam string
	for _, val := range rawParams {
		lastParam = val
	}
	kevVal := strings.Split(lastParam, "=")
	if kevVal[0] != "hv"  {
		w.WriteHeader(http.StatusBadRequest) // DV Guide Page 33
		w.Write([]byte("400 - hv param not present in correct location"))
		return
	}
	values := req.URL.Query()
	url := values.Get("url")
	if url == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - url param not present"))
		return
	}
	ipaddress := values.Get("ip")
	if ipaddress == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - IP Address param not present"))
		return
	}
	partnerId := values.Get("partnerid")
	if partnerId == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Partner ID param not present"))
		return
	}
	userAgent := values.Get("useragent")
	if userAgent == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - User agent param not present"))
		return
	}
	hashVal := kevVal[1]
	// Check the hash value. See DV Guide Page 30
	compareStr := fmt.Sprintf("/dv-iqc?partnerid=%s&url=%s&useragent=%s&ip=%s%s", partnerId, url, userAgent, ipaddress, salt)
	h := sha256.Sum256([]byte(compareStr))
	encodedStr := hex.EncodeToString(h[:])
	if encodedStr != hashVal {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Incorrect parameter value for: hv"))
		return
	}
	// Everything is good. Return DV codes. Add to the switch below as required for testing.
	w.WriteHeader(http.StatusOK)
	switch{
	case url == "https://www.wattpad.com/story/5095707-after":
		w.Write([]byte("80023001,80312001,80013001"))
	case url == "https://www.wattpad.com/amp/248297765":
		w.Write([]byte("883032002,83032003,83032004"))
	}
}
// Mimic the hash function on DV Guide Page 31
func computeHash(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	data := values.Get("data")
	if data == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - data param not present"))
		return
	}
	salt := values.Get("salt")
	if salt == ""  {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - salt param not present"))
		return
	}
	h := sha256.Sum256([]byte(data+salt))
	encodedStr := hex.EncodeToString(h[:])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(encodedStr))
}

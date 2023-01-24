package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go/goserver.io/constants"
	"go/goserver.io/utils"
)

/**
 * Simple encryption using URL and Base64 encoding
 */
func Encode(clear string, key ...string) string {
	var enc strings.Builder
	key_s := GetSecret(key)
	for i, _ := range clear {
		key_c := key_s[i%len(key_s)]
		enc_c := rune(int(clear[i]) + int(key_c)%256)
		enc.WriteRune(enc_c)
	}
	return base64.URLEncoding.EncodeToString([]byte(enc.String()))
}

/**
 * Simple decryption using URL and Base64 encoding
 */
func Decode(enc string, key ...string) string {
	b64, err := base64.URLEncoding.DecodeString(enc)
	if err != nil {
		return utils.LogError(err)
	}

	u8 := []rune(string(b64))
	key_s := GetSecret(key)
	var dec strings.Builder
	for i, _ := range u8 {
		key_c := key_s[i%len(key_s)]
		enc_c := rune((256 + int(u8[i]) - int(key_c)) % 256)
		dec.WriteRune(enc_c)
	}
	return dec.String()
}

/**
 * Return encryption key
 * See: constants.Constants.SECRET_KEY
 */
func GetSecret(arguments []string) string {
	if len(arguments) == 0 {
		return constants.SECRET_KEY
	}
	return arguments[0]
}

/**
 * Proxy request function between UI/Portal and API/Server
 */
func ApiRequest(w http.ResponseWriter, r *http.Request, path string, req_body ...map[string]interface{}) (*constants.Response, error) {
	if len(req_body) > 0 {
		return ApiRequestImpl(w, r, path, true, req_body[0])
	} else {
		return ApiRequestImpl(w, r, path, true, nil)
	}
}

func ApiRequestImpl(w http.ResponseWriter, r *http.Request, path string, isLogBody bool, req_body map[string]interface{}) (*constants.Response, error) {
	var headers map[string]interface{}
	var io_body io.Reader
	method := "GET"
	url := constants.API_BASEPATH + path
	if req_body != nil {
		method = "POST"
		b, err := json.Marshal(req_body)
		if err != nil {
			utils.LogError(err)
			return nil, err
		}
		io_body = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, io_body)
	req.Header.Set("APIKEY", constants.API_KEY)
	for _, val := range constants.HEADERS_REQUIRED {
		req.Header.Set(val, r.Header.Get(val))
	}

	log.Printf("ApiRequest %s: %s", method, url)

	if req_body != nil && isLogBody {
		log.Printf("ApiRequest Body: %s", req_body)
	}
	if headers != nil {
		for key, val := range headers {
			req.Header.Set(key, fmt.Sprintf("%v", val))
		}
	}
	log.Printf("ApiRequest Headers: %s\n", req.Header)

	client := &http.Client{
		Timeout: time.Duration(constants.HTTP_TIMEOUT * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		utils.LogError(err)
		return nil, errors.New(constants.DEFAULT_HTTP_ERROR)
	}
	defer res.Body.Close()

	if w != nil {
		for name, values := range res.Header {
			w.Header()[name] = values
		}
		w.WriteHeader(res.StatusCode)
		_, err = io.Copy(w, res.Body)
		if err != nil {
			utils.LogError(err)
		}
		return nil, err
	}

	var i interface{}
	err = json.NewDecoder(res.Body).Decode(&i)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	body := i.(map[string]interface{})
	resp := &constants.Response{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       body,
	}
	msg, found1 := body["error"]
	_, found2 := body["code"]
	if found1 && found2 {
		resp.Decode(&resp.Error)
		resp.Error.Message = msg.(string)
		utils.LogError(resp.Error)
	}
	return resp, nil
}

/**
 * Remove HTML tags from string
 */
func FormValue(r *http.Request, key string) string {
	return constants.REMOVE_TAGS_REGEXP.ReplaceAllString(r.FormValue(key), "")
}

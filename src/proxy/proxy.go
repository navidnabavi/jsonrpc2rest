package Proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

//JSONRPCRequest stores Request of JSONRPC to unmarshal
type JSONRPCRequest struct {
	JSONrpc string        `json:"jsonrpc"` //version JSONRPC protocol
	Method  string        `json:"method"`  //method name of rpc
	Params  []interface{} `json:"params"`  //parameters of JSONRPC in array form
	ID      int64         `json:"id"`      //ID of request
}

//JSONRPCError stores error data of JSONRPC
type JSONRPCError struct {
	Code    int                    `json:"code"`    //error code of json rpc
	Message string                 `json:"message"` // error message
	Data    map[string]interface{} `json:"data"`    //data of stacktrace or details of error
}

//JSONRPCResponse stores response of JSONRPC to marshal
type JSONRPCResponse struct {
	Error  *JSONRPCError `json:"error"`  //Error object of JSONRPC
	Result interface{}   `json:"result"` // Result of RPC
	ID     int64         `json:"id"`     // ID of request
}

//ProxyManager handles  jsonrpc requests and manipulate data to rest mode
type ProxyManager struct {
	config *ConfigManager.Configuration //Configuration from json file
	client *http.Client                 //http client to handle requests
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//resolve matches the weighted server to proxy pass
func (p *ProxyManager) resolve(hostname string) (*string, bool) {

	value, ok := p.config.Hosts[hostname]
	if ok {
		sum := 0
		random := rand.Intn(sum)

		for _, v := range value {
			sum += v
		}

		sum = 0
		for k, v := range value {
			if sum <= random && sum+v > random {
				return &k, ok
			}
			sum += v
		}
	}

	return nil, ok
}

//addSlashIfMissed adds slash to the end of url
func addSlashIfMissed(_url *string) {
	if (*_url)[len(*_url)-1] != '/' {
		*_url += "/"
	}
}

//makeRequestParams creates url and payload for request
func (p *ProxyManager) makeRequestParams(upstream *ConfigManager.Upstream, params []interface{}) (string, map[string]interface{}) {

	_url := upstream.URL
	query := "?"
	addSlashIfMissed(&_url)

	payloadParams := make(map[string]interface{})

	for i, p := range upstream.Params {
		switch upstream.ParamTypes[i] {
		case ConfigManager.ParamTypeURL:
			_url = strings.Replace(_url, ":"+p, fmt.Sprint(params[i]), -1)
		case ConfigManager.ParamTypePayload:
			payloadParams[p] = params[i]
		case ConfigManager.ParamTypeQuery:
			urlEncodedValue := url.PathEscape(fmt.Sprintf("%v", params[i]))
			query += fmt.Sprintf("%s=%s&", p, urlEncodedValue)
		}
	}
	if len(query) > 1 {
		_url += query
	}

	return _url, payloadParams
}

//proxyPass passes a json rpc request to rest service
func (p *ProxyManager) proxyPass(method string, params []interface{}, header *http.Header) (*io.ReadCloser, int, *http.Header) {
	// var buffer bytes.Buffer
	upstream := p.config.Upstreams[method]

	_url, payloadParams := p.makeRequestParams(upstream, params)

	httpMethod := upstream.Method
	payload, err := EncodePayload(payloadParams, 0)

	check(err)

	urlParsed, _ := url.Parse(_url)

	host, ok := p.resolve(urlParsed.Hostname())

	if ok {
		urlParsed.Host = *host
	}

	fmt.Printf("url is %s", urlParsed.String())
	req, err := http.NewRequest(httpMethod, urlParsed.String(), bytes.NewReader(payload))
	req.Header = *header
	req.Header.Set("CONTENT_TYPE", "application/json")

	check(err)

	res, err := p.client.Do(req)

	check(err)

	return &res.Body, res.StatusCode, &res.Header
}

func (p *ProxyManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var msg JSONRPCRequest
	var responseObj map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)

	check(err)

	fmt.Println(msg)
	fmt.Println(msg.Method)

	response, statusCode, _ := p.proxyPass(msg.Method, msg.Params, &r.Header)

	defer (*response).Close()

	w.WriteHeader(statusCode)

	responseBytes, err := ioutil.ReadAll(*response)

	check(err)

	json.Unmarshal(responseBytes, &responseObj)

	jsonRPCResult := JSONRPCResponse{ID: msg.ID, Result: responseObj["result"], Error: nil}
	jsonRPCResponseBytes, err := json.Marshal(jsonRPCResult)

	check(err)

	w.Write(jsonRPCResponseBytes)
}

//Serve starts to handle requests
func (p *ProxyManager) Serve() {
	http.Handle("/json/", p)
	log.Fatal(http.ListenAndServe(p.config.Bind, nil))
}

//NewProxy constructs a ProxyManager
func NewProxy(config *ConfigManager.Configuration) ProxyManager {
	return ProxyManager{
		config: config,
		client: &http.Client{}}
}

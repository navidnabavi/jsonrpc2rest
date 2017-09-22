package ConfigManager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	//ParamTypeURL represents a param which appears in url
	ParamTypeURL int8 = 0
	//ParamTypePayload represents a param which appears in payload
	ParamTypePayload int8 = 1
	//ParamTypeQuery represents a param which appears in query string in url
	ParamTypeQuery int8 = 2
)

//Upstream contains the config of upstream
type Upstream struct {
	URL           string   `json:"url"`
	Method        string   `json:"method"`
	Params        []string `json:"params"`
	PayloadParams []string `json:"payload_params"`
	ParamTypes    []int8
	ParamCount    []uint8
}

//Configuration contains the configuration
type Configuration struct {
	Hosts     map[string]map[string]int `json:"hosts"`     // host aliases and weights
	Upstreams map[string]*Upstream      `json:"upstreams"` //list of upstreams
	Bind      string                    `json:"bind"`      //bind address of server
}

func isInPayload(param string, params []string) bool {
	for _, v := range params {
		if v == param {
			return true
		}
	}
	return false
}

func (u *Upstream) initialCounts() {
	u.ParamCount = make([]uint8, 3)
	u.ParamCount[ParamTypeURL] = 0
	u.ParamCount[ParamTypePayload] = 0
	u.ParamCount[ParamTypeQuery] = 0
}

//Comment it later
func (u *Upstream) computeAddtionalConfig() {
	u.ParamTypes = make([]int8, len(u.Params))
	u.initialCounts()

	for i, v := range u.Params {
		fmt.Println(v)
		if strings.Contains(u.URL, ":"+v) {
			if isInPayload(v, u.PayloadParams) {
				panic(fmt.Sprintf("Parameter %s defined in both url and payload params", v))
			}
			u.ParamTypes[i] = ParamTypeURL
			u.ParamCount[ParamTypeURL]++
		} else {
			isPayloadParams := false
			for _, p := range u.PayloadParams {
				if p == v {
					isPayloadParams = true
				}
			}
			if !isPayloadParams {
				u.ParamTypes[i] = ParamTypeQuery
				u.ParamCount[ParamTypeQuery]++
			} else {
				u.ParamTypes[i] = ParamTypePayload
				u.ParamCount[ParamTypePayload]++
			}
		}
	}
}

//LoadConfiguration loads json data to Configuration struct
func LoadConfiguration(filename string) Configuration {
	var config Configuration
	data, err := ioutil.ReadFile(filename)
	check(err)
	err = json.Unmarshal(data, &config)
	check(err)

	for k := range config.Upstreams {
		config.Upstreams[k].computeAddtionalConfig()
	}
	fmt.Println("-----")
	fmt.Println(config)
	return config
}

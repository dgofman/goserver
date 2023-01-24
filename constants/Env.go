package constants

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"

	"go/goserver.io/utils"
)

/**
 * Environment properties defined in properties file or passing via command line arguments
 */

var (
	HEADERS_REQUIRED = []string{"X-Auth-Subject", "X-Auth-User-Email", "X-Auth-User-Fn", "X-Auth-User-Ln", "X-Auth-User-Name"}

	REMOVE_TAGS_REGEXP = regexp.MustCompile("<[^>]*>|/>")

	Props        = *initProps()
	SECRET_TOKEN = utils.UUID()

	API_BASEPATH = utils.Get(Props.Envs, Props.ENV+"_API_BASEPATH", "").(string)
	API_KEY      = utils.Get(Props.Envs, Props.ENV+"_API_KEY", "").(string)

	Whitelist_urls     = regexp.MustCompile(fmt.Sprintf(`^%s|^%s|^%s`, URL_Static, URL_Login, URL_Logout))
	Whitelist_reg_urls = regexp.MustCompile(`^/favicon.ico|^/media/|^/admin/|^/users/|^/public/`)
)

var instance *Properties
var once sync.Once

func initProps() *Properties {
	once.Do(func() {
		instance = &Properties{}
		f, err := os.Open("properties.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(instance)
		if err != nil {
			panic(err)
		}
		flag.IntVar(&instance.PORT, "port", instance.PORT, "Port Number")
		flag.BoolVar(&instance.LOGFILE, "file", instance.LOGFILE, "Output to log file")
		flag.BoolVar(&instance.CONSOLE, "console", instance.CONSOLE, "Output to console")
		flag.StringVar(&instance.VIEW, "view", instance.VIEW, "'full' - show header and footer otherwise use 'light'")
		flag.StringVar(&instance.ENV, "env", instance.ENV, "API environment name")
		flag.Parse()

		if instance.VIEW != "light" && instance.VIEW != "full" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
	return instance
}

/**
 * Parse response body
 * @param key  - name in the dictionary
 * @param defVal - default value if key is not found
 * @return - value stored in the response BODY or NIL
 */
func (resp *Response) Get(key string, defVal ...interface{}) interface{} {
	return utils.Get(resp.Body, key, defVal)
}

/**
 * Decode reads the response body value from its
 * input and stores it in the value pointed to by v.
 */
func (resp *Response) Decode(v interface{}) error {
	b, err := json.Marshal(resp.Body)
	if err != nil {
		utils.LogError(err)
		return err
	}
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&v)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

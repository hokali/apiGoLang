package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// OktaIntrospectPayload response fro mthe SPA ... now we get the JWT from Okta
type OktaIntrospectPayload struct {
	Location string `json:"Location"`
	Token    string `json:"token"`
}

// Introspect response from mthe SPA ...
type Introspect struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	Username  string `json:"username"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Sub       string `json:"sub"`
	Aud       string `json:"aud"`
	Iss       string `json:"iss"`
	Jti       string `json:"jti"`
	TokenType string `json:"token_type"`
	ClientID  string `json:"client_id"`
	UID       string `json:"uid"`
	Mail      string `json:"mail"`
}

// OktaToken given
type OktaToken struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	IDToken     string `json:"id_token"`
}

// CheckToken used to see if token from okta is active
func CheckToken(w http.ResponseWriter, r *http.Request) {
	// Declare a new Introspect Paload Struct.
	introspect := OktaIntrospectPayload{}
	err := json.NewDecoder(r.Body).Decode(&introspect)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	var url string

	url = chkLocation(introspect.Location)

	tokenResource := "/v1/introspect"
	introspectJSON := introspectCheck((url + tokenResource), introspect)

	var resp = map[string]interface{}{}
	resp["introspect"] = introspectJSON
	_ = json.NewEncoder(w).Encode(resp)
}

// JwtMiddlewareChk used to see if token from okta is active
func JwtMiddlewareChk(w http.ResponseWriter, r *http.Request) bool {
	// Declare a new Introspect Paload Struct.
	introspect := OktaIntrospectPayload{}
	err := json.NewDecoder(r.Body).Decode(&introspect)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		_ = json.NewEncoder(w).Encode(resp)
		return false
	}

	var url string

	url = chkLocation(introspect.Location)

	tokenResource := "/v1/introspect"
	introspectJSON := introspectCheck((url + tokenResource), introspect)

	return introspectJSON.Active
}

// Check Location resource
func chkLocation(location string) string {
	if location != "Prod" {
		return os.Getenv("OKTA_STAGED")
	} else {
		return os.Getenv("OKTA_PROD")
	}
}

// introspectCheck get the token status from okta
func introspectCheck(apiURL string, o OktaIntrospectPayload) Introspect {
	var username, passwd, payloadtxt, basic string
	client := &http.Client{}

	if o.Location != "Prod" {
		username = os.Getenv("OKTA_STAGED_USERNAME")
		passwd = os.Getenv("OKTA_STAGED_PASSWD")
		basic = os.Getenv("OKTA_STAGED_BASIC")
	} else {
		username = os.Getenv("OKTA_PROD_USERNAME")
		passwd = os.Getenv("OKTA_PROD_PASSWD")
		basic = os.Getenv("OKTA_PROD_BASIC")
	}
	payloadtxt = "token=" + o.Token + "&token_type_hint=access_token"
	payload := strings.NewReader(payloadtxt)
	req, err := http.NewRequest("POST", apiURL, payload)
	if err != nil {
		fmt.Println(err)
	}

	req.SetBasicAuth(username, passwd)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", basic)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	oktaJSON := Introspect{}
	_ = json.Unmarshal([]byte(bodyText), &oktaJSON)

	return oktaJSON
}

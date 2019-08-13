package hydra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BasicLoginHandler struct {
	hydraURI string
	client   *http.Client
}

func (b BasicLoginHandler) LoginRequest(challenge string) map[string]interface{} {
	return get(challenge, "login", b.hydraURI, b.client)
}

func (b BasicLoginHandler) ConsentRequest(challenge string) map[string]interface{} {
	return get(challenge, "consent", b.hydraURI, b.client)
}

func (b BasicLoginHandler) AcceptLogin(challenge string, body map[string]interface{}) string {
	return getRedirect(put(challenge, "login", "accept", body, b.hydraURI, b.client))
}

func (b BasicLoginHandler) DenyLogin(challenge string, body map[string]interface{}) string {
	return getRedirect(put(challenge, "login", "deny", body, b.hydraURI, b.client))
}

func (b BasicLoginHandler) AcceptConsent(challenge string, body map[string]interface{}) string {
	return getRedirect(put(challenge, "consent", "accept", body, b.hydraURI, b.client))
}

func (b BasicLoginHandler) DenyConsent(challenge string, body map[string]interface{}) string {
	return getRedirect(put(challenge, "consent", "deny", body, b.hydraURI, b.client))
}

func getRedirect(resp map[string]interface{}) string {
	r := resp["redirect_to"]
	switch v := r.(type) {
	case string:
		if v != "" {
			return v
		}
		log.Printf("response: %v\n", resp)
		log.Fatalln("Invalid redirect uri recieved")
	default:
		log.Printf("response: %v\n", resp)
		log.Fatalf("invalid redirect")
	}
	return ""
}

func get(challenge string, flow string, uri string, client *http.Client) map[string]interface{} {
	requestURL := fmt.Sprintf("%s/oauth2/auth/requests/%s/%s", uri, flow, challenge)
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatalf("unable to create request: %v\n", err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("unable to get response: %v\n", err)
	}
	respMap := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&respMap)

	if err != nil {
		log.Fatal(err.Error())
	}
	return respMap
}

func put(challenge string, flow string, action string, body map[string]interface{}, uri string, client *http.Client) map[string]interface{} {
	requestURL := fmt.Sprintf("%s/oauth2/auth/requests/%s/%s", uri, flow, challenge)
	s, err := json.Marshal(body)
	by := bytes.NewBuffer(s)
	request, err := http.NewRequest("PUT", requestURL, by)

	if err != nil {
		log.Fatalf("unable to create request: %v\n", err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("unable to get response: %v\n", err)
	}
	respMap := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&respMap)

	if err != nil {
		log.Fatal(err.Error())
	}
	return respMap
}

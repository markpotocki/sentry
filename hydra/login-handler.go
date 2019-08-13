package hydra

type LoginHandler interface {
	LoginRequest(challenge string) map[string]interface{}
	ConsentRequest(challenge string) map[string]interface{}
	AcceptLogin(challenge string, body map[string]interface{}) string
	DenyLogin(challenge string, body map[string]interface{}) string
	AcceptConsent(challenge string, body map[string]interface{}) string
	DenyConsent(challenge string, body map[string]interface{}) string
}

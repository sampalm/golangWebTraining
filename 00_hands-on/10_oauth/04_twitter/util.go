package oauthweb

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

// EncodeParams encodes the given string according to RFC 3986.
func EncodeParams(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

// EncodeSignature encodes the request header and returns a oauth_signature.
func EncodeSignature(method string, requestURI string, auth *oAuthV4) (string, url.Values) {
	t := time.Now().Unix()
	values := make(url.Values)
	values.Add("oauth_consumer_key", auth.ClientID)
	values.Add("oauth_nonce", auth.ANonce)
	values.Add("oauth_signature_method", "HMAC-SHA1")
	values.Add("oauth_timestamp", fmt.Sprintf("%v", t))
	values.Add("oauth_token", auth.OTokenID)
	values.Add("oauth_version", "1.0")
	params := strings.Replace(values.Encode(), "&", "%26", -1)
	params = strings.Replace(params, "=", "%3D", -1)
	params = fmt.Sprintf("%s&%s&%s", method, EncodeParams(requestURI), strings.Replace(params, "+", "%20", -1))
	log.Printf("EncodeSignature BODY: %s\n", params)
	// GET SIGNIN KEY
	signKEY := fmt.Sprintf("%s&%s", EncodeParams(auth.SecretID), EncodeParams(auth.OTokenSecretID))
	// GENERATE HASH
	mac := hmac.New(sha1.New, []byte(signKEY))
	mac.Write([]byte(params))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, values
}

// EncodeNewHeader returns a new encoded header to authorization requests.
func EncodeNewHeader(signature string, values url.Values) string {
	return fmt.Sprintf("OAuth oauth_consumer_key=%s, oauth_nonce=%s, oauth_signature=%s, oauth_signature_method=%s, oauth_timestamp=%s, oauth_token=%s, oauth_version=%s", values.Get("oauth_consumer_key"), values.Get("oauth_nonce"), EncodeParams(signature), values.Get("oauth_signature_method"), values.Get("oauth_timestamp"), values.Get("oauth_token"), values.Get("oauth_version"))
}

// EncodeNewHeaderHTTP returns a new encoded HTTP header to authorization requests.
func EncodeNewHeaderHTTP(requestURI string, signature string, values url.Values) string {
	return fmt.Sprintf("%s?oauth_consumer_key=%s&oauth_nonce=%s&oauth_signature=%s&oauth_signature_method=%s&oauth_timestamp=%s&oauth_token=%s&oauth_version=%s", requestURI, values.Get("oauth_consumer_key"), values.Get("oauth_nonce"), EncodeParams(signature), values.Get("oauth_signature_method"), values.Get("oauth_timestamp"), values.Get("oauth_token"), values.Get("oauth_version"))
}

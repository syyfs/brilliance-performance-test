package chaincodefunction

import "testing"

func TestHttpsPost(t *testing.T) {
	caCertPath := "./keystore/unionbank.crt"
	HttpsPost(caCertPath)
}

package passwordcrypto

import "testing"

const testKey = "0123456789abcdef0123456789abcdef"

func TestRoundTripAndRandomNonce(t *testing.T) {
	c, err := New(testKey)
	if err != nil {
		t.Fatal(err)
	}
	one, err := c.Encrypt("example-password")
	if err != nil {
		t.Fatal(err)
	}
	two, err := c.Encrypt("example-password")
	if err != nil {
		t.Fatal(err)
	}
	if one == two {
		t.Fatal("encryptions must differ because each uses a random nonce")
	}
	if !c.Matches(one, "example-password") || c.Matches(one, "wrong-password") {
		t.Fatal("password matching returned an unexpected result")
	}
}

func TestRejectsInvalidKeyAndCiphertext(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Fatal("expected empty key error")
	}
	c, _ := New(testKey)
	if _, err := c.Decrypt("v1:not-base64!"); err == nil {
		t.Fatal("expected invalid ciphertext error")
	}
}

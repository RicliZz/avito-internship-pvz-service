package pass

import "testing"

func TestComparePassWithHash(t *testing.T) {
	my_password := "mysecurepassword"
	hash_my_password, err := CreateHash(my_password)
	if err != nil {
		t.Error(err)
	}
	if ComparePassWithHash("WRONG", hash_my_password) {
		t.Error("ComparePassWithHash fail")
	}
	if !ComparePassWithHash(my_password, hash_my_password) {
		t.Error("expected password not to match hash")
	}
}

func TestCreateHash(t *testing.T) {
	password := "mysecurepassword"

	hash, err := CreateHash(password)
	if err != nil {
		t.Fatalf("failed to create hash: %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	match := ComparePassWithHash(password, hash)
	if !match {
		t.Error("expected password to match generated hash")
	}
}

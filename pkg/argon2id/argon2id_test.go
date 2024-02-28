package argon2id

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestCreateHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultParams)
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		t.Fatalf("expected 6 parts, got %d", len(parts))
	}

	if parts[1] != "argon2id" {
		t.Errorf("expected variant argon2id, got %s", parts[1])
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		t.Fatal(err)
	}

	if version != argon2.Version {
		t.Errorf("expected version %d, got %d", argon2.Version, version)
	}

	var m, i uint32
	var p uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &i, &p)
	if err != nil {
		t.Fatal(err)
	}

	if m != DefaultParams.Memory {
		t.Errorf("expected memory %d, got %d", DefaultParams.Memory, m)
	}
	if i != DefaultParams.Iterations {
		t.Errorf("expected iterations %d, got %d", DefaultParams.Iterations, i)
	}
	if p != DefaultParams.Parallelism {
		t.Errorf("expected parallelism %d, got %d", DefaultParams.Parallelism, p)
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		t.Fatal(err)
	}

	if uint32(len(salt)) != DefaultParams.SaltLength {
		t.Errorf("expected salt length %d, got %d", DefaultParams.SaltLength, len(salt))
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultParams)
	if err != nil {
		t.Fatal(err)
	}

	match, err := ComparePasswordAndHash("pa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}

	if !match {
		t.Error("expected password and hash to match")
	}

	match, err = ComparePasswordAndHash("otherPa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}

	if match {
		t.Error("expected password and hash to not match")
	}
}

func TestDecodeHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultParams)
	if err != nil {
		t.Fatal(err)
	}

	params, _, _, err := DecodeHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	if *params != *DefaultParams {
		t.Fatalf("expected %#v got %#v", *DefaultParams, *params)
	}
}

func TestCheckHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultParams)
	if err != nil {
		t.Fatal(err)
	}

	ok, params, err := CheckHash("pa$$word", hash)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected password to match")
	}
	if *params != *DefaultParams {
		t.Fatalf("expected %#v got %#v", *DefaultParams, *params)
	}
}

func TestVariant(t *testing.T) {
	// Hash contains wrong variant
	_, _, err := CheckHash("pa$$word", "$argon2i$v=19$m=32768,t=2,p=2$OciL0ybwe8Cw4LHQ/zfBEg$lZrRjdNLsc11ileWikeFrt0ULA5ZwirvU+OyXQ1c4Hw")
	if err != ErrIncompatibleVariant {
		t.Fatalf("expected error %s", ErrIncompatibleVariant)
	}
}

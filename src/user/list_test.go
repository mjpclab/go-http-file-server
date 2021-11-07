package user

import (
	"testing"
)

var list *List

func init() {
	list = NewList()
}

func TestUserPlain(t *testing.T) {
	list.AddPlain("plain_user", "123")
	if !list.Auth("plain_user", "123") {
		t.Error()
	}
	if list.Auth("plain_user", "12") {
		t.Error()
	}
}

func TestUserBase64(t *testing.T) {
	list.AddBase64("base64_user", "MjM0")
	if !list.Auth("base64_user", "234") {
		t.Error()
	}
	if list.Auth("base64_user", "23") {
		t.Error()
	}
}

func TestUserMd5(t *testing.T) {
	list.AddMd5("md5_user", "d81f9c1be2e08964bf9f24b15f0e4900")
	if !list.Auth("md5_user", "345") {
		t.Error()
	}
	if list.Auth("md5_user", "34") {
		t.Error()
	}
}

func TestUserSha1(t *testing.T) {
	list.AddSha1("sha1_user", "51eac6b471a284d3341d8c0c63d0f1a286262a18")
	if !list.Auth("sha1_user", "456") {
		t.Error()
	}
	if list.Auth("sha1_user", "45") {
		t.Error()
	}
}

func TestUserSha256(t *testing.T) {
	list.AddSha256("sha256_user", "97a6d21df7c51e8289ac1a8c026aaac143e15aa1957f54f42e30d8f8a85c3a55")
	if !list.Auth("sha256_user", "567") {
		t.Error()
	}
	if list.Auth("sha256_user", "56") {
		t.Error()
	}
}

func TestUserSha512(t *testing.T) {
	list.AddSha512("sha512_user", "c7d57e5c0b0792b154d573089792d80f5b64d2bc0cf4d7d1f551a9e4a4966e925d06b253cc9662c01df76623fdfecb812a2a0604119cb1ac37c47e8027e94cb5")
	if !list.Auth("sha512_user", "678") {
		t.Error()
	}
	if list.Auth("sha512_user", "67") {
		t.Error()
	}
}

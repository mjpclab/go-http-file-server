package user

import (
	"testing"
)

func TestUserPlain(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "plain_user"
	list.AddPlain(username, "123")
	if id, u, ok = list.Auth("plain_user", "123"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Plain_User", "123"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("plain_user", "12"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "123"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

func TestUserBase64(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "base64_user"
	list.AddBase64(username, "MjM0")
	if id, u, ok = list.Auth("base64_user", "234"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Base64_user", "234"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("base64_user", "23"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "MjM0"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

func TestUserMd5(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "md5_user"
	list.AddMd5(username, "d81f9c1be2e08964bf9f24b15f0e4900")
	if id, u, ok = list.Auth("md5_user", "345"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Md5_user", "345"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("md5_user", "34"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "d81f9c1be2e08964bf9f24b15f0e4900"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

func TestUserSha1(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "sha1_user"
	list.AddSha1(username, "51eac6b471a284d3341d8c0c63d0f1a286262a18")
	if id, u, ok = list.Auth("sha1_user", "456"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Sha1_user", "456"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("sha1_user", "45"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "51eac6b471a284d3341d8c0c63d0f1a286262a18"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

func TestUserSha256(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "sha256_user"
	list.AddSha256(username, "97a6d21df7c51e8289ac1a8c026aaac143e15aa1957f54f42e30d8f8a85c3a55")
	if id, u, ok = list.Auth("sha256_user", "567"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Sha256_user", "567"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("sha256_user", "56"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "97a6d21df7c51e8289ac1a8c026aaac143e15aa1957f54f42e30d8f8a85c3a55"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

func TestUserSha512(t *testing.T) {
	var list = NewList()
	var id int
	var u string
	var ok bool

	username := "sha512_user"
	list.AddSha512(username, "c7d57e5c0b0792b154d573089792d80f5b64d2bc0cf4d7d1f551a9e4a4966e925d06b253cc9662c01df76623fdfecb812a2a0604119cb1ac37c47e8027e94cb5")
	if id, u, ok = list.Auth("sha512_user", "678"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("Sha512_user", "678"); id != 0 || u != username || !ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("sha512_user", "67"); id != 0 || u != username || ok {
		t.Error(id, u, ok)
	}
	if id, u, ok = list.Auth("fake_user", "c7d57e5c0b0792b154d573089792d80f5b64d2bc0cf4d7d1f551a9e4a4966e925d06b253cc9662c01df76623fdfecb812a2a0604119cb1ac37c47e8027e94cb5"); id != -1 || u != "" || ok {
		t.Error(id, u, ok)
	}
}

package user

import (
	"testing"
)

var listAccurate *List
var listNoCase *List

func init() {
	listAccurate = NewList(true)
	listNoCase = NewList(false)
}

func TestUserPlain(t *testing.T) {
	// accurate
	listAccurate.AddPlain("plain_user", "123")
	if !listAccurate.Auth("plain_user", "123") {
		t.Error()
	}
	if listAccurate.Auth("Plain_user", "123") {
		t.Error()
	}
	if listAccurate.Auth("plain_user", "12") {
		t.Error()
	}

	// no case
	listNoCase.AddPlain("plain_user", "123")
	if !listNoCase.Auth("plain_user", "123") {
		t.Error()
	}
	if !listNoCase.Auth("Plain_User", "123") {
		t.Error()
	}
	if listNoCase.Auth("plain_user", "12") {
		t.Error()
	}
}

func TestUserBase64(t *testing.T) {
	// accurate
	listAccurate.AddBase64("base64_user", "MjM0")
	if !listAccurate.Auth("base64_user", "234") {
		t.Error()
	}
	if listAccurate.Auth("Base64_user", "234") {
		t.Error()
	}
	if listAccurate.Auth("base64_user", "23") {
		t.Error()
	}

	// no case
	listNoCase.AddBase64("base64_user", "MjM0")
	if !listNoCase.Auth("base64_user", "234") {
		t.Error()
	}
	if !listNoCase.Auth("Base64_user", "234") {
		t.Error()
	}
	if listNoCase.Auth("base64_user", "23") {
		t.Error()
	}

}

func TestUserMd5(t *testing.T) {
	// accurate
	listAccurate.AddMd5("md5_user", "d81f9c1be2e08964bf9f24b15f0e4900")
	if !listAccurate.Auth("md5_user", "345") {
		t.Error()
	}
	if listAccurate.Auth("Md5_user", "345") {
		t.Error()
	}
	if listAccurate.Auth("md5_user", "34") {
		t.Error()
	}

	// no case
	listNoCase.AddMd5("md5_user", "d81f9c1be2e08964bf9f24b15f0e4900")
	if !listNoCase.Auth("md5_user", "345") {
		t.Error()
	}
	if !listNoCase.Auth("Md5_user", "345") {
		t.Error()
	}
	if listNoCase.Auth("md5_user", "34") {
		t.Error()
	}

}

func TestUserSha1(t *testing.T) {
	// accurate
	listAccurate.AddSha1("sha1_user", "51eac6b471a284d3341d8c0c63d0f1a286262a18")
	if !listAccurate.Auth("sha1_user", "456") {
		t.Error()
	}
	if listAccurate.Auth("Sha1_user", "456") {
		t.Error()
	}
	if listAccurate.Auth("sha1_user", "45") {
		t.Error()
	}

	// no case
	listNoCase.AddSha1("sha1_user", "51eac6b471a284d3341d8c0c63d0f1a286262a18")
	if !listNoCase.Auth("sha1_user", "456") {
		t.Error()
	}
	if !listNoCase.Auth("Sha1_user", "456") {
		t.Error()
	}
	if listNoCase.Auth("sha1_user", "45") {
		t.Error()
	}
}

func TestUserSha256(t *testing.T) {
	// accurate
	listAccurate.AddSha256("sha256_user", "97a6d21df7c51e8289ac1a8c026aaac143e15aa1957f54f42e30d8f8a85c3a55")
	if !listAccurate.Auth("sha256_user", "567") {
		t.Error()
	}
	if listAccurate.Auth("Sha256_user", "567") {
		t.Error()
	}
	if listAccurate.Auth("sha256_user", "56") {
		t.Error()
	}

	// no case
	listNoCase.AddSha256("sha256_user", "97a6d21df7c51e8289ac1a8c026aaac143e15aa1957f54f42e30d8f8a85c3a55")
	if !listNoCase.Auth("sha256_user", "567") {
		t.Error()
	}
	if !listNoCase.Auth("Sha256_user", "567") {
		t.Error()
	}
	if listNoCase.Auth("sha256_user", "56") {
		t.Error()
	}
}

func TestUserSha512(t *testing.T) {
	// accurate
	listAccurate.AddSha512("sha512_user", "c7d57e5c0b0792b154d573089792d80f5b64d2bc0cf4d7d1f551a9e4a4966e925d06b253cc9662c01df76623fdfecb812a2a0604119cb1ac37c47e8027e94cb5")
	if !listAccurate.Auth("sha512_user", "678") {
		t.Error()
	}
	if listAccurate.Auth("Sha512_user", "678") {
		t.Error()
	}
	if listAccurate.Auth("sha512_user", "67") {
		t.Error()
	}

	// no case
	listNoCase.AddSha512("sha512_user", "c7d57e5c0b0792b154d573089792d80f5b64d2bc0cf4d7d1f551a9e4a4966e925d06b253cc9662c01df76623fdfecb812a2a0604119cb1ac37c47e8027e94cb5")
	if !listNoCase.Auth("sha512_user", "678") {
		t.Error()
	}
	if !listNoCase.Auth("Sha512_user", "678") {
		t.Error()
	}
	if listNoCase.Auth("sha512_user", "67") {
		t.Error()
	}
}

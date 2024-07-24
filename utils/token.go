package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuBTMMu9kdFcQsNhYvIYjaMzxUIMmQBar5FF0PRxPQpOzJF2W
PPh1p6XgNyWhjzFjNOfY+KS1mdwSP4HVh+29IZjqufhhhFW4/4Knbl2zFG6VcVsV
E+hevxm6UYDC+latXowoEZ4CYzLqCve0OwlPx0x27VIMNpSEC88U5TfCCXU7g/yr
1mhsv6+sTl7fVwF7GrRCqTuNZjMrdtdUaZo1Ipp+2PB3uJoM8NOcRsTZ87vrLwlB
KHDfoYeOguxf5rac+3nw44vXjynT0SG86ij/ahppGe08zKKrEN/whn8PGpVj6plf
zO3vElFMDhKzrlnWFsO0gyu3+Lp1j7ZWdfN7DQIDAQABAoIBACueOfiDwxDSJJFS
4kvhmPhXP4LtYJ3lAVv7dQiZ8an754Hhbx0JXo+X/XXzw8FzWWeK3F2mYaWQgP5t
XFR2H8+bNVtVoH3D9i8NHXFIYIuh3GVcyLTL5c9wYDa5xPmemjwCB+iMwVISGWT5
5snqXe4Cj1eSjnRMYrGNowzhVmf46Vc/TfZfqcOo2tJAEMx1iBEusLkl727bfLHo
EBsJAMUdEHa/v9TU+ClTLDSyr8Wt6dpFxoi/51lRertIjubxDzPnaVZWbrEi1+qt
kDTdQs8oI7fA9ICTqtkIIfAbkzRwRKH5BgrtGzX+y4phytRnWeGuGWWJHgwIMHpH
+DTEpR0CgYEA5SaEAVSKI7+/3mBXfInMpqpjwnAjeZYWTaqbAOVXdeoOodLut4gJ
i+JZ+5sQs35KzUsSFgr5xe801PR/7rJENxU5lIrw0WSlo+1BWFSlR5y6ZsGr2fWV
xzQXPZ/K3ExzZmSM4YKz2tFbi1C+2QDApq9EdhjYYPveesYMApwOEssCgYEAzaZo
eVp9uoMZljD6Otld9+t2qnYDR6VmlFlvpQ0eqfP0LDjLY6G5BmrrSPArO3BBEeLZ
wwjq14mXQNrgqRC1GFpt+LaTUzFEOSmq1Yi/SOHiMg4Pf1WWrv/EHc3mzZoQf447
LezNY62dSLUQHNjJbtKaLfPak1JIOuhi9GfadocCgYEAoO04nVqKnOqHy5srNZns
sEtPPfjU4QmHZknfC3UExBl45yqkXR3bXnK7MNjIlNWnoJ8M95ADs373QmrnAXIO
OATe6DPfRZ6COSpgzrC7Vhx6R7nRf4NaCYjKnYt/wtCp5onM6n6I4q5OtPsi3HEL
2sORt8JhC1M2/k/hlV+U/psCgYBkf40QuOs2aXjoj9jJR46HaJdeKDvkGG1v0+Ee
fLHehixuK/chIlhETZ3b0BqgenQiJIUcrc/uMvwqoowlstd9JjwVzkti3XGkqbsl
jSVFnbWnln12UcJIlQ8nLYc8NK0ZWM2M3OtmaeKyNGHCZyLROLRF/qRzWEOaHhS1
scbuIwKBgQC4idVhpLxyGPhMX1Z6qrrl4Vhul6Dev+TUZ9DZc3LRlMaBtu4KrBI3
s2TCwENcCgp0MZKSIROrPEywsrg0HCyhIEjL6No2rPlqgzP9Lx2npit2RxGgvbO/
Ad5wsCe5c4Kfs3wrbxEXcIVNE39ihjAfhXPlz2uUFAAOwHP8H/hG7w==
-----END RSA PRIVATE KEY-----
`)

type Claims struct {
	Email string `json:"email"`
	Group string `json:"group"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, id string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: email,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

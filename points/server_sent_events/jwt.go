package main

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type claims struct {
	ClientID string `json:"clientID"`
	jwt.StandardClaims
}

func parseToken(jwtToken string) (*claims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if token == nil {
		return nil, err
	}

	claimsIns, ok := token.Claims.(*claims)
	if ok && token.Valid {
		return claimsIns, nil
	}

	return claimsIns, err
}

func generateToken(userID string) (string, error) {
	claimsIns := claims{
		ClientID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(cookieValidPeriod)).Unix(),
		},
	}

	tokenIns := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsIns)
	return tokenIns.SignedString(jwtSecret)
}

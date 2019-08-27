package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func JWTGenToken(key string,m map[string]interface{}) string{
	//init a Token struct
	token := jwt.New(jwt.SigningMethodHS256)

	//init a MapClaims struct
	claims := make(jwt.MapClaims)

	for index,value :=range m {
		claims[index] = value
	}
	token.Claims = claims
	tokenString,_ := token.SignedString([]byte(key))
	return tokenString
}


func JWTParseToken(tokenString string,key string) (interface{}, bool){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return "", false
	}
}

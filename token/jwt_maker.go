package token

// import (
// 	"github.com/golang-jwt/jwt/v5"
// 	"time"
// 	"fmt"
// )

// type JWTMaker struct{
// 	secretKey string
// }
// const(
// 	minSecretKeySize=6
// )
// func NewJWTMaker(secretKey string) (Maker, error){
// 	if len(secretKey) < minSecretKeySize{
// 		return nil, fmt.Errorf("invalid key size: must be at least %d characters",minSecretKeySize)
// 	}
// 	return &JWTMaker{secretKey}, nil
// }	

// func (maker *JWTMaker) CreateToken(username string,duration time.Duration) (string,error){
// 	payload,err := NewPayload(username,duration)
// 	if err != nil{
// 		return "",nil
// 	}
// 	mapClaims:=jwt.MapClaims{
// 		"payload":payload,
// 	}
// 	jwtToken:=jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
// 	return jwtToken.SignedString([]byte(maker.secretKey))
// }
	
// func (maker *JWTMaker) VerifyToken(token string) (*Payload,error){
// 	keyfunc := func(token *jwt.Token) (interface{},error){
// 		_,ok :=token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok{
// 			return nil, fmt.Errorf("unexpected signing method: %v",token.Header["alg"])
// 		}
// 		return []byte(maker.secretKey),nil
// 	}
// 	jwtToken,err := jwt.ParseWithClaims(token,&jwt.MapClaims{},keyfunc)
// 	if err !=nil{
// 		return nil,fmt.Errorf("invalid token: %w",err)
// 	}
// 	payload,ok :=jwtToken.Claims.(*jwt.MapClaims)
// 	if !ok || jwtToken.Valid{
// 		return nil,fmt.Errorf("invalid token")
// 	}
// 	return payload["payload"],nil

// }
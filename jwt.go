package zlib

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

func GetJwtToken(secretKey string,appId int ,uid int ,username string)string{
	header := JwtDataHeader{
		Alg: "HS256",
		Typ:"JWT",
	}
	current := time.Now().Second()
	ExpireTime := int(current) + (   2 * 60 * 60 )
	payload := JwtDataPayload{
		Id:uid,
		Expire: ExpireTime,
		ATime: current,
		AppId: appId,
		Username: username,
	}

	headerJson,_ := json.Marshal(header)
	payloadJson ,_ := json.Marshal(payload)

	//fmt.Println("json : ",string(headerJson),string(payloadJson))

	base64HeaderJson := EncodeSegment(headerJson)
	base64PayloadJson := EncodeSegment(payloadJson)

	base64HeaderPayload := base64HeaderJson + "." + base64PayloadJson

	//fmt.Println("base64HeaderPayload : ",base64HeaderPayload)
	hasher := hmac.New(crypto.SHA256.New , []byte(secretKey))
	hasher.Write([]byte(base64HeaderPayload))

	sign := hasher.Sum(nil)

	base64Sign :=  EncodeSegment(sign)
	//fmt.Println(  " base64Sign : " , base64Sign)
	jwtString := base64HeaderPayload + "." + base64Sign
	fmt.Println("myself : ",jwtString)


	type jwtCustomClaims struct {
		jwt.StandardClaims
		Id 			int
		Expire 		int
		ATime		int
		AppId		int
		Username	string
		// 追加自己需要的信息
		//Uid   uint `json:"uid"`
		//Admin bool `json:"admin"`
	}

	claims := &jwtCustomClaims{
		//StandardClaims: jwt.StandardClaims{
		//	ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
		//},
		Id:uid,
		Expire: ExpireTime,
		ATime: current,
		AppId: appId,
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	fmt.Println("jwt-go : ",tokenString)


	myToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiRXhwaXJlIjo3MjA5LCJBVGltZSI6OSwiQXBwSWQiOjEsIlVzZXJuYW1lIjoid2FuZyJ9.U57wFFeADbDdRj0MuEF0mNfSZ_JgD3wUEzGhRE02jOI"
	rs,err := ParseToken(myToken,[]byte(secretKey))
	fmt.Println("ParseToken : ",rs ,err)

	myParseToken(myToken,secretKey)

	return jwtString
}
func myParseToken(tokenStr string, SecretKey string){
	if tokenStr == ""{

	}
	tokenStr = strings.Trim(tokenStr," ")
	if tokenStr == ""{

	}

	tokenArr := strings.Split(tokenStr,".")
	if len(tokenArr) != 3{

	}
	fmt.Println(tokenArr)
	headerBase64 := tokenArr[0]
	headerJsonStr,err := DecodeSegment(headerBase64)
	if err != nil{

	}

	if  len(headerJsonStr) == 0{

	}

	fmt.Println(string(headerJsonStr))
	headerStruct := JwtDataHeader{}
	//err := json.Unmarshal(headerJsonStr,&headerStruct)
	//if err != nil{
	//	fmt.Println("err")
	//	os.Exit(-100)
	//}
	//payloadBase64 := tokenArr[1]

	fmt.Println(headerStruct)
	os.Exit(-100)
}

func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}


func ParseToken(tokenSrt string, SecretKey []byte) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims = token.Claims
	return
}

//type JwtData struct {
//	header 	JwtDataHeader
//	payload	JwtDataPayload
//	sign 	string
//}

type JwtDataHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtDataPayload struct {
	Id 			int
	Expire 		int
	ATime		int
	AppId		int
	Username	string
}

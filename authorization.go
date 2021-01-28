package xlib

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bearer string = "bearer"
)

//TokenDetails structure describing all properties need for authorization
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

//AccessDetails contains access token unique indentifier and user id
type AccessDetails struct {
	TokenUUID string
	UserID    primitive.ObjectID
}

type tokenType int

const (
	//AccessTokenType type of access token
	AccessTokenType tokenType = iota
	//RefreshTokenType type of refresh token
	RefreshTokenType
)

//CreateToken create access and refresh token jwt
func CreateToken(userid primitive.ObjectID) (*TokenDetails, error) {
	td := &TokenDetails{}
	//Access token expires after 15 minutes
	iat := time.Now()
	iatu := iat.Unix()
	td.AtExpires = iat.Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	//Refresh token expires after 1 week
	td.RtExpires = iat.Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["iat"] = iatu
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	atClaims["iat"] = iatu
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

//ExtractToken get token from header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//VerifyToken check the signing method
func VerifyToken(tokenString string, t tokenType) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if t == AccessTokenType {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil

	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//TokenValid check if token has expired
func TokenValid(tokenString string, t tokenType) (*AccessDetails, error) {
	token, err := VerifyToken(tokenString, t)
	if err != nil {
		return nil, err
	}
	aUUID := "access_uuid"
	tstr := "Access token"
	if t == RefreshTokenType {
		aUUID = "refresh_uuid"
		tstr = "Refresh token"
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenUUID, ok := claims[aUUID].(string)
		if !ok {
			return nil, errors.New(aUUID + " error")
		}
		userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			TokenUUID: tokenUUID,
			UserID:    userID,
		}, nil
	}
	return nil, errors.New(tstr + " expired")
}

//ExtractTokenMetadata get the metadata stored in the token
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	return TokenValid(ExtractToken(r), AccessTokenType)
}

//CreateAuth set token to mem
func CreateAuth(m IMemCash, userid primitive.ObjectID, td *TokenDetails) error {
	//converting Unix to UTC(to Time object)
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	errAccess := m.SetKey(td.AccessUUID, userid.Hex(), at.Sub(now))
	if errAccess != nil {
		return errAccess
	}
	errRefresh := m.SetKey(td.RefreshUUID, userid.Hex(), rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

//FetchAuth fetch token from mem
func FetchAuth(m IMemCash, authD *AccessDetails) (primitive.ObjectID, error) {
	userid, err := m.GetKey(authD.TokenUUID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	userID, _ := primitive.ObjectIDFromHex(userid)
	return userID, nil
}

//DeleteAuth delete token in redis
func DeleteAuth(m IMemCash, givenUUID string) (int64, error) {
	deleted, err := m.DelKey(givenUUID)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

//HashAndSalt hash password
func HashAndSalt(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

//ComparePasswords compare hashed Password with giaven plain password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		return false
	}

	return true
}

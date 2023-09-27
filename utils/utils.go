package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"math/rand"
	"path/filepath"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

var letters = []rune("0123456789")

const length = 6

func RandSeq() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MonthFromNow(how uint) time.Time {
	du := time.Hour * time.Duration(24*30*how)
	//TODO make duration from tonight 12:00
	return time.Now().Add(du)
}

func ToHstore(in map[string]string) map[string]*string {
	values := make(map[string]*string)
	for k, v := range in {
		values[k] = &v
	}
	return values
}

func HstoreToMap(in map[string]*string) map[string]string {
	values := make(map[string]string)
	for k, v := range in {
		values[k] = *v
	}
	return values
}

//KeyPair preparing key/cert files
func GRPCKeyPair() (string, string) {
	keyFile := config.Config().GetString(KeyGRPCKeyFile)
	certFile := config.Config().GetString(KeyGRPCCertFile)
	return filepath.Join("certs", certFile), filepath.Join("certs", keyFile)
}

//Get secretkey
var secretKey string

func init() {
	secretKey = config.Config().GetString(KeySecretKey)
}

// valid validates the authorization.
func AuthValid(ctx context.Context, md metadata.MD) (context.Context, bool) {
	authorization := md[KeyAuthToken]
	log.Logger().Info(md)
	if len(authorization) < 1 {
		return ctx, false
	}
	tokenString := authorization[0]
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	log.Logger().Debugf("Token: %v, err: %v", token, err)
	if err == nil && token.Valid {
		// add user ID to the context
		usrID := uint(token.Claims.(jwt.MapClaims)["sub"].(float64))
		newCtx := context.WithValue(ctx, KeyEmail, usrID)
		return newCtx, true
	}
	return ctx, false
}

func TimestampMD5() string {
	return MD5(time.Now().Format(time.StampNano))
}

// MD5 hashes using md5 algorithm
func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ConvertUint(in []uint32) []int64 {
	out := []int64{}
	for _, u := range in {
		out = append(out, int64(u))
	}
	return out
}

func ConvertInt64(in []int64) []uint32 {
	out := []uint32{}
	for _, u := range in {
		out = append(out, uint32(u))
	}
	return out
}

var nameRegex = regexp.MustCompile("^[a-z]+[a-z0-9-_]*[a-z0-9]+$")

func NameValid(name string) error {
	if len(name) < 3 || !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ToTimestamp(t *time.Time) *pb.Timestamp {
	return &pb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   0,
	}
}

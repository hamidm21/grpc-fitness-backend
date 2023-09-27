package controller

import (
	"context"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"
	// "context"
	// "regexp"
	// "time"
	// "github.com/jinzhu/gorm"
	// "gitlab.com/mefit/mefit-server/entity"
	// "gitlab.com/mefit/mefit-server/utils/log"
	// pb_admin "gitlab.com/mefit/mefit-server/proto"
	// jwt "github.com/dgrijalva/jwt-go"
	// pb "gitlab.com/mefit/mefit-api/proto"
	// "gitlab.com/mefit/mefit-server/utils"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (s *Controller) AnonySignUp(ctx context.Context, in *pb.AnonyReq) (*pb.AnonyRes, error) {
	usr := entity.User{}
	//try renew anonymous token
	if len(in.AnonyId) > 0 {
		usr.AnonymousID = in.AnonyId
		if err := entity.SimpleCrud(usr).Get(usr); err == nil {
			token, _ := createJWT(&usr)
			log.Logger().Print("JWT is ...", token)
			return &pb.AnonyRes{AnonyId: in.AnonyId, Token: token}, nil
		}
	}
	//Create a new anonmous user
	anonyID, err := createAnonyID()
	if err != nil {
		return nil, err
	}
	usr.Anonymous = true
	usr.AnonymousID = anonyID
	//FIXME: quick hack
	usr.Email = anonyID
	//Fetch the very first plan
	plan1 := entity.Plan{}

	if err := entity.SimpleCrud(entity.Plan{MetaName: "first"}).Get(&plan1); err != nil {
		return nil, utils.ErrNotFound
	}
	pro := &entity.Profile{}
	pro.User = usr
	// pro.Name = anonyID
	pro.PlanID = plan1.ID
	if err := entity.SimpleCrud(pro).Save(); err != nil {
		return nil, err
	}
	if err := entity.SimpleCrud(usr).Get(&usr); err != nil {
		return nil, utils.ErrInternal
	}
	token, _ := createJWT(&usr)
	log.Logger().Print("JWT is ...", token)
	return &pb.AnonyRes{AnonyId: anonyID, Token: token}, nil
}

// Register register a new user
func (s *Controller) SignUp(ctx context.Context, in *pb.SignUpReq) (*pb.Empty, error) {
	//validate mail
	if !emailRegex.MatchString(in.Email) {
		return nil, utils.ErrInvalidEmail
	}
	q := &entity.User{
		Email: in.Email,
	}

	if err := entity.SimpleCrud(q).Get(q); err == nil {
		return nil, utils.ErrExists
	}
	///////////////////////////
	anonyId := in.AnonyId
	anonUser := &entity.User{}
	if err := entity.SimpleCrud(entity.User{AnonymousID: anonyId}).Get(anonUser); err != nil {
		log.Logger().Errorf("Signup: Check anoy user exists: %v", err)
		return nil, utils.ErrNotFound
	}

	//generate unique hash
	token := utils.GetMD5Hash(in.Email) + utils.RandSeq()

	anonUser.Email = in.Email
	anonUser.Password = in.Password
	anonUser.ConfirmToken = token
	anonUser.AnonymousID = anonyId
	anonUser.Anonymous = false
	anonUser.Confirmed = true

	if err := entity.SimpleCrud(anonUser).Save(); err != nil {
		log.Logger().Errorf("Signup: %v", err)
		return nil, utils.ErrInternal
	}
	return &pb.Empty{}, nil
}

//Login user controller
func (s *Controller) SignIn(ctx context.Context, in *pb.AuthReq) (*pb.SignInRes, error) {
	if !emailRegex.MatchString(in.Email) {
		return nil, utils.ErrInvalidEmail
	}
	usr := &entity.User{Email: in.Email}
	//Fetch user for this email address
	if err := entity.SimpleCrud(usr).Get(usr); err != nil {
		return nil, utils.ErrInvalidCreds
	}
	hashedPassword, _ := utils.HashPassword(in.Password)
	if !utils.CheckPasswordHash(usr.Password, hashedPassword) {
		return nil, utils.ErrInvalidCreds
	}
	token, _ := createJWT(usr)
	return &pb.SignInRes{Token: token}, nil
}

func (s *Controller) Logout(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	//TODO: logout logic here
	return &pb.Empty{}, nil
}

func createJWT(usr *entity.User) (string, error) {
	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := make(jwt.MapClaims)
	claims["sub"] = usr.ID
	// Expire in 5 mins
	claims["exp"] = time.Now().Add(time.Hour * 24 * 300).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", utils.ErrUnavailable
	}
	return tokenString, nil
}

func createAnonyID() (string, error) {
	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := make(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	id, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", utils.ErrUnavailable
	}
	return id, nil
}

package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata        = status.Errorf(codes.InvalidArgument, "missing metadata")
	ErrUnauthenticated        = status.Errorf(codes.Unauthenticated, "registration required")
	ErrPermissionDenied       = status.Errorf(codes.PermissionDenied, "vip registration required")
	ErrInvalidName            = status.Errorf(codes.InvalidArgument, "invalid name")
	ErrInvalidEmail           = status.Errorf(codes.InvalidArgument, "invalid email address")
	ErrInvalidCreds           = status.Errorf(codes.Unauthenticated, "invalid account/password")
	ErrConfirm                = status.Errorf(codes.InvalidArgument, "password does not match the confirm password")
	ErrInvalidAttachment      = status.Errorf(codes.InvalidArgument, "invalid attachment name")
	ErrInternal               = status.Errorf(codes.Internal, "internal error")
	ErrExists                 = status.Errorf(codes.AlreadyExists, "Record already exist or configured!")
	ErrAppNotInstalled        = status.Errorf(codes.Internal, "app can not be installed")
	ErrUserHaveNoPlan         = status.Errorf(codes.InvalidArgument, "user have no plan")
	ErrInvalidVars            = status.Errorf(codes.InvalidArgument, "invalid variables")
	ErrInvalidUrl             = status.Errorf(codes.InvalidArgument, "invalid url")
	ErrInvalidMountPath       = status.Errorf(codes.InvalidArgument, "invalid mount path")
	ErrInvalidDN              = status.Errorf(codes.InvalidArgument, "invalid domain name")
	ErrWrongPlanVariable      = status.Errorf(codes.InvalidArgument, "missing/invalid input plan")
	ErrInvalidToken           = status.Errorf(codes.Unauthenticated, "invalid token")
	ErrInvalidActivationToken = status.Errorf(codes.Unauthenticated, "invalid activation code")
	ErrUnavailable            = status.Errorf(codes.Unavailable, "unavailable error")
	ErrUnknown                = status.Errorf(codes.Unknown, "unknown error")
	ErrInvalidRating          = status.Errorf(codes.Unknown, "invalid rating range")
	ErrBadScale               = status.Errorf(codes.InvalidArgument, "scale.min must be lower than maximum scaling")
	ErrBadPort                = status.Errorf(codes.InvalidArgument, "port must be lower than 65535")
	ErrNotFound               = status.Errorf(codes.NotFound, "not found")
	ErrNotReady               = status.Errorf(codes.Unavailable, "resource not ready")
)

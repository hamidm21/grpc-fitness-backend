package controller

import (
	"context"

	"gitlab.com/mefit/mefit-server/utils/log"

	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"

	pb "gitlab.com/mefit/mefit-api/proto"
)

//ProfileUpdate controller
func (s *Controller) ProfileUpdate(ctx context.Context, in *pb.Profile) (*pb.Profile, error) {
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	pro := &entity.Profile{}
	log.Logger().Infof("ProfileUpdate: userID: %v", usrID)
	if err := entity.Crud(pro, usrID).Get(pro); err != nil {
		return nil, err
	}
	profileEntityFrom(pro, in)
	if err := entity.Crud(pro, usrID).Save(); err != nil {
		return nil, err
	}
	return profileInfoFrom(usrID, pro)
}

//ProfileInfo user controller
func (s *Controller) ProfileInfo(ctx context.Context, in *pb.Empty) (*pb.Profile, error) {
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	return profileInfoFrom(usrID, &entity.Profile{})
}

func profileEntityFrom(out *entity.Profile, in *pb.Profile) {
	//Just boilerplate

	if len(in.Name) > 0 {
		out.Name = in.Name
	}
	log.Logger().Debugf("####### Gender in: %v", in.Gender)
	if in.Gender > 0 && in.Gender < 4 {
		out.Gender = entity.Gender(in.Gender)
	}
	log.Logger().Debugf("####### Gender in: %v", out.Gender)

	if in.Age > 0 {
		out.Age = in.Age
	}
	if in.Height > 0 {
		out.Height = in.Height
	}
	if in.Waist > 0 {
		out.Waist = in.Waist
	}
	if in.Neck > 0 {
		out.Neck = in.Neck
	}
	if in.Hip > 0 {
		out.Hip = in.Hip
	}
	if in.Arm > 0 {
		out.Arm = in.Arm
	}
	if in.Leg > 0 {
		out.Leg = in.Leg
	}
	if in.CurrentWeight > 0 {
		out.CurrentWeight = in.CurrentWeight
	}
	if in.TargetWeight > 0 {
		out.TargetWeight = in.TargetWeight
	}
	//Ignore invalid values
	if in.ActivityLevel > 0 && in.ActivityLevel < 6 {
		out.ActivityLevel = entity.ActivityLevel(in.ActivityLevel)
	}
	//Ignore invalid values
	if in.Goal > 0 && in.Goal < 4 {
		out.Goal = entity.Goal(in.Goal)
	}
	if in.DaysOfWeek != nil && len(in.DaysOfWeek) > 0 && len(in.DaysOfWeek) < 8 {
		for _, day := range in.DaysOfWeek {
			if day == 0 || day > 8 {
				return
			}
		}

		// out.DaysOfWeek = pq.Int64Array(utils.ConvertUint(in.DaysOfWeek))
	}
}
func profileInfoFrom(usrID uint, in *entity.Profile) (*pb.Profile, error) {
	if in.ID == 0 {
		if err := entity.Crud(*in, usrID).Get(in); err != nil {
			return nil, err
		}
	}

	return &pb.Profile{
		Name:          in.Name,
		Gender:        uint32(in.Gender),
		Age:           in.Age,
		Height:        in.Height,
		Waist:         in.Waist,
		Neck:          in.Neck,
		Hip:           in.Hip,
		Arm:           in.Arm,
		Leg:           in.Leg,
		CurrentWeight: in.CurrentWeight,
		TargetWeight:  in.TargetWeight,
		ActivityLevel: uint32(in.ActivityLevel),
		Goal:          uint32(in.Goal),
		// DaysOfWeek:    utils.ConvertInt64(in.DaysOfWeek),
	}, nil
}

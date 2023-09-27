package controller

// import (
// 	"context"

// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/entity"
// 	"gitlab.com/mefit/mefit-server/utils"

// 	pb "gitlab.com/mefit/mefit-api/proto"
// )

// func (s *Controller) GetProgram(ctx context.Context, in *pb.Empty) (*pb.ProgramRes, error) {
// 	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
// 	usr := entity.User{}
// 	usr.ID = usrID
// 	if err := entity.SimpleCrud(usr).Get(&usr); err != nil {
// 		return nil, err
// 	}
// 	levels := []entity.Level{}
// 	if err := entity.GetDB().Preload("SubLevels", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("mefit_sub_level.no asc")
// 	}).Preload("SubLevels.Plans", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("mefit_plan.no asc")
// 	}).Order("no asc").Find(&levels).Error; err != nil {
// 		return nil, utils.ErrNotFound
// 	}
// 	return listLevelFrom(&usr, levels)
// }

// func listLevelFrom(usr *entity.User, items []entity.Level) (*pb.ProgramRes, error) {
// 	program := &pb.ProgramRes{}
// 	program.Levels = []*pb.Level{}
// 	// program.Levels.SubLevels = []*pb.SubLevel{}
// 	for _, level := range items {
// 		lv := &pb.Level{
// 			Name:        level.Name,
// 			Description: level.Description,
// 			LevelId:     uint32(level.ID),
// 			Value:       uint32(level.No),
// 		}
// 		lv.SubLevels = []*pb.SubLevel{}
// 		for _, subLevels := range level.SubLevels {
// 			sl := &pb.SubLevel{
// 				Name:  subLevels.Name,
// 				Value: uint32(subLevels.No),
// 			}
// 			sl.Plans = []*pb.SubLevelPlan{}
// 			for _, subLevelPlans := range subLevels.Plans {
// 				sp := &pb.SubLevelPlan{
// 					ThumbnailUrl: subLevelPlans.ThumbnailUrl,
// 					Name:         subLevelPlans.Name,
// 					Id:           uint32(subLevelPlans.ID),
// 					VipLock:      subLevelPlans.VipLock && !usr.VIP,
// 				}
// 				sl.Plans = append(sl.Plans, sp)
// 			}
// 			lv.SubLevels = append(lv.SubLevels, sl)
// 		}
// 		program.Levels = append(program.Levels, lv)
// 	}
// 	return program, nil
// }

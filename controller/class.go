package controller

import (
	"context"

	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"

	pb "gitlab.com/mefit/mefit-api/proto"
)

func (s *Controller) GetClassMovements(ctx context.Context, in *pb.FindByIdAndPage) (*pb.ListMovementsRes, error) {
	movements := []entity.Movement{}
	q := entity.Movement{}
	q.ClassID = uint(in.Id)
	if err := entity.SimpleCrud(q).LimitedList(&movements, in.Page, ""); err != nil {
		return nil, utils.ErrNotFound
	}
	return ListMovementsFrom(movements, uint(in.Id))
}

func (s *Controller) GetClasses(ctx context.Context, in *pb.Empty) (*pb.ListClassRes, error) {
	classes := []entity.Class{}
	if err := entity.SimpleCrud(entity.Class{}).List(&classes, nil); err != nil {
		return nil, utils.ErrNotFound
	}
	return listClassesFrom(classes)
}

func (s *Controller) ClassMovementInfo(ctx context.Context, in *pb.FindByIdReq) (*pb.Movement, error) {
	Movement := entity.Movement{}
	if err := entity.SimpleCrud(entity.Movement{}).ID(in.Id).Get(&Movement); err != nil {
		return nil, utils.ErrNotFound
	}
	if err := entity.SimpleCrud(entity.Movement{}).ID(in.Id).Related(&Movement.MuscleGroups, "MuscleGroups"); err != nil {
		return nil, utils.ErrNotFound
	}
	return movementFrom(Movement)
}

func ListMovementsFrom(items []entity.Movement, id uint) (*pb.ListMovementsRes, error) {
	var count uint
	if err := entity.GetDB().Model(entity.Movement{ClassID: id}).Count(&count).Error; err != nil {
		return nil, utils.ErrInternal
	}
	listMovements := &pb.ListMovementsRes{}
	listMovements.Movements = []*pb.Movement{}
	for _, move := range items {
		MoveMent := &pb.Movement{
			Name:         move.Name,
			ThumbnailUrl: move.GetThumbnailUrl(),
			Description:  move.Description,
			Instruction:  move.Instruction,
			VideoUrl:     move.GetVideoUrl(),
			Tips:         move.Tips,
			Id:           uint32(move.ID),
			MuscleGroups: []*pb.MuscleGroup{},
			NameFa:       move.NameFa,
		}
		for _, MuscleGroup := range move.MuscleGroups {
			mg := &pb.MuscleGroup{
				Name: MuscleGroup.Name,
			}
			MoveMent.MuscleGroups = append(MoveMent.MuscleGroups, mg)
		}
		// for _, Article := range move.Articles {
		// 	at := &pb.Article {
		// 		ThumbnailUrl: Article.ThumbnailUrl,
		// 		CoverUrl:     Article.CoverUrl,
		// 		Title:        Article.Title,
		// 		Body:         Article.Body,
		// 		Id:           uint32(Article.ID),
		// 	}
		// 	MoveMent.Articles = append(MoveMent.Articles, at)
		// }
		listMovements.Movements = append(listMovements.Movements, MoveMent)
	}
	listMovements.TotalCount = int32(count)
	return listMovements, nil
}

func listClassesFrom(items []entity.Class) (*pb.ListClassRes, error) {
	classes := &pb.ListClassRes{}
	classes.Classes = []*pb.Class{}
	for _, class := range items {
		cls := &pb.Class{
			Name:        class.Name,
			Id:          uint32(class.ID),
			CoverUrl:    class.GetCoverUrl(),
			Description: class.Description,
		}
		for _, MuscleGroup := range class.MuscleGroups {
			mg := &pb.MuscleGroup{
				Name: MuscleGroup.Name,
			}
			cls.MuscleGroups = append(cls.MuscleGroups, mg)
		}
		classes.Classes = append(classes.Classes, cls)
	}
	return classes, nil
}

func movementFrom(item entity.Movement) (*pb.Movement, error) {
	move := &pb.Movement{
		Name:   item.Name,
		NameFa: item.NameFa,

		ThumbnailUrl: item.GetThumbnailUrl(),
		Description:  item.Description,
		Instruction:  item.Instruction,
		Tips:         item.Tips,
		VideoUrl:     item.GetVideoUrl(),
		Id:           uint32(item.ID),
		MuscleGroups: []*pb.MuscleGroup{},
	}
	for _, MuscleGroup := range item.MuscleGroups {
		mg := &pb.MuscleGroup{
			Name: MuscleGroup.Name,
		}
		move.MuscleGroups = append(move.MuscleGroups, mg)
	}
	return move, nil
}

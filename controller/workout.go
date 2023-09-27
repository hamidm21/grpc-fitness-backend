package controller

import (
	"context"

	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"

	pb "gitlab.com/mefit/mefit-api/proto"
)

func (s *Controller) GetWorkout(ctx context.Context, in *pb.FindByIdReq) (*pb.Workout, error) {
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	pro := &entity.Profile{UserID: usrID}
	if err := entity.SimpleCrud(pro).Get(pro, "Plan"); err != nil {
		return nil, err
	}
	workout := entity.Workout{}
	workout.ID = uint(in.Id)
	if err := entity.SimpleCrud(workout).Get(&workout, "WorkoutType", "ExerciseSections", "ExerciseSections.Exercises", "ExerciseSections.Exercises.Movement"); err != nil {
		return nil, utils.ErrNotFound
	}
	return workoutFrom(pro, &workout)
}

func (s *Controller) GetPromotedWorkouts(ctx context.Context, in *pb.Empty) (*pb.WorkoutList, error) {
	wl := []*pb.Workout{}
	return &pb.WorkoutList{
		Workouts: wl,
	}, nil
}

func (s *Controller) FinishWorkout(ctx context.Context, in *pb.FeedbackReq) (*pb.Empty, error) {
	log.Logger().Print("Finish workout started")
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	pro := &entity.Profile{UserID: usrID}
	if err := entity.SimpleCrud(pro).Get(pro, "Plan"); err != nil {
		return nil, err
	}
	workout := &entity.Workout{}
	workout.ID = uint(in.Id)
	//Double check if workout exists!
	if err := entity.SimpleCrud(workout).Get(&workout); err != nil {
		return nil, utils.ErrNotFound
	}
	WorkoutHistory := &entity.WorkoutHistory{
		ProfileID:  pro.ID,
		WorkoutID:  workout.ID,
		Rating:     uint(in.Rate),
		Difficulty: entity.Difficulty(in.Difficulty),
	}
	if !WorkoutHistory.Valid() {
		return nil, utils.ErrInvalidRating
	}

	if err := entity.SimpleCrud(WorkoutHistory).Create(); err != nil {
		return nil, utils.ErrInternal
	}
	return &pb.Empty{}, nil
}

func workoutFrom(pro *entity.Profile, item *entity.Workout) (*pb.Workout, error) {
	lock := false
	if item.VipLock == true {
		product := entity.Product{PlanID: pro.PlanID}
		if err := entity.SimpleCrud(product).Get(&product); err != nil {
			lock = true
		}
		purchased := entity.PurchasedProduct{}
		purchased.ProfileID = pro.ID
		purchased.ProductID = product.ID
		if err := entity.SimpleCrud(purchased).Get(&purchased); err != nil {
			lock = true
		}
	}

	//First check if workout is done by the user
	workout := &pb.Workout{
		Name:        item.WorkoutType.Name,
		Description: item.WorkoutType.Description,
		Instruction: item.Instruction,
		Calorie:     uint32(item.Calorie),
		Id:          uint32(item.ID),
		Duration:    uint32(item.Duration),
		Viplock:     lock,
	}

	workout.ExerciseSections = []*pb.ExerciseSection{}
	for _, sections := range item.ExerciseSections {
		es := &pb.ExerciseSection{
			Round: uint32(sections.Round),
		}

		//Fix bug related to gorm preload eager loading
		movements := map[uint]*entity.Movement{}
		for _, ex := range sections.Exercises {
			if ex.Movement != nil {
				movements[ex.MovementID] = ex.Movement
			}
		}
		es.Exercises = []*pb.ExercisePeriod{}
		for _, exercise := range sections.Exercises {
			ep := &pb.ExercisePeriod{
				Name:         exercise.Movement.NameFa,
				ThumbnailUrl: exercise.Movement.GetThumbnailUrl(),
				MovementId:   uint32(exercise.MovementID),
				ExerciseType: uint32(exercise.ExerciseType),
				Value:        uint32(exercise.Value),
			}
			//set rest
			ep.Rest = 0
			if exercise.Rest != nil {
				ep.Rest = *exercise.Rest
			}
			log.Logger().Debugf("ep: %+v", ep)
			exerciseMovement := movements[exercise.MovementID]
			mv := &pb.Movement{
				Name:         exerciseMovement.NameFa,
				ThumbnailUrl: exerciseMovement.GetThumbnailUrl(),
				Description:  exerciseMovement.Description,
				Instruction:  exerciseMovement.Instruction,
				VideoUrl:     exerciseMovement.GetVideoUrl(),
				Tips:         exerciseMovement.Tips,
				Id:           uint32(exerciseMovement.ID),
			}

			ep.Movement = mv
			es.Exercises = append(es.Exercises, ep)
		}
		workout.ExerciseSections = append(workout.ExerciseSections, es)
	}
	log.Logger().Debugf("%v", *workout)
	return workout, nil
}

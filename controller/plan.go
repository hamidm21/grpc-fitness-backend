package controller

import (
	"context"

	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"

	// "gitlab.com/mefit/mefit-server/utils/log"

	pb "gitlab.com/mefit/mefit-api/proto"
)

func (s *Controller) GetPlan(ctx context.Context, in *pb.FindByIdReq) (*pb.PlanRes, error) {
	plan := entity.Plan{}
	plan.ID = uint(in.Id)
	if err := entity.SimpleCrud(plan).Get(&plan); err != nil {
		return nil, utils.ErrNotFound
	}
	pro := &entity.Profile{UserID: ctx.Value(utils.KeyEmail).(uint)}
	if err := entity.SimpleCrud(pro).Get(pro, "Plan"); err != nil {
		return nil, err
	}
	return planFrom(&plan, pro)
}

func (s *Controller) GetPlans(ctx context.Context, in *pb.Empty) (*pb.PlansRes, error) {
	plans := []entity.Plan{}
	p := entity.Plan{}
	if err := entity.SimpleCrud(p).List(&plans, nil); err != nil {
		return nil, utils.ErrNotFound
	}
	pro := &entity.Profile{UserID: ctx.Value(utils.KeyEmail).(uint)}
	if err := entity.SimpleCrud(pro).Get(pro, "Plan"); err != nil {
		return nil, err
	}
	return listPlansFrom(plans, pro)
}

func (s *Controller) GetCurrentPlan(ctx context.Context, in *pb.Empty) (*pb.CurrentPlanRes, error) {
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	log.Logger().Infof("userID::: %v", usrID)
	pro := &entity.Profile{UserID: usrID}
	if err := entity.SimpleCrud(pro).Get(pro, "Plan"); err != nil {
		return nil, err
	}
	return currentPlanFrom(pro, pro.Plan)
}

func (s *Controller) JoinPlan(ctx context.Context, in *pb.JoinPlanReq) (*pb.JoinPlanRes, error) {
	usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	pro := &entity.Profile{UserID: usrID}
	if err := entity.SimpleCrud(pro).Get(pro, "User"); err != nil {
		return nil, err
	}
	plan := entity.Plan{}
	plan.ID = uint(in.Id)
	//Double check if plan exists!
	if err := entity.SimpleCrud(plan).Get(&plan); err != nil {
		return nil, utils.ErrNotFound
	}
	//Initialize response with worst case
	res := &pb.JoinPlanRes{
		State: pb.State_DONE,
	}
	//TODO: check for vip plan here!
	if plan.VipLock == true && plan.HasTrial == false {
		product := entity.Product{PlanID: plan.ID}
		if err := entity.SimpleCrud(product).Get(&product); err != nil {
			log.Logger().Error(err)
			res.State = pb.State_NEED_PURCHASE
		}
		purchased := entity.PurchasedProduct{}
		purchased.ProfileID = pro.ID
		purchased.ProductID = product.ID
		if err := entity.SimpleCrud(purchased).Get(&purchased); err != nil {
			log.Logger().Error(err)
			res.State = pb.State_NEED_PURCHASE
		}
	}
	if err := entity.SimpleCrud(pro).Updates(entity.Profile{PlanID: plan.ID, CurrentWorkoutNo: 1}); err != nil {
		return nil, utils.ErrInternal
	}
	return res, nil
}

func planFrom(item *entity.Plan, pro *entity.Profile) (*pb.PlanRes, error) {
	lock := false
	if item.VipLock == true {
		product := entity.Product{PlanID: item.ID}
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
	// get total calories for workouts inside plan
	var totalCalories uint32
	workouts := []entity.Workout{}
	if err := entity.GetDB().Model(&entity.Workout{}).Preload("WorkoutType").Where(([]int64)(item.WorkoutArray)).Find(&workouts).Error; err != nil {
		log.Logger().Errorf("cant get workout for plan: %v", item.ID)
	}
	for i := 0; i < len(workouts); i++ {
		totalCalories += uint32(workouts[i].Calorie)
	}
	plan := &pb.PlanRes{
		ThumbnailUrl:  item.ThumbnailUrl,
		TotalCalorie:  totalCalories,
		Name:          item.Name,
		Description:   item.Description,
		Id:            uint32(item.ID),
		WorkoutCounts: uint32(len(item.WorkoutArray)),
		Weeks:         uint32(item.Weeks),
		Level:         uint32(item.Level),
		HasTrial:      item.HasTrial,
		VipLock:       lock,
		CoverUrl:      item.CoverUrl,
	}
	return plan, nil
}

func listPlansFrom(items []entity.Plan, pro *entity.Profile) (*pb.PlansRes, error) {

	ListPlans := &pb.PlansRes{}
	ListPlans.Plans = []*pb.PlanRes{}
	for _, plan := range items {
		lock := false
		if plan.VipLock == true {
			product := entity.Product{PlanID: plan.ID}
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
		log.Logger().Print("lock is.... ", lock)
		pl := &pb.PlanRes{
			ThumbnailUrl:  plan.ThumbnailUrl,
			Name:          plan.Name,
			Description:   plan.Description,
			Id:            uint32(plan.ID),
			WorkoutCounts: uint32(len(plan.WorkoutArray)),
			Weeks:         uint32(plan.Weeks),
			Level:         uint32(plan.Level),
			VipLock:       lock,
			HasTrial:      plan.HasTrial,
			CoverUrl:      plan.CoverUrl,
		}
		ListPlans.Plans = append(ListPlans.Plans, pl)
	}
	return ListPlans, nil
}

// func workoutCount(item *entity.Plan) uint32 {
// 	var count uint32
// 	if err := entity.GetDB().Model(&entity.Workout{}).Where("plan_id = ?", item.ID).Count(&count).Error; err != nil {
// 		log.Logger().Errorf("cant count workout for plan: %v", item.ID)
// 		return 0
// 	}
// 	return count
// }

// type _count struct {
// 	Count uint32
// }

// func completedWorkoutCount(pro *entity.Profile, item *entity.Plan) uint32 {
// 	var count _count
// 	//TODO: move this to entity package as a func
// 	if err := entity.GetDB().Table("mefit_workout_feedback").Select("count(*) as Count").Joins("LEFT JOIN mefit_workout on mefit_workout.id = mefit_workout_feedback.workout_id").
// 		Where("mefit_workout.plan_id = ?", item.ID).Scan(&count).Error; err != nil {
// 		log.Logger().Errorf("cant count workout history for plan: %v", item.ID)
// 		return 0
// 	}
// 	return count.Count
// }

func currentPlanFrom(pro *entity.Profile, item *entity.Plan) (*pb.CurrentPlanRes, error) {
	warray := 0
	if item.WorkoutArray != nil {
		warray = len(item.WorkoutArray)
	}
	completedWorkouts := warray
	if pro.CurrentWorkoutNo != 0 {
		completedWorkouts = pro.CurrentWorkoutNo - 1
	}
	var totalCalories uint32
	workouts := []entity.Workout{}
	if err := entity.GetDB().Model(&entity.Workout{}).Preload("WorkoutType").Where(([]int64)(item.WorkoutArray)).Find(&workouts).Error; err != nil {
	}
	log.Logger().Print("workoutsss length", len(workouts))
	for i := 0; i < len(workouts); i++ {
		totalCalories += uint32(workouts[i].Calorie)
	}
	plan := &pb.CurrentPlanRes{
		Plan: &pb.PlanRes{
			ThumbnailUrl:  item.ThumbnailUrl,
			Name:          item.Name,
			Id:            uint32(item.ID),
			Weeks:         uint32(item.Weeks),
			Level:         uint32(item.Level),
			WorkoutCounts: uint32(warray),
			TotalCalorie:  totalCalories,
			CoverUrl:      item.CoverUrl,
		},
		TotalWorkouts:     uint32(warray),
		CompletedWorkouts: uint32(completedWorkouts),
	}
	workout := entity.Workout{}
	if warray < pro.CurrentWorkoutNo {
		workout.ID = 0
	} else {
		workout.ID = uint(item.WorkoutArray[pro.CurrentWorkoutNo])
	}

	if err := entity.SimpleCrud(workout).Get(&workout, "WorkoutType"); err != nil {
		log.Logger().Errorf("cant get current workout for plan: %v, profile: %v", item.ID, pro.ID)
		return nil, utils.ErrNotFound
	}
	lock := false
	if item.VipLock == true {
		product := entity.Product{PlanID: item.ID}
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
	// log.Logger().Print("workout id is : ", pro.CurrentWorkoutNo, " workout array id is : ", item.WorkoutArray[0])
	if item.HasTrial == true && pro.CurrentWorkoutNo == int(item.WorkoutArray[0]) {
		lock = false
	}
	plan.Workout = &pb.Workout{
		Id:          uint32(workout.ID),
		Name:        workout.WorkoutType.Name,
		Description: workout.WorkoutType.Description,
		Calorie:     uint32(workout.Calorie),
		Duration:    uint32(workout.Duration),
		Viplock:     lock,
	}
	return plan, nil
}

package entity

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//User all users are organizations and there is only one Admin Service Account
type User struct {
	gorm.Model
	Anonymous   bool
	AnonymousID string `gorm:"varchar(32);unique"`
	// Username     string `gorm:"varchar(100);not null;unique" structs:"-"`
	Email    string `gorm:"varchar(100);unique" structs:"-"`
	Password string `gorm:"varchar(100)"`
	//TODO: complete vip
	VIP               bool
	VIPExpirationDate *time.Time
	ConfirmToken      string
	Confirmed         bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time
}

//ActivityLevel ranges from 1 to 5
type ActivityLevel uint32

const (

	//Rarely (1)
	Rarely ActivityLevel = iota + 1
	//Light (2): 1-3 times a week
	Light
	//Moderate (3): 3-5 times a week
	Moderate
	//Heavy (4): 6-7 times a week
	Heavy
	//Elite (5): Several times per day
	Elite
)

//Goal ranges from 1 to 3
type Goal uint32

const (
	//Lose fat (1)
	Lose Goal = iota + 1
	//Maintain or Get fitter (2)
	Maintain
	//Gain muscle (3)
	Gain
)

//DaysToWorkout ranges from 1 to 3
type DaysToWorkout uint32

const (
	//One2Three 1 to 3 times a week
	One2Three DaysToWorkout = iota + 1
	//Three2Five 3 to 5 times a week
	Three2Five
	//Five2Seven 5 to 7 times a week
	Five2Seven
)

type Gender uint32

const (
	Male Gender = iota + 1
	Female
	Others
)

func (in Profile) GetGender() string {
	switch in.Gender {
	case Male:
		return "Male"
	case Female:
		return "Female"
	default:
		return "Others"
	}
}

type Profile struct {
	gorm.Model
	User          User
	UserID        uint `gorm:"unique_index"`
	Name          string
	Gender        Gender `gorm:"default:'3'"`
	Age           uint32
	Height        uint32
	Waist         uint32
	Neck          uint32
	Hip           uint32
	Arm           uint32
	Leg           uint32
	CurrentWeight float32
	TargetWeight  float32
	ActivityLevel ActivityLevel
	//Days of week that user working out
	DaysToWorkout DaysToWorkout
	Goal          Goal
	//Current plan for this profile, might be null!
	Plan             *Plan
	PlanID           uint `gorm:"index"`
	CurrentWorkoutNo int  `gorm:"index; default:1"`
	// WorkoutFeedbacks []WorkoutFeedback
}

// the purchasable items in the application
type Product struct {
	gorm.Model
	Name  string `gorm:"unique"`
	SKU   string
	Price uint
	//TODO: change to *Plan if nesseccery
	Plan        *Plan
	PlanID      uint `gorm:"index"`
	ExpireTime  uint
	Description string
	Off         uint
	// to be used if the product was a special product
	Label       string
	Recommended bool
	CoverUrl    string
}

type Payment struct {
	gorm.Model
	UserID     uint   `gorm:"index"`
	Authority  string `gorm:"unique"`
	PaymentURL string `gorm:"unique"`
	Plan       *Plan
	PlanID     uint `gorm:"index"`
	Refrence   string
	//id of the purchased product
	ProductID   uint `gorm:"index"`
	ProductName string
	Status      uint
	Paid        bool
	Price       uint
	ExpireTime  *time.Time
}

type BazaarPayment struct {
	gorm.Model
	UserID     uint `gorm:"index"`
	Token      string
	Plan       *Plan
	PaymentURL string
	PlanID     uint `gorm:"index"`
	//id of the purchased product
	ProductID   uint `gorm:"index"`
	ProductName string
	Status      uint
	DevPayload  string
	Paid        bool
	Price       uint
	RSA         string
	SKU         string
	ExpireTime  *time.Time
}

type Keyword struct {
	gorm.Model
	Key string `gorm:"unique;not null"`
	// Programs  []Program  `gorm:"many2many:keyword_program;"`
	Classes []Class `gorm:"many2many:keyword_class;"`
	// Articles  []Article  `gorm:"many2many:keyword_article;"`
	Movements []Movement `gorm:"many2many:keyword_movement;"`
}

type MuscleGroup struct {
	gorm.Model
	Name      string     `gorm:"unique;not null"`
	Classes   []Class    `gorm:"many2many:mg_class;"`
	Movements []Movement `gorm:"many2many:mg_movement;"`
	// Articles  []Article  `gorm:"many2many:mg_article;"`
}

// type Article struct {
// 	gorm.Model
// 	ThumbnailUrl string
// 	CoverUrl     string
// 	Slug         string     `gorm:"unique;not null"`
// 	Movements    []Movement `gorm:"many2many:article_movements;"`
// 	// HTML body
// 	Body         string        `gorm:"type:text;not null"`
// 	Keywords     []Keyword     `gorm:"many2many:keyword_article;"`
// 	MuscleGroups []MuscleGroup `gorm:"many2many:mg_article;"`
// 	ShareURL     string
// 	// publish2.Version
// 	// publish2.Schedule
// 	// publish2.Visible
// }

type Class struct {
	gorm.Model
	Name         string `gorm:"unique;not null"`
	CoverUrl     string
	Description  string        `gorm:"text;not null"`
	Keywords     []Keyword     `gorm:"many2many:keyword_class;"`
	MuscleGroups []MuscleGroup `gorm:"many2many:mg_class;"`
	Movements    []Movement    //FIXME: better to use `gorm:"many2many:class_movements;"`
}

// type FaultNFix{

// }

type Movement struct {
	gorm.Model
	Name         string `gorm:"unique;not null"`
	NameFa       string `gorm:"unique;not null"`
	ThumbnailUrl string
	// Description for movement
	Description string `gorm:"text;not null"`
	// Instruction on how to do the movement
	Instruction string `gorm:"text;not null"`
	Tips        string `gorm:"text;not null"`
	VideoUrl    string
	//FIXME: better to use Classes      []Class       `gorm:"many2many:class_movements;"`
	Class   Class
	ClassID uint `gorm:"index"`
	// Articles     []Article     `gorm:"many2many:article_movements;"`
	MuscleGroups []MuscleGroup `gorm:"many2many:mg_movement;"`
	Keywords     []Keyword     `gorm:"many2many:keyword_movement;"`
	// FaultNFix FaultNFix
}

type Difficulty uint32

const (
	Easy Difficulty = iota + 1
	Suitable
	Hard
)

type WorkoutHistory struct {
	gorm.Model
	Workout   *Workout
	WorkoutID uint
	Profile   *Profile
	ProfileID uint
	//Rating is between 1 to 5
	Rating     uint
	Difficulty Difficulty
}

type PurchasedProduct struct {
	gorm.Model
	Product     Product
	ProductID   uint
	Profile     Profile
	ProfileID   uint
	PaymentType string
	PaymentID   uint
}

// type Program struct {
// 	gorm.Model
// 	Title       string `gorm:"unique;not null"`
// 	Description string `gorm:"text;not null"`
// 	CoverUrl    string
// 	Keywords    []Keyword `gorm:"many2many:keyword_program;"`
// 	Levels      []Level   `gorm:"foreignkey:ProgramID;"`
// }
// type Level struct {
// 	gorm.Model
// 	Name        string `gorm:"unique;not null"`
// 	Description string
// 	SubLevels   []SubLevel
// 	// //hhio618 Added this for redundancy
// 	// Plans     []Plan
// 	No        uint `gorm:"unique_index;not null"`
// 	ProgramID uint `gorm:"index"`
// 	Program   Program
// }

// type SubLevel struct {
// 	gorm.Model
// 	Name    string `gorm:"unique;not null"`
// 	LevelID uint   `gorm:"index"`
// 	Level   Level
// 	No      uint `gorm:"index;not null"`
// 	Plans   []Plan
// }

type Plan struct {
	gorm.Model
	ThumbnailUrl string
	CoverUrl     string
	Name         string `gorm:"not null"`
	Description  string `gorm:"text;not null"`
	//What is the priority of this plan in the list
	No uint `gorm:"index;not null"`
	// how many weeks it takes
	Weeks uint `gorm:"not null"`
	// requirements e.g. some gym tools
	SubLevelID uint `gorm:"index"`
	Level      uint
	MetaName   string
	// SubLevel   SubLevel
	// Workouts []Workout
	WorkoutArray pq.Int64Array `gorm:"type:integer[]"`

	VipLock  bool
	HasTrial bool
}

type WorkoutType struct {
	gorm.Model
	Name        string `gorm:"text;not null"`
	Description string `gorm:"text;not null"`
}

type Workout struct {
	gorm.Model
	Plan          *Plan
	PlanID        uint `gorm:"index"`
	Name          string
	WorkoutType   *WorkoutType
	WorkoutTypeID uint
	// Description string `gorm:"text;not null"`
	Instruction string `gorm:"text;not null"`
	Calorie     uint   `gorm:"not null"`
	//duration in second
	Duration         uint `gorm:"not null"`
	ExerciseSections []ExerciseSection
	VipLock          bool
}

type ExerciseSection struct {
	gorm.Model
	WorkoutID uint `gorm:"index"`
	Workout   Workout
	Exercises []Exercise
	Round     uint `gorm:"not null"`
}

type ExerciseType uint32

const (
	Duration ExerciseType = iota + 1
	Repitation
)

type Exercise struct {
	gorm.Model
	ExerciseSection   ExerciseSection
	ExerciseSectionID uint `gorm:"index"`
	MovementID        uint `gorm:"index"`
	Movement          *Movement
	// for duration(1) or repitation (2);
	ExerciseType ExerciseType `gorm:"not null"`
	//Value is second for duration, or count for rep. e.g. 10x pushups, 120 sec of running
	Value uint `gorm:"not null"`
	//Rest value zero means no rest included
	Rest *int32 `gorm:"not null"`
}

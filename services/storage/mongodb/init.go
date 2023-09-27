package mongodb

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"gitlab.com/mefit/mefit-server/utils"

	"gitlab.com/mefit/mefit-server/utils/assert"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Image struct {
	Small       string
	Medium      string
	Large       string
	Original    string
}
type Article struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string
	Slug      string
	Body      string
	Thumb     string
	Sum       string             `bson:"sum"`
	Time      uint64             `bson:"time"`
	Author    string             `bson:"author"`
	Images    Image
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetArticles(page int) (out []Article, err error) {
	out = []Article{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Logger().Errorf("In fetching article: connecting: %v", err)
		return nil, utils.ErrInternal
	}
	collection := client.Database(db).Collection("courses")
	projection := bson.D{
		{"_id", 1},
		{"title", 1},
		{"time", 1},
		{"slug", 1},
		{"images", 1},
		{"time", 1},
		{"body", 1},
		{"thumb", 1},
		{"sum", 1},
		{"createdAt", 1},
		{"updatedAt", 1},
	}
	sort := bson.D{
		{"_id", -1},
	}
	cur, err := collection.Find(context.Background(), bson.D{}, options.Find().SetSkip(10 * int64(page - 1)).SetLimit(10).SetProjection(projection).SetSort(sort))
	if err != nil {
		log.Logger().Errorf("While fetching article cursor: %v", err)
		return nil, utils.ErrInternal
	}
	defer cur.Close(ctx)

	var result Article
	for cur.Next(ctx) {
		_ = cur.Decode(&result)
		if err != nil {
			log.Logger().Errorf("While looping over article cursor: %v", err)
			return nil, utils.ErrInternal
		}

		// do something with result....
		out = append(out, result)
	}
	if err := cur.Err(); err != nil {
		log.Logger().Errorf("While closing article cursor: %v", err)
		return nil, utils.ErrInternal
	}
	return out, nil
}

func GetArticle(id string) (out *Article, err error) {
	out = &Article{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Logger().Errorf("In fetching single article: connecting: %v", err)
		return nil, utils.ErrInternal
	}
	collection := client.Database(db).Collection("courses")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(out)
	if err != nil {
		log.Logger().Errorf("While fetching article: %v", err)
		return nil, utils.ErrInternal
	}
	return out, nil
}

func ArticleCount() (count int, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Logger().Errorf("In fetching single article: connecting: %v", err)
		return 0, utils.ErrInternal
	}
	collection := client.Database(db).Collection("courses")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Logger().Errorf("While fetching article: %v", err)
		return 0, utils.ErrInternal
	}
	defer cur.Close(context.Background())
	articleCount := 0
	for cur.Next(context.Background()) {
		articleCount++
	}
	return articleCount, nil
}

func getMongoURI() (string, string) {
	log.Logger().Print("Parsing mongodb credentials")
	var (
		uri string
		db  string
	)
	if _, ok := os.LookupEnv("VCAP_SERVICES"); ok {
		vcapServices := make(map[string]interface{})
		json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &vcapServices)
		creds := vcapServices["mongodb"].([]interface{})[0].(map[string]interface{})["credentials"].(map[string]interface{})
		uri = creds["uri"].(string)
	} else {
		uri = config.Config().GetString("mongodb_uri")
	}
	parts := strings.Split(uri, "/")
	db = parts[len(parts)-1]
	return uri, db
}

type manager struct{}

var (
	uri string
	db  string
)

func (manager) Initialize() func() {
	uri, db = getMongoURI()
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	assert.Nil(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	assert.Nil(err)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	assert.Nil(err)
	log.Logger().Info("Connected to MongoDB!")
	return nil
}

func init() {
	initializer.Register(manager{}, initializer.VeryHighPriority)
}

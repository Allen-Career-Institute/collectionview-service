package mongo

import (
	"collectionview-service/internal/conf"
	"collectionview-service/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ProviderSet is mongoDB providers.
var ProviderSet = wire.NewSet(NewData, NewMongoCollectionImpl, NewLogger)

// Data .
type Data struct {
	client *mongo.Client
}

type DatabaseCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLogger() *log.Helper {
	logger := log.NewStdLogger(os.Stdout)
	logHelper := log.NewHelper(logger)
	return logHelper
}

// Create a MongoDB client
func NewData(c *conf.Data, log *log.Helper) (*Data, error) {
	credentials := ReadCredentials(c.Mongo.CredFileLocation)
	if credentials == nil {
		err := errors.New(500, utils.MongoErrorTag, "Error fetching mongo credentials")
		return nil, err
	}
	credential := options.Credential{
		Username: credentials.Username,
		Password: credentials.Password,
	}
	log.Info("Started connecting with mongo server")
	ctx := context.TODO()
	//options.Client().ApplyURI(c.Mongo.GetConnection()).SetAuth(credential).SetMonitor(otelmongo.NewMonitor())
	clientOptions := options.Client().ApplyURI(c.Mongo.GetConnection()).SetAuth(credential).SetMonitor(otelmongo.NewMonitor())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Error while connecting to mongo db")
		return nil, errors.New(500, utils.MongoErrorTag, fmt.Sprintf("Error on connecting Mongo server: %v", err))
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Errorf("Error while connecting to mongo db")
		return nil, errors.New(500, utils.MongoErrorTag, fmt.Sprintf("Error pinging Mongo server: %v", err))
	}
	log.Info("Completed  connecting with mongo server successfully")
	return &Data{client: client}, nil
}

func ReadCredentials(fileName string) *DatabaseCreds {
	// read our opened jsonFile as a byte array.
	byteValue, err := os.ReadFile(fileName)

	if err != nil {
		log.Errorf("Error while reading the file for mongo creds")
		return nil
	}

	// we initialize our Users array
	var creds DatabaseCreds

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &creds)

	if err != nil {
		log.Errorf("Error while marshaling  the file for mongo creds")
		return nil
	}

	return &creds
}

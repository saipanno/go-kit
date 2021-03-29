package client

import (
	"context"
	"reflect"
	"time"

	"github.com/saipanno/go-kit/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// CreateMongoConn ...
func CreateMongoConn(conf *DBConfig) (client *mongo.Client, err error) {

	logger.Infof("create mongo connect %s", conf.URI)

	rb := bson.NewRegistryBuilder()
	option := options.Client().ApplyURI(conf.URI)

	rb.RegisterTypeDecoder(reflect.TypeOf(time.Time{}),
		bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true)))

	option.SetRegistry(rb.Build())

	client, err = mongo.Connect(context.Background(), option)
	if err != nil {
		logger.Errorf("create mongodb session failed, message is %s", err.Error())
		return
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Errorf("ping mongodb failed, message is %s", err.Error())
	}

	return
}

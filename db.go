package xlib

import (
	"context"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type DBOptions struct {
	Start     int64
	Offset    int64
	Order     string
	OrderType int
}

//IDBConn database interface
type IDBConn interface {
	Find(table string, field string, value interface{}, res interface{}) error
	FindMany(table string, field string, value interface{}, dbopt *DBOptions, res interface{}) error
	Insert(table string, value interface{}) (primitive.ObjectID, error)
	Update(table string, data interface{}, field string, value interface{}) (int64, error)
	Delete(table string, field string, value interface{}) (int64, error)
	CreateUniqueIndex(collection string, keys ...string) error
}

type dbConn struct {
	session *mongo.Client
	db      *mongo.Database
}

func (conn *dbConn) FindMany(table string, field string, value interface{}, dbopt *DBOptions, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &options.FindOptions{}
	if dbopt.Offset > 0 {
		opts.Skip = &dbopt.Start
		opts.Limit = &dbopt.Offset
	}
	opts.SetSort(bson.D{{dbopt.Order, dbopt.OrderType}})

	filter := bson.M{}
	if field != "" {
		filter = bson.M{field: value}
	}
	cur, err := conn.db.Collection(table).Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cur.All(ctx, res)
}

func (conn *dbConn) Find(table string, field string, value interface{}, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{field: value}
	return conn.db.Collection(table).FindOne(ctx, filter).Decode(res)
}

func (conn *dbConn) Insert(table string, value interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.db.Collection(table).InsertOne(ctx, value)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (conn *dbConn) Delete(table string, field string, value interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{field: value}
	res, err := conn.db.Collection(table).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (conn *dbConn) Update(table string, data interface{}, field string, value interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{field: value}

	pByte, err := bson.Marshal(data)
	if err != nil {
		return 0, err
	}
	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return 0, err
	}
	upd := bson.D{{Key: "$set", Value: update}}
	res, err := conn.db.Collection(table).UpdateOne(ctx, filter, upd)
	if err != nil {
		return 0, err
	}
	return res.MatchedCount, nil
}

//CreateUniqueIndex create a unique index
func (conn *dbConn) CreateUniqueIndex(collection string, keys ...string) error {
	db := conn.db.Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}

	// Composite index
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// Create index
	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		opts,
	)
	if result == "" || err != nil {
		return err
	}
	return nil
}

//InitDatabase init database connection
func InitDatabase() (IDBConn, func(), error) {

	//Get enviromental variables
	connMongo := os.Getenv("MONGO_DSN")
	if len(connMongo) == 0 {
		connMongo = "mongodb://127.0.0.1:27017"
	}

	dbMongo := os.Getenv("MONGO_DB")
	if len(dbMongo) == 0 {
		dbMongo = os.Args[0]
	}

	//Create session
	session, err := mongo.NewClient(options.Client().ApplyURI(connMongo))
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Connect to mongo
	err = session.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Return disconnect as func to defer in main
	f := func() { session.Disconnect(ctx) }

	//Ping Database
	err = session.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, f, err
	}
	db := session.Database(dbMongo)
	return &dbConn{session: session, db: db}, f, nil

}

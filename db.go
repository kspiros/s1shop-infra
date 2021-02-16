package xlib

import (
	"context"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

//DBOptions database options
type DBOptions struct {
	Start     int64
	Offset    int64
	Order     string
	OrderType int
}

//IDBConn database interface
type IDBConn interface {
	FindOne(table string, sel interface{}, filters map[string]interface{}, res interface{}) error
	FindMany(table string, sel interface{}, filters map[string]interface{}, dbopt *DBOptions, res interface{}) error
	FindOneAndUpdate(table string, sel interface{}, filters map[string]interface{}, update interface{}, res interface{}) error
	FindOneAndDelete(table string, sel interface{}, filter interface{}, res interface{}) error
	InsertOne(table string, document interface{}) (interface{}, error)
	InsertMany(table string, documents []interface{}) ([]interface{}, error)
	UpdateOne(table string, filter interface{}, update interface{}) (int64, error)
	UpdateMany(table string, filter interface{}, update interface{}) (int64, error)
	BulkReplace(table string, filters []interface{}, documents []interface{}) (int64, error)
	DeleteOne(table string, filter interface{}) (int64, error)
	DeleteMany(table string, filter interface{}) (int64, error)
	Aggregate(table string, pipeline interface{}, res interface{}) error
	CreateUniqueIndex(collection string, keys ...string) error
}

type dbConn struct {
	session *mongo.Client
	db      *mongo.Database
}

func (conn *dbConn) FindMany(table string, sel interface{}, filters map[string]interface{}, dbopt *DBOptions, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &options.FindOptions{}
	if dbopt.Offset > 0 {
		opts.Skip = &dbopt.Start
		opts.Limit = &dbopt.Offset
	}
	opts.SetSort(bson.D{{Key: dbopt.Order, Value: dbopt.OrderType}})
	if sel != nil {
		opts.SetProjection(sel)
	}

	cur, err := conn.db.Collection(table).Find(ctx, filters, opts)
	if err != nil {
		return err
	}

	return cur.All(ctx, res)
}

func (conn *dbConn) FindOneAndUpdate(table string, sel interface{}, filters map[string]interface{}, update interface{}, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &options.FindOneAndUpdateOptions{}
	if sel != nil {
		opts.SetProjection(sel)
	}

	return conn.db.Collection(table).FindOneAndUpdate(ctx, filters, update, opts).Decode(res)
}

func (conn *dbConn) FindOne(table string, sel interface{}, filters map[string]interface{}, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &options.FindOneOptions{}
	if sel != nil {
		opts.SetProjection(sel)
	}

	return conn.db.Collection(table).FindOne(ctx, filters, opts).Decode(res)
}

func (conn *dbConn) InsertMany(table string, documents []interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.db.Collection(table).InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

func (conn *dbConn) InsertOne(table string, document interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.db.Collection(table).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (conn *dbConn) DeleteOne(table string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := conn.db.Collection(table).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (conn *dbConn) FindOneAndDelete(table string, sel interface{}, filter interface{}, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &options.FindOneAndDeleteOptions{}
	if sel != nil {
		opts.SetProjection(sel)
	}

	return conn.db.Collection(table).FindOneAndDelete(ctx, filter, opts).Decode(res)
}

func (conn *dbConn) DeleteMany(table string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := conn.db.Collection(table).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (conn *dbConn) UpdateOne(table string, filter interface{}, update interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.db.Collection(table).UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.MatchedCount, nil
}

func (conn *dbConn) UpdateMany(table string, filter interface{}, update interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.db.Collection(table).UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.MatchedCount, nil
}

func (conn *dbConn) BulkReplace(table string, filters []interface{}, documents []interface{}) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var operations []mongo.WriteModel

	for i, doc := range documents {
		m := mongo.NewReplaceOneModel()
		m.SetFilter(filters[i])
		m.Replacement = doc
		operations = append(operations, m)
	}

	res, err := conn.db.Collection(table).BulkWrite(ctx, operations)
	if err != nil {
		return 0, err
	}
	return res.MatchedCount, nil
}

func (conn *dbConn) Aggregate(table string, pipeline interface{}, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := conn.db.Collection(table).Aggregate(ctx, pipeline)
	defer cur.Close(ctx)

	if err != nil {
		return err
	}

	return cur.All(ctx, res)
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
	f := func() {
		session.Disconnect(ctx)
	}

	//Ping Database
	err = session.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, f, err
	}
	db := session.Database(dbMongo)
	return &dbConn{session: session, db: db}, f, nil

}

package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
	dbName string
}

func NewMongoStorage(client *mongo.Client, dbName string) *MongoStorage {
	return &MongoStorage{
		client: client,
		dbName: dbName,
	}
}

func (s *MongoStorage) PutDocument(ctx context.Context, collectionName string, document map[string]any) error {
	coll := s.client.Database(s.dbName).Collection(collectionName)

	var filter bson.M
	if id, ok := document["_id"]; ok {
		filter = bson.M{"_id": id}
	} else if id, ok := document["id"]; ok {
		filter = bson.M{"id": id}
	} else {
		_, err := coll.InsertOne(ctx, document)
		return err
	}

	update := bson.M{"$set": document}
	opts := options.Update().SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *MongoStorage) GetDocument(ctx context.Context, collectionName string, filter map[string]any) (map[string]any, error) {
	coll := s.client.Database(s.dbName).Collection(collectionName)

	var doc map[string]any
	err := coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return doc, nil
}

func (s *MongoStorage) ListDocuments(ctx context.Context, collectionName string, filter map[string]any) ([]map[string]any, error) {
	coll := s.client.Database(s.dbName).Collection(collectionName)

	f := bson.M{}
	if len(filter) > 0 {
		f = filter
	}

	cursor, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []map[string]any
	for cursor.Next(ctx) {
		var doc map[string]any
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if documents == nil {
		documents = []map[string]any{}
	}
	return documents, nil
}

func (s *MongoStorage) DeleteDocument(ctx context.Context, collectionName string, filter map[string]any) (int64, error) {
	coll := s.client.Database(s.dbName).Collection(collectionName)
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (s *MongoStorage) CreateCollection(ctx context.Context, collectionName string) error {
	return s.client.Database(s.dbName).CreateCollection(ctx, collectionName)
}

func (s *MongoStorage) ListCollections(ctx context.Context) ([]string, error) {
	collections, err := s.client.Database(s.dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if collections == nil {
		collections = []string{}
	}
	return collections, nil
}

func (s *MongoStorage) DeleteCollection(ctx context.Context, collectionName string) error {
	return s.client.Database(s.dbName).Collection(collectionName).Drop(ctx)
}

func (s *MongoStorage) CreateIndex(ctx context.Context, collectionName string, fieldName string, unique bool, sparse bool) (string, error) {
	coll := s.client.Database(s.dbName).Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: fieldName, Value: 1}},
	}

	if unique || sparse {
		opts := options.Index()
		if unique {
			opts.SetUnique(true)
		}
		if sparse {
			opts.SetSparse(true)
		}
		indexModel.Options = opts
	}

	return coll.Indexes().CreateOne(ctx, indexModel)
}

func (s *MongoStorage) DeleteIndex(ctx context.Context, collectionName string, indexName string) error {
	_, err := s.client.Database(s.dbName).Collection(collectionName).Indexes().DropOne(ctx, indexName)
	return err
}

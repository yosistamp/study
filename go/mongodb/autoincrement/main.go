package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AutoIncrement struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Next int                `bson:"next"`
}

type YourDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

func main() {
	// credential
	credential := options.Credential{
		Username: "root",
		Password: "password",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// MongoDBに接続
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 自動増分用のコレクション
	autoIncrementCollection := client.Database("mydb").Collection("autoIncrement")
	count, err := autoIncrementCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if count < 1 {
		// 自動増分の初期値を挿入
		_, err = autoIncrementCollection.InsertOne(ctx, AutoIncrement{Next: 1})
		if err != nil {
			log.Fatal(err)
		}
	}

	// ドキュメントを挿入するための関数
	insertDocument := func(doc YourDocument) error {
		// 自動増分値を取得
		var autoInc AutoIncrement
		err := autoIncrementCollection.FindOneAndUpdate(
			ctx,
			bson.M{},
			bson.M{"$inc": bson.M{"next": 1}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&autoInc)
		if err != nil {
			return err
		}

		// ドキュメントに自動増分値を設定
		doc.ID = autoInc.Next

		// ドキュメントを挿入
		_, err = client.Database("mydb").Collection("yourCollection").InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		return nil
	}

	// ドキュメントを挿入
	err = insertDocument(YourDocument{Name: "Example 1"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ドキュメントが挿入されました。")
}

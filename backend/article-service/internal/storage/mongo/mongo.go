package mongo

import (
	"context"
	"fmt"
	"time"

	articleService "github.com/NekNB/CyberNavigate/backend/article-service/internal/services/article"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	log               *logrus.Logger
	Client            *mongo.Client
	DB                *mongo.Database
	ArticleCollection *mongo.Collection
}

type Article struct {
	ID   primitive.ObjectID `bson:"_id"`
	Text string             `bson:"text"`
}

var _ articleService.ArticleContentProvider = (*MongoStorage)(nil)

func CreateConnection(log *logrus.Logger, uri, dbName, collectionName string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаем клиент
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.WithError(err).Error("Ошибка подключения к MongoDB")
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверяем подключение
	if err = client.Ping(ctx, nil); err != nil {
		log.WithError(err).Error("MongoDB не отвечает на ping")
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Info("Успешно подключено к MongoDB")

	// Инициализируем базу данных и коллекции
	db := client.Database(dbName)
	articleCollection := db.Collection(collectionName)

	return &MongoStorage{
		log:               log,
		Client:            client,
		DB:                db,
		ArticleCollection: articleCollection,
	}, nil
}

func (m *MongoStorage) SaveText(ctx context.Context, articleText string) (string, error) {

	m.log.Debug(articleText)
	res, err := m.ArticleCollection.InsertOne(ctx, bson.D{{Key: "text", Value: articleText}})
	if err != nil {
		return "", fmt.Errorf("insert: %w", err)
	}

	insertID := res.InsertedID.(primitive.ObjectID)

	return insertID.Hex(), nil
}

func (m *MongoStorage) UpdateText(ctx context.Context, id string, articleText string) error {
	textID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.log.Error(err)
		return err
	}
	filter := bson.M{"_id": textID}
	update := bson.M{"$set": bson.M{"text": articleText}}

	// SetUpsert(true) создаст документ, если его нет
	opts := options.Update().SetUpsert(false)
	_, err = m.ArticleCollection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (m *MongoStorage) ArticleTextById(ctx context.Context, id string) (string, error) {
	var article Article
	textID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.log.Error(err)
		return "", nil
	}

	filter := bson.M{"_id": textID}

	err = m.ArticleCollection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		m.log.Error(err)
		if err == mongo.ErrNoDocuments {
			return "", nil // Не найдено
		}
		return "", fmt.Errorf("find one: %w", err)
	}
	return article.Text, nil
}

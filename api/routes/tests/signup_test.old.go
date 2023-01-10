package tests

import (
	"accounts.sidooh/api"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var collection *mongo.Collection
var ctx context.Context

func TestMain(m *testing.M) {
	mongoServer, err := memongo.Start("5.0.4")
	if err != nil {
		log.Fatal(err)
	}
	defer mongoServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		panic(err)
	}

	collection = client.Database("auth").Collection("users")

	e, port, s := api.Setup()
	e.Logger.Fatal(e.StartH2CServer(":"+port, s))
	defer e.Close()

	os.Exit(m.Run())
}

func TestSignUp(t *testing.T) {

	collection.InsertOne(context.TODO(), bson.D{{"name", "Alice"}})

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	userJSON := `{"name":"Jon Snow","email":"jon@labstack.com"}`

	fmt.Println("Sending payload")

	res, err := http.Post("http://localhost:35663/api/users/signup", echo.MIMEApplicationJSON, strings.NewReader(userJSON))
	fmt.Println(err)
	fmt.Println(res)

	//req := httptest.NewRequest(http.MethodPost, "/api/users/signup", strings.NewReader(userJSON))
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//
	//rec := httptest.NewRecorder()
	//
	//err = handlerFunc(e.NewContext(req, rec))
	//
	//if assert.NoError(t, err) {
	assert.Equal(t, http.StatusCreated, res)
	assert.Equal(t, userJSON, res.Header)
	//}
}

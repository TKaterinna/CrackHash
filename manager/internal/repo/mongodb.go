package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestStatusDoc struct {
	ID         uuid.UUID       `bson:"_id"`
	TasksCount int             `bson:"tasks_count"`
	TasksReady map[string]bool `bson:"tasks_ready"`
	StartTime  time.Time       `bson:"start_time"`
	Status     string          `bson:"status"`
	Results    []string        `bson:"results"`
}

type WorkerTasksDoc struct {
	ID         uuid.UUID `bson:"_id"`
	RequestID  uuid.UUID `bson:"request_id"`
	StartIndex int64     `bson:"start_index"`
	Count      int64     `bson:"count"`
	Alphabet   string    `bson:"alphabet"`
	MaxLen     int64     `bson:"max_len"`
	TargetHash string    `bson:"target_hash"`
}

type MongoTaskRepo struct {
	client      *mongo.Client
	requests    *mongo.Collection
	workerTasks *mongo.Collection
	errorDelay  time.Duration
}

func InitDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	fmt.Println("MongoDB connected successfully")

	return client, nil
}

func NewMongoTaskRepo(uri string, dbName string, errorDelay time.Duration) *MongoTaskRepo {
	client, err := InitDB(uri)
	if err != nil {
		panic("mongo client creation failed")
	}

	db := client.Database(dbName)

	_, _ = db.Collection("requests").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{Keys: bson.D{{"_id", 1}}, Options: options.Index().SetUnique(true)},
	)
	_, _ = db.Collection("worker_tasks").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{Keys: bson.D{{"request_id", 1}}, Options: options.Index()},
	)

	return &MongoTaskRepo{
		client:      client,
		requests:    db.Collection("requests"),
		workerTasks: db.Collection("worker_tasks"),
		errorDelay:  errorDelay,
	}
}

func (r *MongoTaskRepo) Close() error {
	return r.client.Disconnect(context.Background())
}

func (r *MongoTaskRepo) SaveRequest(id uuid.UUID, tasks []*models.CrackTaskRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasksReady := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		tasksReady[t.TaskId.String()] = false
	}

	requestDoc := RequestStatusDoc{
		ID:         id,
		TasksCount: len(tasks),
		TasksReady: tasksReady,
		StartTime:  time.Now(),
		Status:     models.StatusInProgress,
		Results:    make([]string, 0),
	}

	_, err := r.requests.InsertOne(ctx, requestDoc)
	if err != nil {
		return fmt.Errorf("insert request status: %w", err)
	}

	if len(tasks) > 0 {
		var workerDocs []any
		for _, t := range tasks {
			workerDocs = append(workerDocs, WorkerTasksDoc{
				ID:         t.TaskId,
				RequestID:  id,
				StartIndex: t.StartIndex,
				Count:      t.Count,
				Alphabet:   t.Alphabet,
				MaxLen:     t.MaxLen,
				TargetHash: t.TargetHash,
			})
		}
		_, err = r.workerTasks.InsertMany(ctx, workerDocs)
		if err != nil {
			return fmt.Errorf("insert worker tasks: %w", err)
		}
	}

	return nil
}

func (r *MongoTaskRepo) GetStatus(id uuid.UUID) (string, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req RequestStatusDoc
	err := r.requests.FindOne(ctx, bson.M{"_id": id}).Decode(&req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil, fmt.Errorf("request not found")
		}
		return "", nil, err
	}

	status := req.Status
	var results []string

	if status == models.StatusREADY {
		results = req.Results
	} else {
		results = nil
		if time.Now().After(req.StartTime.Add(r.errorDelay)) && status != models.StatusERROR {
			_, _ = r.requests.UpdateOne(ctx,
				bson.M{"_id": id, "status": bson.M{"$ne": models.StatusERROR}},
				bson.M{"$set": bson.M{"status": models.StatusERROR}},
			)
			status = models.StatusERROR
		}
	}

	return status, results, nil
}

func (r *MongoTaskRepo) UpdateResult(reqId uuid.UUID, taskId uuid.UUID, results []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	taskKey := taskId.String()

	filter := bson.M{
		"_id":                                  reqId,
		"status":                               bson.M{"$ne": models.StatusERROR},
		fmt.Sprintf("tasks_ready.%s", taskKey): false,
	}

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("tasks_ready.%s", taskKey): true,
		},
		"$inc": bson.M{
			"tasks_count": -1,
		},
	}

	if len(results) > 0 {
		update["$push"] = bson.M{
			"results": bson.M{
				"$each": results,
			},
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // return docs version after updating
	var updated RequestStatusDoc

	err := r.requests.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var req RequestStatusDoc
			if findErr := r.requests.FindOne(ctx, bson.M{"_id": reqId}).Decode(&req); findErr == nil {
				if req.Status == models.StatusERROR {
					return fmt.Errorf("this task was canceled by timeout")
				}
				if req.TasksReady[taskKey] {
					return fmt.Errorf("duplicated task result")
				}
			}
			return fmt.Errorf("request not found")
		}
		return fmt.Errorf("update failed: %w", err)
	}

	if updated.TasksCount == 0 {
		_, _ = r.requests.UpdateOne(ctx,
			bson.M{"_id": reqId, "status": models.StatusInProgress},
			bson.M{"$set": bson.M{"status": models.StatusREADY}},
		)
	}

	return nil
}

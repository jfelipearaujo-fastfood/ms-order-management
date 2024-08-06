package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

var opts = godog.Options{
	Format:      "pretty",
	Paths:       []string{"features"},
	Output:      colors.Colored(os.Stdout),
	Concurrency: 4,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

// Steps
const featureKey CtxKeyType = "feature"

type feature struct {
	HostApi    string
	OrderId    string
	Items      []string
	StateTitle string
	Token      string
}

var state = NewState[feature](featureKey)

func iCreateAnOrder(ctx context.Context) (context.Context, error) {
	feat := state.retrieve(ctx)

	token, err := generateToken(uuid.NewString(), time.Minute*10)
	if err != nil {
		return ctx, err
	}

	body := `{
		"customer_id": "1387d7f1-732e-4ab4-8c0a-adb13b0d7797"
	}`

	route := fmt.Sprintf("%s/orders", feat.HostApi)
	req, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return ctx, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}

	if res.StatusCode != http.StatusCreated {
		return ctx, fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	defer res.Body.Close()

	var order map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return ctx, err
	}

	orderId, ok := order["id"].(string)
	if !ok {
		return ctx, fmt.Errorf("Order ID not found")
	}

	feat.Token = token
	feat.OrderId = orderId

	return state.enrich(ctx, feat), nil
}

func iAddedAnItemToTheOrder(ctx context.Context) (context.Context, error) {
	feat := state.retrieve(ctx)

	body := `{
		"items": [
			{
				"id": "b88014db-320d-4ac9-99b1-422774d56106",
				"name": "Test item",
				"unit_price": 10.5,
				"quantity": 1
			}
		]
	}`

	route := fmt.Sprintf("%s/orders/%s/items", feat.HostApi, feat.OrderId)
	req, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return ctx, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", feat.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}

	if res.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	return ctx, nil
}

func iRetrieveTheOrder(ctx context.Context) (context.Context, error) {
	feat := state.retrieve(ctx)

	route := fmt.Sprintf("%s/orders/%s", feat.HostApi, feat.OrderId)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return ctx, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", feat.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}

	if res.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("Expected status code 201, got %d", res.StatusCode)
	}

	defer res.Body.Close()

	var order map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return ctx, err
	}

	items, ok := order["items"].([]interface{})
	if !ok {
		return ctx, fmt.Errorf("Items not found")
	}

	for _, item := range items {
		itemId, ok := item.(map[string]interface{})["id"].(string)
		if !ok {
			return ctx, fmt.Errorf("Item ID not found")
		}
		feat.Items = append(feat.Items, itemId)
	}

	stateTitle, ok := order["state_title"].(string)
	if !ok {
		return ctx, fmt.Errorf("State not found")
	}

	feat.StateTitle = stateTitle

	return state.enrich(ctx, feat), nil
}

func theOrderShouldHaveTheItem(ctx context.Context) (context.Context, error) {
	feat := state.retrieve(ctx)

	if len(feat.Items) != 1 {
		return ctx, fmt.Errorf("Expected 1 item, got %d", len(feat.Items))
	}

	return ctx, nil
}

func theOrderStateShouldBe(ctx context.Context, stateTitle string) (context.Context, error) {
	feat := state.retrieve(ctx)

	if feat.StateTitle != stateTitle {
		return ctx, fmt.Errorf("Expected state %s, got %s", stateTitle, feat.StateTitle)
	}

	return ctx, nil
}

type testContext struct {
	network    *testcontainers.DockerNetwork
	containers []testcontainers.Container
}

var (
	containers = make(map[string]testContext)
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		network, err := network.New(ctx, network.WithCheckDuplicate(), network.WithDriver("bridge"))
		if err != nil {
			return ctx, err
		}

		postgresContainer, ctx, err := createPostgresContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		localstack, ctx, err := createLocalstackContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		apiContainer, ctx, err := createApiContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		containers[sc.Id] = testContext{
			network: network,
			containers: []testcontainers.Container{
				postgresContainer,
				localstack,
				apiContainer,
			},
		}

		return ctx, nil
	})

	ctx.Step(`^I create an order$`, iCreateAnOrder)
	ctx.Step(`^I added an item to the order$`, iAddedAnItemToTheOrder)
	ctx.Step(`^I retrieve the order$`, iRetrieveTheOrder)
	ctx.Step(`^the order should have the item$`, theOrderShouldHaveTheItem)
	ctx.Step(`^the order state should be "([^"]*)"$`, theOrderStateShouldBe)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			return ctx, err
		}

		tc := containers[sc.Id]

		for _, c := range tc.containers {
			err := c.Terminate(ctx)
			if err != nil {
				return ctx, err
			}
		}

		err = tc.network.Remove(ctx)

		return ctx, err
	})
}

func generateToken(userId string, expire time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(expire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("my-secret"))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenString), nil
}

func createApiContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../",
				Dockerfile: "Dockerfile",
				KeepImage:  true,
			},
			ExposedPorts: []string{
				"8080",
			},
			Env: map[string]string{
				"API_PORT":                     "8080",
				"API_ENV_NAME":                 "development",
				"API_VERSION":                  "v1",
				"DB_URL":                       "todo",
				"DB_URL_SECRET_NAME":           "db-secret-url",
				"AWS_ACCESS_KEY_ID":            "test",
				"AWS_SECRET_ACCESS_KEY":        "test",
				"AWS_REGION":                   "us-east-1",
				"AWS_BASE_ENDPOINT":            "http://test:4566",
				"AWS_ORDER_PAYMENT_TOPIC_NAME": "OrderPaymentTopic",
				"AWS_UPDATE_ORDER_QUEUE_NAME":  "UpdateOrderQueue",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			WaitingFor: wait.ForLog("Server started").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, ctx, err
	}

	ports, err := container.Ports(ctx)
	if err != nil {
		return nil, ctx, err
	}

	if len(ports["8080/tcp"]) == 0 {
		return nil, ctx, fmt.Errorf("Port 8080/tcp not found")
	}

	port := ports["8080/tcp"][0].HostPort

	res, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		return nil, ctx, err
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ctx, err
		}
		defer res.Body.Close()

		fmt.Printf("Body: %s", string(body))

		return nil, ctx, fmt.Errorf("API health check failed with status: %d", res.StatusCode)
	}

	return container, state.enrich(ctx, &feature{
		HostApi: fmt.Sprintf("http://localhost:%s/api/v1", port),
	}), nil
}

func createLocalstackContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	snsScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sns.sh"))
	if err != nil {
		return nil, ctx, err
	}

	sqsScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sqs.sh"))
	if err != nil {
		return nil, ctx, err
	}

	smScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sm.sh"))
	if err != nil {
		return nil, ctx, err
	}

	snsScriptReader, err := os.Open(snsScript)
	if err != nil {
		return nil, ctx, err
	}

	sqsScriptReader, err := os.Open(sqsScript)
	if err != nil {
		return nil, ctx, err
	}

	smScriptReader, err := os.Open(smScript)
	if err != nil {
		return nil, ctx, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "localstack/localstack:latest",
			ExposedPorts: []string{
				"4566",
			},
			Env: map[string]string{
				"SERVICES":       "secretsmanager,sqs,sns",
				"DEFAULT_REGION": "us-east-1",
				"DOCKER_HOST":    "unix:///var/run/docker.sock",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            snsScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sns.sh",
					FileMode:          0777,
				},
				{
					Reader:            sqsScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sqs.sh",
					FileMode:          0777,
				},
				{
					Reader:            smScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sm.sh",
					FileMode:          0777,
				},
			},
			WaitingFor: wait.ForListeningPort("4566/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})

	if err != nil {
		return nil, ctx, err
	}

	return container, ctx, nil
}

func createPostgresContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	dbScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-db.sql"))
	if err != nil {
		return nil, ctx, err
	}

	dbScriptReader, err := os.Open(dbScript)
	if err != nil {
		return nil, ctx, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16.0",
			ExposedPorts: []string{
				"5432",
			},
			Env: map[string]string{
				"POSTGRES_DB":       "order_db",
				"POSTGRES_USER":     "order",
				"POSTGRES_PASSWORD": "order",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            dbScriptReader,
					ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
					FileMode:          0644,
				},
			},
			WaitingFor: wait.ForLog("PostgreSQL init process complete; ready for start up").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to start postgres container: %w", err)
	}

	postgresIp, err := container.Host(ctx)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to get postgres ip: %w", err)
	}

	postgresPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to get postgres port: %w", err)
	}

	connStr := fmt.Sprintf("postgres://order:order@%s:%s/order_db?sslmode=disable", postgresIp, postgresPort.Port())

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, ctx, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return container, ctx, nil
}

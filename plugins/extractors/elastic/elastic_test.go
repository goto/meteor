//go:build plugins
// +build plugins

package elastic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/elastic"
	"github.com/goto/meteor/test/mocks"
	"github.com/goto/meteor/test/utils"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	host     = "http://localhost:9200"
	pass     = "secret_pass"
	user     = "elastic_meteor"
	urnScope = "test-elasticsearch"
)

var (
	client *elasticsearch.Client
	ctx    = context.TODO()
)

func TestMain(m *testing.M) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
		Username: user,
		Password: pass,
	}
	var err error
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// setup test
	opts := dockertest.RunOptions{
		Repository: "elasticsearch",
		Tag:        "7.13.2",
		Env: []string{
			"discovery.type=single-node",
			"ES_JAVA_OPTS=-Xms512m -Xmx512m",
			"ELASTIC_USER=" + user,
			"ELASTIC_PASSWORD=" + pass,
			"xpack.security.enabled=false",
		},
		ExposedPorts: []string{"9200"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9200": {
				{HostIP: "0.0.0.0", HostPort: "9200"},
			},
		},
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	retryFn := func(resource *dockertest.Resource) (err error) {
		res, err := client.Cluster.Health(
			client.Cluster.Health.WithLevel("indices"),
		)
		if err != nil {
			return
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("received %d status code", res.StatusCode)
		}
		return
	}
	purgeFn, err := utils.CreateContainer(opts, retryFn)
	if err != nil {
		log.Fatal(err)
	}
	if err := setup(); err != nil {
		log.Fatal(err)
	}

	// run tests
	code := m.Run()

	// clean tests
	if err := purgeFn(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestInit(t *testing.T) {
	t.Run("should return error if no host in config", func(t *testing.T) {
		err := newExtractor().Init(ctx, plugins.Config{
			URNScope: urnScope,
			RawConfig: map[string]interface{}{
				"password": "pass",
			},
		})
		assert.ErrorAs(t, err, &plugins.InvalidConfigError{})
	})
}

func TestExtract(t *testing.T) {
	t.Run("should return mockdata we generated with service running on localhost", func(t *testing.T) {
		extr := newExtractor()
		err := extr.Init(ctx, plugins.Config{
			URNScope: urnScope,
			RawConfig: map[string]interface{}{
				"host":     host,
				"user":     user,
				"password": pass,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		emitter := mocks.NewEmitter()
		err = extr.Extract(ctx, emitter.Push)
		assert.NoError(t, err)
		utils.AssertEqualProtos(t, getExpectedVal(t), utils.SortedAssets(emitter.GetAllData()))
	})
}

type MeteorMockElasticDocs struct {
	SomeStr string
}

func setup() (err error) {
	doc1 := MeteorMockElasticDocs{}
	doc1.SomeStr = "Value1"

	doc2 := MeteorMockElasticDocs{}
	doc2.SomeStr = "Value2"
	docStr1 := jsonStruct(doc1)
	docStr2 := jsonStruct(doc2)
	err = populateElasticSearch("index1", "1", docStr1)
	if err != nil {
		return
	}
	err = populateElasticSearch("index2", "1", docStr2)
	if err != nil {
		return
	}
	return
}

func populateElasticSearch(index, id, data string) error {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	res.Body.Close()
	return err
}

func jsonStruct(doc MeteorMockElasticDocs) string {
	docStruct := &MeteorMockElasticDocs{
		SomeStr: doc.SomeStr,
	}
	b, err := json.Marshal(docStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(b)
}

func newExtractor() *elastic.Extractor {
	return elastic.New(utils.Logger)
}

func getExpectedVal(t *testing.T) []*v1beta2.Asset {
	table1, err := anypb.New(&v1beta2.Table{
		Columns: []*v1beta2.Column{
			{
				Name:     "SomeStr",
				DataType: "text",
			},
		},
		Profile: &v1beta2.TableProfile{
			TotalRows: 1,
		},
		Attributes: &structpb.Struct{},
	})
	if err != nil {
		t.Fatal("error creating Any struct for test: %w", err)
	}
	table2, err := anypb.New(&v1beta2.Table{
		Columns: []*v1beta2.Column{
			{
				Name:     "SomeStr",
				DataType: "text",
			},
		},
		Profile: &v1beta2.TableProfile{
			TotalRows: 1,
		},
		Attributes: &structpb.Struct{},
	})
	if err != nil {
		t.Fatal("error creating Any struct for test: %w", err)
	}
	return []*v1beta2.Asset{
		{
			Urn:     "urn:elasticsearch:test-elasticsearch:index:index1",
			Name:    "index1",
			Type:    "table",
			Data:    table1,
			Service: "elasticsearch",
		},
		{
			Urn:     "urn:elasticsearch:test-elasticsearch:index:index2",
			Name:    "index2",
			Type:    "table",
			Data:    table2,
			Service: "elasticsearch",
		},
	}
}

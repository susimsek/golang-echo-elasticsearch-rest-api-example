package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"golang-echo-elasticsearch-rest-api-example/repository"
	"log"
)

func ElasticsearchConnection() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL(ElasticsearchUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	if err != nil {
		log.Fatal(err)
	}

	p, code, err := client.Ping(ElasticsearchUrl).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v , %v, %v, %v,  %v", p.ClusterName, p.Name, p.TagLine, p.Version.Number, code)

	fmt.Println("Connected to Elasticsearch!")

	CreateIndexIfDoesNotExist(context.Background(), client, repository.UserIndexName)

	return client, err
}

// CreateIndexIfDoesNotExist ...
func CreateIndexIfDoesNotExist(ctx context.Context, client *elastic.Client, indexName string) error {
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	res, err := client.CreateIndex(indexName).Do(ctx)

	if err != nil {
		return err
	}

	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	return nil
}

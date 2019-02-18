package env

import "gopkg.in/olivere/elastic.v5"

type ElasticSearch struct {
	Host        string
	IndexPrefix string
	Client      *elastic.Client
}

func GetElasticSearchClient() *elastic.Client {
	if Env.ElasticSearch.Client == nil {
		temp, err := elastic.NewClient(elastic.SetURL(Env.ElasticSearch.Host), elastic.SetSniff(false))
		if err == nil {
			Env.ElasticSearch.Client = temp
		} else {
			panic("can't open elasticsearch client:" + err.Error())
		}
	}
	return Env.ElasticSearch.Client
}

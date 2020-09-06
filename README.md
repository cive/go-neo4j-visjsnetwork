# go-neo4j-visjsnetwork

## what is this?

### en

[neo4j-visjs](https://github.com/neo4j-contrib/neovis.js/) project is already exist. But, that use user/pass in javascript (front-end). This is bad for me. So, i made this project which is making visjs data from neo4j object. I'm wating for pull requests🤗

### ja

neo4j-visjsプロジェクトがすでに存在します。ただし、JavaScript (フロントエンド) でuser / passを使用します。これは私にとっては悪いことです。それで、neo4jオブジェクトからvisjsデータを作るこのプロジェクトを作りました。プルリクエストお待ちしてます。🤗

## how to use

```
go get github.com/cive/go-neo4j-visjsnetwork
```

code:
```
func (conf Neo4jConfig) GetActedIn() *visjs.GraphObject {
	conn := conf.Connect(neo4j.AccessModeRead)
	defer conn.Close()
	result, err := conn.session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n)-[r:ACTED_IN]-(m) RETURN n, r, m",
			nil,
		)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		conf := visjs.NewConfig()
		data := visjs.NewGraphData()

		for result.Next() {
			record := result.Record()
			deco := visjs.Deco{}
			deco.LabelFunc = LabelFuncSample1
			deco.TitleFunc = TitleFuncSample1
			node1 := conf.Neo4jNode2VisjsNode(record.GetByIndex(0).(neo4j.Node), &deco)
			node2 := conf.Neo4jNode2VisjsNode(record.GetByIndex(2).(neo4j.Node), &deco)
			edge := conf.Neo4jEdge2VisjsEdge(record.GetByIndex(1).(neo4j.Relationship))
			data.AddNode(node1)
			data.AddNode(node2)
			data.AddEdge(edge)
		}

		return data.Export(), result.Err()
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return result.(*visjs.GraphObject)
}
```

## development

```
docker-compose build
docker-compose up
```

1. access to `localhost:7474` (neo4j browser)
2. make neo4j built-in example (movie & acted_in)
3. access to `localhost`

## please help me...

I use only for development `go get gin, cors`, but go mod contains `gin, cors`... Please let me know if there is a way not to include it.

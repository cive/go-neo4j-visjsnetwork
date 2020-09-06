package internal

import (
	"fmt"
	"log"

	"github.com/cive/go-neo4j-visjsnetwork/visjs"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func LabelFuncSample1(labels []string, props map[string]interface{}) string {
	if labels[0] == "Movie" {
		return fmt.Sprintf("%s", props["title"])
	} else if labels[0] == "Person" {
		return fmt.Sprintf("%s", props["name"])
	}
	return labels[0]
}

func TitleFuncSample1(labels []string, props map[string]interface{}) string {
	if labels[0] == "Movie" {
		if props["tagline"] != nil {
			return fmt.Sprintf("%s", props["tagline"])
		} else {
			return ""
		}
	} else if labels[0] == "Person" {
		if props["born"] != nil {
			return fmt.Sprintf("born in %d", props["born"])
		} else {
			return ""
		}
	}
	return fmt.Sprintf("%v", props)
}

func OptionNodeFunceSample1(labels []string, props map[string]interface{}, node *visjs.Node) {
	if labels[0] == "Movie" {
		node.Size = 10
	} else if labels[0] == "Person" {
		node.Size = 30
	}
}

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

package visjs

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type (
	Node struct {
		Borderwidth         float64 `json:"borderwidth,omitempty"`
		Borderwidthselected float64 `json:"borderwidthselected,omitempty"`
		Color               string  `json:"color,omitempty"`
		Opacity             float64 `json:"opacity,omitempty"`
		Group               string  `json:"group,omitempty"`
		Hidden              bool    `json:"hidden,omitempty"`
		Id                  uint64  `json:"id"`
		Label               string  `json:"label,omitempty"`
		Mass                float64 `json:"mass,omitempty"`
		Physics             bool    `json:"physics,omitempty"`
		Shape               string  `json:"shape,omitempty"`
		Size                float64 `json:"size,omitempty"`
		Title               string  `json:"title,omitempty"`
		Value               float64 `json:"value,omitempty"`
		Kind                string  `json:"kind"` // Original property
	}

	Edge struct {
		Arrows             string  `json:"arrows,omitemty"`
		Arrowstrikethrough bool    `json:"arrowstrikethrough,omitempty"`
		From               int64   `json:"from,omitempty"`
		Hidden             bool    `json:"hidden,omitempty"`
		Hoverwidth         float64 `json:"hoverwidth,omitempty"`
		Id                 int64   `json:"id,omitempty"`
		Label              string  `json:"label,omitempty"`
		Labelhighlightbold bool    `json:"labelhighlightbold,omitempty"`
		Length             float64 `json:"length,omitempty"`
		Physics            bool    `json:"physics,omitempty"`
		Selectionwidth     float64 `json:"selectionwidth,omitempty"`
		Selfreferencesize  float64 `json:"selfreferencesize,omitempty"`
		Title              string  `json:"title,omitempty"`
		To                 int64   `json:"to,omitempty"`
		Value              float64 `json:"value,omitempty"`
		Width              float64 `json:"width,omitempty"`
		Widthconstraint    float64 `json:"widthconstraint,omitempty"`
	}

	GraphData struct {
		NodeSet map[string]Node
		EdgeSet map[string]Edge
	}

	GraphObject struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	Config struct {
		ColorNumber int64
		Colors      []string
		NodeColors  map[string]string
	}

	Deco struct {
		LabelFunc
		TitleFunc
		OptionNodeFunc
	}

	// Labels, Props
	LabelFunc      func([]string, map[string]interface{}) string
	TitleFunc      func([]string, map[string]interface{}) string
	OptionNodeFunc func([]string, map[string]interface{}, *Node)
)

var color_list = []string{"#00bcd4", "#ff5722", "#b2ebf2", "#dd2c00"}

func (node Node) AddNodeSet(nodeSet map[string]Node) {
	nodeSet[fmt.Sprintf("%#v", node)] = node
}

func (edge Edge) AddEdgeSet(edgeSet map[string]Edge) {
	edgeSet[fmt.Sprintf("%#v", edge)] = edge
}

func (data GraphData) AddNode(node Node) {
	node.AddNodeSet(data.NodeSet)
}

func (data GraphData) AddEdge(edge Edge) {
	edge.AddEdgeSet(data.EdgeSet)
}

func (data GraphData) Export() *GraphObject {
	graph := GraphObject{}
	graph.Nodes = make([]Node, 0, 0)
	graph.Edges = make([]Edge, 0, 0)
	for _, node := range data.NodeSet {
		graph.Nodes = append(graph.Nodes, node)
	}
	for _, edge := range data.EdgeSet {
		graph.Edges = append(graph.Edges, edge)
	}
	return &graph
}

func NewConfig() *Config {
	config := Config{}
	config.SetColorScheme([]string{})
	config.NodeColors = make(map[string]string)
	return &config
}

func NewGraphData() *GraphData {
	data := GraphData{}
	data.NodeSet = make(map[string]Node)
	data.EdgeSet = make(map[string]Edge)
	return &data
}

func (conf *Config) SetColorScheme(input []string) {
	var default_colors = []string{"#00bcd4", "#ff5722", "#b2ebf2", "#dd2c00"}
	if len(input) == 0 {
		conf.Colors = default_colors
	} else {
		conf.Colors = input
	}
}

func (conf *Config) MakeColor() string {
	color_index := conf.ColorNumber
	conf.ColorNumber = (conf.ColorNumber + 1) % int64(len(conf.Colors))
	return conf.Colors[color_index]
}

func (conf Config) Neo4jNodes2VisjsNodes(origNodes []neo4j.Node) map[string]Node {
	nodeSet := make(map[string]Node)
	for _, origNode := range origNodes {
		node := conf.Neo4jNode2VisjsNode(origNode, nil)
		node.AddNodeSet(nodeSet)
	}
	return nodeSet
}

func (conf Config) Neo4jEdges2VisjsEdges(origEdges []neo4j.Relationship) map[string]Edge {
	edgeSet := make(map[string]Edge)
	for _, origEdge := range origEdges {
		edge := conf.Neo4jEdge2VisjsEdge(origEdge)
		edge.AddEdgeSet(edgeSet)
	}
	return edgeSet
}

func (conf *Config) Neo4jNode2VisjsNode(origNode neo4j.Node, deco *Deco) Node {
	node := Node{}
	node.Id = uint64(origNode.Id())
	node.Kind = origNode.Labels()[0]
	if deco.LabelFunc != nil {
		if label := deco.LabelFunc(origNode.Labels(), origNode.Props()); label != "" {
			node.Label = label
		}
	} else {
		for _, prop := range origNode.Props() {
			node.Label = fmt.Sprintf("%v", prop)
			break
		}
	}
	if deco.TitleFunc != nil {
		if title := deco.TitleFunc(origNode.Labels(), origNode.Props()); title != "" {
			node.Title = title
		}
	} else {
		node.Title = fmt.Sprintf("%v", origNode.Props())
	}
	if deco.OptionNodeFunc != nil {
		deco.OptionNodeFunc(origNode.Labels(), origNode.Props(), &node)
	}
	if origNode.Labels()[0] != "" {
		if val, ok := conf.NodeColors[origNode.Labels()[0]]; ok {
			node.Color = val
		} else {
			log.Printf("MakeColor to: %s", origNode.Labels()[0])
			node.Color = conf.MakeColor()
			log.Printf("color is: %s", node.Color)
			conf.NodeColors[origNode.Labels()[0]] = node.Color
		}
	}
	return node
}

func (conf Config) Neo4jEdge2VisjsEdge(origEdge neo4j.Relationship) Edge {
	edge := Edge{}
	edge.From = origEdge.StartId()
	edge.To = origEdge.EndId()
	edge.Label = origEdge.Type()
	return edge
}

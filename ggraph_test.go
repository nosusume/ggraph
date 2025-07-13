package ggraph_test

import (
	"testing"

	"github.com/nosusume/ggraph"
	"github.com/stretchr/testify/assert"
)

func TestGraphOperations(t *testing.T) {
	// 测试整数类型图

	intGraph := ggraph.NewGraph[int]()
	intGraph.AddNode(1)
	intGraph.AddNode(2)
	intGraph.AddEdge(1, 2)

	// 测试节点存在性
	assert.True(t, intGraph.HasNode(1), "节点1应存在")
	assert.True(t, intGraph.HasNode(2), "节点2应存在")
	assert.False(t, intGraph.HasNode(3), "节点3不应存在")

	// 测试边存在性
	assert.True(t, intGraph.HasEdge(1, 2), "边1->2应存在")
	assert.False(t, intGraph.HasEdge(2, 1), "边2->1不应存在")

	// 测试获取所有节点
	nodes := intGraph.Nodes()
	assert.ElementsMatch(t, []int{1, 2}, nodes, "节点列表应包含1和2")

	// 测试获取邻居
	neighbors := intGraph.Neighbors(1)
	assert.ElementsMatch(t, []int{2}, neighbors, "节点1的邻居应为2")

	// 测试字符串类型图
	strGraph := ggraph.NewGraph[string]()
	strGraph.AddNode("A")
	strGraph.AddNode("B")
	strGraph.AddEdge("A", "B")

	assert.True(t, strGraph.HasNode("A"), "字符串节点A应存在")
	assert.True(t, strGraph.HasEdge("A", "B"), "字符串边A->B应存在")
}

func TestAddNodeIdempotent(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	initialLen := len(graph.Nodes())
	graph.AddNode(1)
	graph.AddNode(1) // 重复添加
	assert.Equal(t, initialLen+1, len(graph.Nodes()), "重复添加节点不应增加数量")
}

func TestNeighborsEmpty(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	graph.AddNode(1)
	neighbors := graph.Neighbors(1)
	assert.Empty(t, neighbors, "无邻居时应返回空列表")
}
func TestEdgeCountAndNodeCount(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	graph.AddNode(1)
	graph.AddNode(2)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 1)
	assert.Equal(t, 2, graph.NodeCount(), "节点数量应为2")
	assert.Equal(t, 2, graph.EdgeCount(), "边数量应为2")
}

func TestStringRepresentation(t *testing.T) {
	graph := ggraph.NewGraph[string]()
	graph.AddEdge("A", "B")
	str := graph.String()
	assert.Contains(t, str, "A: [B]", "字符串表示应包含A: [B]")
	assert.Contains(t, str, "B: []", "字符串表示应包含B: []")
}

func TestToDTOAndNewGraphByDTO(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	graph.AddEdge(1, 2)
	dto := graph.ToDTO()
	newGraph := ggraph.NewGraphByDTO(dto)
	assert.True(t, newGraph.HasNode(1), "DTO恢复的图应包含节点1")
	assert.True(t, newGraph.HasEdge(1, 2), "DTO恢复的图应包含边1->2")
}

type testNode struct {
	val   int
	edges []ggraph.Edge[int]
}

func (n testNode) This() int {
	return n.val
}
func (n testNode) Edges() []ggraph.Edge[int] {
	return n.edges
}

func TestNewGraphByNodeList(t *testing.T) {
	nodes := []testNode{
		{val: 1, edges: []ggraph.Edge[int]{{From: 1, To: 2}}},
		{val: 2, edges: []ggraph.Edge[int]{}},
	}
	nodePtrs := make([]ggraph.Node[int], len(nodes))
	for i := range nodes {
		nodePtrs[i] = nodes[i]
	}
	graph := ggraph.NewGraphByNodeList(nodePtrs)
	assert.True(t, graph.HasEdge(1, 2), "应通过NodeList添加边1->2")
	assert.True(t, graph.HasNode(2), "应通过NodeList添加节点2")
}

func TestHasEdgeWithMissingNodes(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	assert.False(t, graph.HasEdge(1, 2), "不存在节点时应返回false")
}

func TestNeighborsNonExistentNode(t *testing.T) {
	graph := ggraph.NewGraph[int]()
	neighbors := graph.Neighbors(999)
	assert.Empty(t, neighbors, "不存在节点应返回空邻居列表")
}

package ggraph

import "slices"

// Graph 泛型邻接图结构，T为节点类型（需可比较）
type Graph[T comparable] struct {
	// 节点映射，用于快速查找节点索引
	nodes map[T]int
	// 邻接表，每个索引对应一个节点的邻居索引列表
	adj [][]int
}

// GraphDTO 用于序列化图结构的DTO（Data Transfer Object）
// 适用于JSON序列化，节点类型为interface{}以支持多种类型
type GraphDTO struct {
	// Nodes 存储图中的所有节点
	Nodes []interface{} `json:"nodes"`
	// Adjacency list，存储每个节点的邻居索引
	Adj [][]int `json:"adj"`
}

// Node 泛型节点接口，定义了从和到方法
// 适用于有向图的边表示
type Node[T comparable] interface {
	// This 返回当前节点的值
	This() T
	// From 返回边的起始节点
	From() T
	// To 返回边的终止节点
	To() T
}

func NewGraphByNodeList[T comparable](l []Node[T]) *Graph[T] {
	// 创建一个新的泛型图
	g := NewGraph[T]()
	// 遍历节点列表，添加节点和边
	for _, node := range l {
		g.AddNode(node.This())
		g.AddNode(node.From())
		g.AddNode(node.To())
		g.AddEdge(node.From(), node.To())
	}
	return g
}

// NewGraphByDTO 从GraphDTO创建一个新的泛型图
// 适用于从序列化数据恢复图结构
func NewGraphByDTO(dto *GraphDTO) *Graph[any] {
	// 创建一个新的泛型图
	g := NewGraph[interface{}]()
	// 添加所有节点
	for _, node := range dto.Nodes {
		g.AddNode(node)
	}
	// 添加所有边
	for i, neighbors := range dto.Adj {
		for _, neighborIndex := range neighbors {
			if neighborIndex < len(dto.Nodes) {
				g.AddEdge(dto.Nodes[i], dto.Nodes[neighborIndex])
			}
		}
	}
	return g
}

// NewGraph 初始化一个空的泛型邻接图
func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		nodes: make(map[T]int),
		adj:   make([][]int, 0),
	}
}

// AddNode 向图中添加一个节点（去重）
func (g *Graph[T]) AddNode(node T) {
	// 检查节点是否已存在
	if _, exists := g.nodes[node]; exists {
		return
	}
	// 分配新索引
	index := len(g.nodes)
	g.nodes[node] = index
	// 扩展邻接表（避免索引越界）
	if index >= cap(g.adj) {
		g.adj = append(g.adj, make([][]int, index+1)...)
	}
}

// AddEdge 添加一条从from到to的有向边（自动添加缺失节点）
func (g *Graph[T]) AddEdge(from, to T) {
	g.AddNode(from)
	g.AddNode(to)
	// 获取节点索引
	fromIndex := g.nodes[from]
	toIndex := g.nodes[to]
	// 添加有向边
	g.adj[fromIndex] = append(g.adj[fromIndex], toIndex)
}

// Nodes 返回图中所有节点的切片
func (g *Graph[T]) Nodes() []T {
	nodes := make([]T, 0, len(g.nodes))
	for node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// Neighbors 返回指定节点的所有邻居（邻接表直接映射）
func (g *Graph[T]) Neighbors(node T) []T {
	// 获取节点索引
	index := g.nodes[node]
	// 获取邻居索引列表
	neighborIndices := g.adj[index]
	// 映射邻居索引为节点值
	neighbors := make([]T, 0, len(neighborIndices))
	for _, neighborIndex := range neighborIndices {
		// 遍历邻居索引，获取节点值
		for node, idx := range g.nodes {
			if idx == neighborIndex {
				neighbors = append(neighbors, node)
				break
			}
		}
	}
	return neighbors
}

// HasNode 检查图中是否存在指定节点
func (g *Graph[T]) HasNode(node T) bool {
	// 检查节点索引是否存在
	if _, exists := g.nodes[node]; !exists {
		return false
	}
	_, exists := g.nodes[node]
	return exists
}

// HasEdge 检查是否存在从from到to的有向边
func (g *Graph[T]) HasEdge(from, to T) bool {
	if !g.HasNode(from) || !g.HasNode(to) {
		return false
	}
	// 获取节点索引
	fromIndex := g.nodes[from]
	toIndex := g.nodes[to]
	// 检查边是否存在
	return slices.Contains(g.adj[fromIndex], toIndex)
}

// ToDTO 将图转换为GraphDTO格式，适用于序列化
// 返回的DTO包含所有节点和邻接表
func (g *Graph[T]) ToDTO() *GraphDTO {
	nodes := make([]interface{}, 0, len(g.nodes))
	for node := range g.nodes {
		nodes = append(nodes, node)
	}

	return &GraphDTO{
		Nodes: nodes,
		Adj:   g.adj,
	}
}

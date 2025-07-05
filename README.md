# ggraph
[![Go Version](https://img.shields.io/github/go-mod/go-version/nosusume/ggraph)](https://github.com/nosusume/ggraph)
[![License](https://img.shields.io/github/license/nosusume/ggraph)](https://github.com/nosusume/ggraph/blob/main/LICENSE)

Generic adjacency list graph implementation in Go with serialization support.

## Installation
```bash
go get github.com/nosusume/ggraph
```

## Basic Usage
```go
package main

import (
	"github.com/nosusume/ggraph"
)

func main() {
	// Create a new graph with string nodes
	g := ggraph.NewGraph[string]()
	
	// Add nodes
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	
	// Add edges
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	
	// Query neighbors
	neighbors := g.Neighbors("A") // Returns ["B"]
	
	// Check edge existence
	hasEdge := g.HasEdge("B", "C") // Returns true
}
```

## Serialization
```go
// Convert graph to DTO
dto := g.ToDTO()

// Marshal to JSON
jsonData, err := json.Marshal(dto)
if err != nil {
	// handle error
}

// Create graph from DTO
reconstructedGraph := ggraph.NewGraphByDTO(dto)
```

## API Reference

### Graph[T comparable]
Generic adjacency list graph structure.

**Methods:**
- `func NewGraph[T comparable]() *Graph[T]`  
  Creates new empty graph
- `AddNode(node T)`  
  Adds node (deduplicated)
- `AddEdge(from, to T)`  
  Adds directed edge (auto-adds missing nodes)
- `Nodes() []T`  
  Returns all nodes
- `Neighbors(node T) []T`  
  Returns node's neighbors
- `HasNode(node T) bool`  
  Checks node existence
- `HasEdge(from, to T) bool`  
  Checks edge existence
- `ToDTO() *GraphDTO`  
  Converts to serializable DTO

### GraphDTO
Serializable graph representation.

**Fields:**
- `Nodes []interface{}` - Graph nodes
- `Adj [][]int` - Adjacency list

### Edge[T]
Edge representation interface.

**Methods:**
- `This() T` - Current node value
- `From() T` - Edge source node
- `To() T` - Edge target node

## License
MIT License - See [LICENSE](LICENSE) for details.

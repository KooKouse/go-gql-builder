package core

import (
	"github.com/Finovate/go-gql-builder/pkg/core/argument"
	"github.com/graphql-go/graphql"
)

// Node a Node represents a fundamental business object,
// serving as the core unit for constructing GraphQL queries.
// As an abstract concept, Node allows developers to define and declare their business data models.
// These models are then transformed by the framework into GraphQL query structures.
type Node interface {
	// Name query name
	// users{ <- this is the query name
	//   id
	//   name
	//   email
	// }
	Name() string
	// Type NodeType
	Type() FieldType

	Resolve() graphql.FieldResolveFn
	BuildFields() []*Field
	BuildArgs() []argument.Argument

	IsList() bool

	GetRegistry() *NodeRegistry
	SetRegistry(*NodeRegistry)
}

type BaseNode struct {
	registry *NodeRegistry
}

func (n *BaseNode) GetRegistry() *NodeRegistry {
	return n.registry
}

func (n *BaseNode) SetRegistry(registry *NodeRegistry) {
	n.registry = registry
}

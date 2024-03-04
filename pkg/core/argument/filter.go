package argument

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"strings"
)

var (
	filterArgumentType = graphql.NewScalar(graphql.ScalarConfig{
		Name:        FilterArgumentType,
		Description: "Filter argument",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {

			return parseAstValue(valueAST)
		},
	})
)

var _ SqlArgument = (*FilterArgument)(nil)

type FilterArgument struct {
	operationsMap map[string][]Operation
}

func newFilterArgument() *FilterArgument {
	return &FilterArgument{
		operationsMap: make(map[string][]Operation),
	}
}

func (f *FilterArgument) TypeName() string {
	return FilterArgumentType
}

func (f *FilterArgument) Validate(input interface{}) error {
	argsMap, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("filter argument must be a map[string]interface{}")
	}
	for fieldName, rawMap := range argsMap {
		operationMap, ok := rawMap.(map[string]interface{})
		if !ok {
			return fmt.Errorf("argument for field %s must be a map[string]interface{}", fieldName)
		}

		for op, value := range operationMap {
			operation, err := OperationFactory(op, fieldName, value)
			if err != nil {
				return err
			}
			if f.operationsMap[fieldName] == nil {
				f.operationsMap[fieldName] = make([]Operation, 0)
			}
			f.operationsMap[fieldName] = append(f.operationsMap[fieldName], operation)
		}

	}

	// TODO validate fieldName in Node fields or table columns
	return nil
}

func (f *FilterArgument) GetArgumentType() graphql.Input {
	return filterArgumentType
}

func (f *FilterArgument) ParseSqlValue() string {
	sqlStrings := make([]string, 0, len(f.operationsMap))
	for _, opList := range f.operationsMap {
		for _, operation := range opList {
			sqlStrings = append(sqlStrings, operation.ToSql())
		}
	}

	return strings.Join(sqlStrings, " AND ")
}

// 辅助函数：递归解析 AST 值
func parseAstValue(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.ObjectValue:
		// 创建一个映射来存储对象值
		result := make(map[string]interface{})
		for _, field := range valueAST.Fields {
			// 递归地解析字段值
			result[field.Name.Value] = parseAstValue(field.Value)
		}
		return result
	case *ast.ListValue:
		// 创建一个 slice 来存储列表中的元素
		list := make([]interface{}, len(valueAST.Values))
		for i, value := range valueAST.Values {
			// 递归解析列表中的每个元素
			list[i] = parseAstValue(value)
		}
		return list
	case *ast.StringValue:
		return valueAST.Value
	case *ast.BooleanValue:
		return valueAST.Value
	case *ast.IntValue:
		return valueAST.Value
	case *ast.FloatValue:
		return valueAST.Value
	case *ast.EnumValue:
		return valueAST.Value
	default:
		return nil
	}
}
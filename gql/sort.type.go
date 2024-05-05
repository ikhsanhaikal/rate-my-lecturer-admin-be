package gql

import "github.com/graphql-go/graphql"

func (factory *GqlFactory) SortType() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "Sort",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"order": &graphql.InputObjectFieldConfig{
				Type: graphql.NewEnum(graphql.EnumConfig{
					Name: "order",
					Values: graphql.EnumValueConfigMap{
						"ASC": &graphql.EnumValueConfig{
							Value: "ASC",
						},
						"DESC": &graphql.EnumValueConfig{
							Value: "DESC",
						},
					},
				}),
			},
		},
	})
}

package recipes

import (
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/budgets/types"
	"github.com/cloudquery/plugin-sdk/codegen"
	// "github.com/cloudquery/plugin-sdk/schema"
)

func BudgetsResources() []*Resource {
	resources := []*Resource{
		{
			SubService: "budgets",
			Struct:     &types.Budget{},
			SkipFields: []string{},
			ExtraColumns: append(
				defaultAccountColumns,
				[]codegen.ColumnDefinition{}...),
			Relations: []string{},
		},
	}

	// set default values
	for _, r := range resources {
		r.Service = "budgets"
		r.Multiplex = "client.AccountMultiplex"
		structName := reflect.ValueOf(r.Struct).Elem().Type().Name()
		if strings.Contains(structName, "Wrapper") {
			r.UnwrapEmbeddedStructs = true
		}
	}
	return resources
}

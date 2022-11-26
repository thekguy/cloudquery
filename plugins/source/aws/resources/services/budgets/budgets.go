// Code generated by codegen; DO NOT EDIT.

package budgets

import (
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func Budgets() *schema.Table {
	return &schema.Table{
		Name:      "aws_budgets_budgets",
		Resolver:  fetchBudgetsBudgets,
		Multiplex: client.AccountMultiplex,
		Columns: []schema.Column{
			{
				Name:     "account_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveAWSAccount,
			},
			{
				Name:     "budget_name",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("BudgetName"),
			},
			{
				Name:     "budget_type",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("BudgetType"),
			},
			{
				Name:     "time_unit",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("TimeUnit"),
			},
			{
				Name:     "auto_adjust_data",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("AutoAdjustData"),
			},
			{
				Name:     "budget_limit",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("BudgetLimit"),
			},
			{
				Name:     "calculated_spend",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("CalculatedSpend"),
			},
			{
				Name:     "cost_filters",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("CostFilters"),
			},
			{
				Name:     "cost_types",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("CostTypes"),
			},
			{
				Name:     "last_updated_time",
				Type:     schema.TypeTimestamp,
				Resolver: schema.PathResolver("LastUpdatedTime"),
			},
			{
				Name:     "planned_budget_limits",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("PlannedBudgetLimits"),
			},
			{
				Name:     "time_period",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("TimePeriod"),
			},
		},
	}
}

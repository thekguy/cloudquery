# Table: aws_budgets_budgets

https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_budgets_Budget.html

The primary key for this table is **cq_id**.


## Columns
| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|account_id (PK)|String|
|region|String|
|budget_name|String|
|budget_type|String|
|time_unit|String|
|auto_adjust_data|JSON|
|budget_limit|JSON|
|calculated_spend|JSON|
|cost_filters|JSON|
|cost_types|JSON|
|last_updated_time|Timestamp|
|planned_budget_limits|JSON|
|time_period|JSON|

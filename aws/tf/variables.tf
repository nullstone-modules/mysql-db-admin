variable "name" {
  description = "The name of the lambda function and role"
  type        = string
}

variable "tags" {
  description = "A map of tags that are applied to AWS resources"
  type        = map(string)
}

variable "host" {
  description = "The database cluster host to connect"
  type        = string
}

variable "port" {
  description = "The database cluster port to connect"
  type        = string
  default     = "3306"
}

variable "database" {
  description = "The initial database to connect"
  type        = string
  default     = ""
}

variable "username" {
  description = "Mysql username"
  type        = string
}

variable "password" {
  description = "Mysql password"
  type        = string
}

variable "is_prod_env" {
  type        = bool
  default     = true
  description = <<EOF
When destroying, is_prod_env determines the recovery window for the admin password secret.
If true, a 7-day recovery window will be configured.
If not, secret will be deleted immediately.
EOF
}

variable "alerts" {
  description = <<EOF
Configuration for CloudWatch alarms on the db-admin lambda function.
- enabled: Set to true to create the error-rate alarm (default: false)
- error_rate: Percentage of invocations that error over a 5-minute period to trigger the alarm (default: 5%)
- notification_arn: SNS topic ARN notified when an alarm changes state
EOF

  type = object({
    enabled          = optional(bool, false)
    error_rate       = optional(number, 5)
    notification_arn = optional(string, "")
  })

  default = {}
}

variable "network" {
  description = "Network configuration"
  type = object({
    vpc_id : string
    mysql_security_group_id : string
    security_group_ids : list(string)
    subnet_ids = list(string)
  })
  default = {
    vpc_id                  = ""
    mysql_security_group_id = ""
    security_group_ids      = []
    subnet_ids              = []
  }
}

resource "aws_secretsmanager_secret" "db_admin_mysql" {
  name                    = "${var.name}/conn_url"
  tags                    = var.tags
  recovery_window_in_days = var.is_prod_env ? 7 : 0
}

resource "aws_secretsmanager_secret_version" "db_admin_mysql" {
  secret_id     = aws_secretsmanager_secret.db_admin_mysql.id
  secret_string = "mysql://${urlencode(var.username)}:${urlencode(var.password)}@${var.host}:${var.port}/${urlencode(var.database)}"
}

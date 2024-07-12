
resource "relyt_dwuser" "user1" {
  dwsu_id          = "dwsu-id-from-an-duws-resource"
  account_name     = "demo5"
  account_password = "daf#$dgdfe&Abce%64"

  datalake_aws_lakeformation_role_arn      = "role arn arn=//xxxx"          # option
  async_query_result_location_prefix       = "s3=//bucket-name/abc/def/..." # option
  async_query_result_location_aws_role_arn = "role arn arn=//xxxx"          # option
}
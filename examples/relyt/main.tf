terraform {
  required_providers {
    relyt = {
      source = "hashicorp.com/edu/relyt"
    }
  }
}

provider "relyt" {
  #  至少是域名
  api_host = "http://120.92.213.241:8080"
  auth_key = "9a3727e5b9c0ddabaGbll2HVLVKLLY1AyjOilAqeyPOBAb74A7VlJRAdTi0bJWJd"
  #  可读角色，接口直接传入
  role = "343842875420708874"
}

locals {
  #  size_s = {
  #    id = 2
  #  }
  #  size_edps_s = {
  #    id = 52
  #  }
  cloud_id = {
    id = "ksc"
    #  "ksyun"
  }
  region_id = {
    id = "beijing-cicd"
  }
  BASIC = {
    id = "basic"
  }
}

#dps，类型用type来支持。尽可能平级，一级打平
resource "relyt_dwsu" "qingdeng" {
  cloud     = local.cloud_id.id
  region    = local.region_id.id
  dwsu_type = local.BASIC.id # 默认basic
  domain    = "qing-deng-tf"
  alias     = "qingdeng-test"
  default_dps = {
    name        = "hybrid"
    description = "qingdeng-test" #optional
    engine      = "hybrid"
    size        = "S"
  }
}

resource "relyt_dps" "abc" {
  dwsu_id     = relyt_dwsu.qingdeng.id
  name        = "edps-exp"
  description = "qingdeng-test" #optional
  engine      = "extreme"
  size        = "S"
}

resource "relyt_dwuser" "user1" {
  dwsu_id          = relyt_dwsu.qingdeng.id
  account_name     = "demo5"
  account_password = "daf#$dgdfe&Abce%64"

  #  datalake_aws_lakeformation_role_arn      = "aws's role arn arn=//xxxx"    # option
  datalake_aws_lakeformation_role_arn = "anotherRole2" # option
  async_query_result_location_prefix  = "simple"       # option
  #  async_query_result_location_prefix       = "s3=//bucket-name/abc/def/..." # option
  async_query_result_location_aws_role_arn = "anotherSimple" # option
  #  async_query_result_location_aws_role_arn = "aws's role arn arn=//xxxx"    # option
}


# 需要注意的细节点
#1、状态维护，以及中间某一步骤失败如何处理
#2、provider不同版本如何升级迭代，兼容老的数据？
#3、删除逻辑，tf框架如何保证先删除dps，再删除dwsu
#4、异步资源，如何维护状态。
#5、tf框架，变更的状态粒度。比如password变更，如何更新
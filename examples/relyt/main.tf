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
  role_id  = "343842875420708874"
}

variable "size_s" {
  id = 2
}
variable "cloud_id" {
  id = "ksc"
}
variable "region_id" {
  id = "beijing-cicd"
}
variable "BASIC" {
  id = "basic"
}

#dps，类型用type来支持。尽可能平级，一级打平
resource "relyt_dwsu" "qingdeng" {
  cloud: "ksyun"
  region: "beijing-cicd"
  dwsu_type: var.BASIC.id # 默认basic
  defaultDps : {
    name : "hybrid",
    description : "qingdeng-test", #optional
    engine : "hybrid",
    size: var.size_s.id
  }
}

resource "relyt_dps" "abc" {
  dwsu: relyt_dwsu.qingdeng
  name : "edps-exp",
  description : "qingdeng-test", #optional
  engine : "extream",
  size: var.size_s.id
}

resource "relyt_dwuser" "user1"{
  dwsu: relyt_dwsu.qingdeng
#  dwsu/
  account_name = "demo"
  account_password = "ssss"

#  account_type = "Super"
  datalake_aws_lakeformation_role_arn: "aws's role arn arn://xxxx" # option
  async_query_result_location_prefix: "s3://bucket-name/abc/def/..." # option
  async_query_result_location_aws_role_arn: "aws's role arn arn://xxxx" # option
}



# 需要注意的细节点
#1、状态维护，以及中间某一步骤失败如何处理
#2、provider不同版本如何升级迭代，兼容老的数据？
#3、删除逻辑，tf框架如何保证先删除dps，再删除dwsu
#4、异步资源，如何维护状态。
#5、tf框架，变更的状态粒度。比如password变更，如何更新
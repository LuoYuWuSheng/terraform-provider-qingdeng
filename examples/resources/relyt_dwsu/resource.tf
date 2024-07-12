

resource "relyt_dwsu" "dwsu" {
  cloud     = "ksc"
  region    = "beijing-cicd"
  dwsu_type = "basic"
  domain    = "qing-deng-tf"
  alias     = "qingdeng-test"
  default_dps = {
    name        = "hybrid"
    description = "qingdeng-test"
    engine      = "hybrid"
    size        = "S"
  }
}
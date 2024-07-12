terraform {
  required_providers {
    relyt = {
      source = "hashicorp.com/edu/relyt"
    }
  }
}

provider "relyt" {
  api_host = "http://localhost:19090"
  auth_key = "education"
  role_id  = "test123"
}

#data "hashicups_example" "edu" {}

data "relyt_coffees" "ccff" {}

output "edu_coffees" {
  value = data.relyt_coffees.ccff
}

resource "relyt_order" "order" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "edu_order" {
  value = relyt_order.order
}

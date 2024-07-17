# Terraform Provider Relyt 

```
example：

```

[//]: # ()
[//]: # (_This template repository is built on the [Terraform Plugin Framework]&#40;https://github.com/hashicorp/terraform-plugin-framework&#41;. The template repository built on the [Terraform Plugin SDK]&#40;https://github.com/hashicorp/terraform-plugin-sdk&#41; can be found at [terraform-provider-scaffolding]&#40;https://github.com/hashicorp/terraform-provider-scaffolding&#41;. See [Which SDK Should I Use?]&#40;https://developer.hashicorp.com/terraform/plugin/framework-benefits&#41; in the Terraform documentation for additional information._)

[//]: # ()
[//]: # (This repository is a *template* for a [Terraform]&#40;https://www.terraform.io&#41; provider. It is intended as a starting point for creating Terraform providers, containing:)

[//]: # ()
[//]: # (- A resource and a data source &#40;`internal/provider/`&#41;,)

[//]: # (- Examples &#40;`examples/`&#41; and generated documentation &#40;`docs/`&#41;,)

[//]: # (- Miscellaneous meta files.)

[//]: # ()
[//]: # (These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Developer]&#40;https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework&#41; platform. _Terraform Plugin Framework specific guides are titled accordingly._)

[//]: # ()
[//]: # (Please see the [GitHub template repository documentation]&#40;https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template&#41; for how to create a new repository from this template on GitHub.)

[//]: # ()
[//]: # (Once you've written your provider, you'll want to [publish it on the Terraform Registry]&#40;https://developer.hashicorp.com/terraform/registry/providers/publishing&#41; so that others can use it.)

[//]: # ()
[//]: # (## Requirements)

[//]: # ()
[//]: # (- [Terraform]&#40;https://developer.hashicorp.com/terraform/downloads&#41; >= 1.0)

[//]: # (- [Go]&#40;https://golang.org/doc/install&#41; >= 1.21)

[//]: # ()
[//]: # (## Building The Provider)

[//]: # ()
[//]: # (1. Clone the repository)

[//]: # (1. Enter the repository directory)

[//]: # (1. Build the provider using the Go `install` command:)

[//]: # ()
[//]: # (```shell)

[//]: # (go install)

[//]: # (```)

[//]: # ()
[//]: # (## Adding Dependencies)

[//]: # ()
[//]: # (This provider uses [Go modules]&#40;https://github.com/golang/go/wiki/Modules&#41;.)

[//]: # (Please see the Go documentation for the most up to date information about using Go modules.)

[//]: # ()
[//]: # (To add a new dependency `github.com/author/dependency` to your Terraform provider:)

[//]: # ()
[//]: # (```shell)

[//]: # (go get github.com/author/dependency)

[//]: # (go mod tidy)

[//]: # (```)

[//]: # ()
[//]: # (Then commit the changes to `go.mod` and `go.sum`.)

[//]: # ()
[//]: # (## Using the provider)

[//]: # ()
[//]: # (Fill this in for each provider)

[//]: # ()
[//]: # (## Developing the Provider)

[//]: # ()
[//]: # (If you wish to work on the provider, you'll first need [Go]&#40;http://www.golang.org&#41; installed on your machine &#40;see [Requirements]&#40;#requirements&#41; above&#41;.)

[//]: # ()
[//]: # (To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.)

[//]: # ()
[//]: # (To generate or update documentation, run `go generate`.)

[//]: # ()
[//]: # (In order to run the full suite of Acceptance tests, run `make testacc`.)

[//]: # ()
[//]: # (*Note:* Acceptance tests create real resources, and often cost money to run.)

[//]: # ()
[//]: # (```shell)

[//]: # (make testacc)

[//]: # (```)

stack "foobar" {
  source = "github.com/ConsciousML/terragrunt-template-stack//stacks/foo-bar?ref=v0.0.4"
  path   = "services"

  values = {
    output_dir = get_terragrunt_dir()
    content    = "Hello from foo in live repo"
  }
}
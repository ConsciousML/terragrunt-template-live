locals {
    environment_hcl = find_in_parent_folders("environment.hcl")
    environment = read_terragrunt_config(local.environment_hcl).locals.environment
}

stack "foobar" {
    source = "github.com/ConsciousML/terragrunt-template-stack//stacks/foo-bar"
    path = "services"

    values = {
        output_dir = get_terragrunt_dir()
        content = "Hello from foo in ${local.environment}!"
    }
}
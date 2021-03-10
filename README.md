# terraform-provider-wistia

A [Terraform](https://www.terraform.io/) provider for interacting with the Wistia API to create projects, media, etc.

## Why?

Maybe you want to create a lot of projects and populate them with a bunch of media. Maybe you want to do this repeatedly,
perhaps for QA purposes or for automated testing. Computers are really good at this type of thing, so let's have them do
the heavy lifting.

## Getting started

First, install Terraform:

```
brew install terraform
```

Then, look through the examples and create your own configuration in a file like `lenny_learns_terraform.tf`.

See the [Terraform docs](https://www.terraform.io/docs/language/index.html) for details about the language and how to
create resource configurations.

## Examples

The `examples/` directory has working configuration examples that you can use to create projects and media. You'll need 
to have an access token for your Wistia account (Account -> Settings -> API Access) that has read, update, delete, and 
upload permissions. Drop that token into the commands below, as indicated.

```
cd examples
terraform init
WISTIA_ACCESS_TOKEN="your access token goes here" terraform plan
```

The video files aren't included in this repository, so you'll need to edit the file paths to point at something on your
computer or upload by URL instead.
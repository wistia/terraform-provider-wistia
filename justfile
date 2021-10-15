install os_arch="darwin_amd64":
  go build -o ~/.terraform.d/plugins/registry.wistia.io/wistia/wistia/0.0.1/{{os_arch}}/terraform-provider-wistia_v0.0.1

clean:
  go clean

default_os_arch := `uname -s | tr 'A-Z' 'a-z'` + "_" + `uname -m`

install os_arch=default_os_arch:
  go build -o ~/.terraform.d/plugins/registry.wistia.io/wistia/wistia/0.0.1/{{os_arch}}/terraform-provider-wistia_v0.0.1

clean:
  go clean

name = waypoint
organization = hashicorp
version = 0.1.0
arch = darwin_amd64

build:
	go build -o bin/terraform-provider-$(name)_v$(version)

install: build
	mkdir -p ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)
	mv bin/terraform-provider-$(name)_v$(version) ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/$(arch)/
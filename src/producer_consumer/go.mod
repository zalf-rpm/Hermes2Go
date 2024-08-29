module github.com/zalf-rpm/Hermes2Go/src/producer_consumer

go 1.19

require (
	capnproto.org/go/capnp/v3 v3.0.1-alpha.2
	github.com/zalf-rpm/Hermes2Go/hermes v0.0.0-20210816164110-4329e56e99f8
	github.com/zalf-rpm/mas-infrastructure/capnproto_schemas/gen/go/test v0.0.0-20240726160202-0a44b9807771
)

replace github.com/zalf-rpm/Hermes2Go/hermes => ../../hermes

require (
	github.com/colega/zeropool v0.0.0-20230505084239-6fb4a4f75381 // indirect
	golang.org/x/sync v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

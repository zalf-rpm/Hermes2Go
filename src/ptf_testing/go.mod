module github.com/zalf-rpm/Hermes2Go/src/ptf_testing

go 1.17

require (
	github.com/go-echarts/go-echarts/v2 v2.2.4
	github.com/zalf-rpm/Hermes2Go/hermes v0.0.0-20221208175533-78d6bd21bdc0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

replace github.com/zalf-rpm/Hermes2Go/hermes => ../../hermes

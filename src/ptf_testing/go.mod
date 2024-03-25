module github.com/zalf-rpm/Hermes2Go/src/ptf_testing

go 1.19

require (
	github.com/go-echarts/go-echarts/v2 v2.3.3
	github.com/zalf-rpm/Hermes2Go/hermes v0.0.0-20240321153557-74f6c29be1ef
)

require (
	github.com/kr/text v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/zalf-rpm/Hermes2Go/hermes => ../../hermes

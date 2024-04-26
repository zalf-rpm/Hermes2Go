module github.com/zalf-rpm/Hermes2Go/src/renderservice

go 1.19

require (
	github.com/go-echarts/go-echarts/v2 v2.2.4
	github.com/zalf-rpm/Hermes2Go/hermes v0.0.0-20210825090813-daca9091f65f

)

replace github.com/zalf-rpm/Hermes2Go/hermes => ../../hermes

require gopkg.in/yaml.v3 v3.0.1 // indirect

module github.com/mats9693/study

go 1.17

require (
	github.com/pkg/errors v0.9.1
	github.com/mats9693/utils v0.0.0
)

replace (
	github.com/mats9693/utils => ../utils
)

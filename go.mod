module github.com/neurafuse/neuracli

go 1.16

replace github.com/neurafuse/neuracli => /Users/djw/CloudSync/AI/neurafuse/git/neuracli@v0.1.1

replace github.com/neurafuse/neurakube => /Users/djw/CloudSync/AI/neurafuse/git/neurakube@v0.1.1

replace github.com/neurafuse/tools-go => /Users/djw/CloudSync/AI/neurafuse/git/tools-go@v0.1.1

replace github.com/containers/podman/v2 => github.com/containers/libpod/v2 v2.2.1

require (
	github.com/c-bata/go-prompt v0.2.6
	github.com/neurafuse/neurakube v0.0.0-00010101000000-000000000000
	github.com/neurafuse/tools-go v0.1.1
	github.com/spf13/cobra v1.2.1
)

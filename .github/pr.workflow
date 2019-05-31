workflow "Build and Test" {
  resolves = ["Test"]
  on = "pull_request"
}

action "Build" {
  uses = "actions-contrib/go@master"
  args = "build ./..."
}

action "Test" {
  needs = "Build"
  uses = "actions-contrib/go@master"
  args = "test -v -short -race ./..."
}

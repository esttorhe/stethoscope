workflow "Build and Test" {
  on = "push"
  resolves = ["Test"]
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

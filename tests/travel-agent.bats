#!/usr/bin/env bats

source $BATS_TEST_DIRNAME/../bin/common.sh
load assertions

travel_example_dir=/tmp/travel-agent-example

setup(){
  rm -rf $travel_example_dir
}

teardown(){
  rm -rf $travel_example_dir
}

@test "It should book when it receives book argument" {
  run $BATS_TEST_DIRNAME/../bin/travel-agent book
  assert_line 0 "Booking..." 

}

@test "Target should fail when concourse target input is missing" {
  run target ""
  assert_line 0 "${red}You must provide a concourse target${reset}" 
  assert_failure
}


@test "It books when project project_path is a travel agent project" {
  DRY_RUN=true run $BATS_TEST_DIRNAME/../bin/travel-agent book $BATS_TEST_DIRNAME/../manifest/assets/travel-agent.yml
  assert_match "Travel agent project found"
  assert_not_match "This does not look like a travel agent project"
  assert_success
}


@test "It does not book when project folder does not look like a travel agent project" {
  run $BATS_TEST_DIRNAME/../bin/travel-agent book
  assert_match "This does not look like a travel agent project"
  assert_not_match "Travel agent project found"
  assert_failure
}

@test "When booking it removes previous /tmp/project-name if exists" {
  mkdir -p $travel_example_dir/banana
  DRY_RUN=true run $BATS_TEST_DIRNAME/../bin/travel-agent book $BATS_TEST_DIRNAME/../manifest/assets/travel-agent.yml
  assert_success
}

@test "When git_project is a git url it performs a clone" {
  run clone_project "https://github.com/compozed/travel-agent-example.git"
  assert_success
  assert_dir_exists $travel_example_dir/.git
}

@test "When git_project is a path it skips git clone" {
  mkdir -p /tmp/travel-agent-example/ci/manifest
  run clone_project "/tmp/travel-agent-example"
  assert_success
  assert_dir_not_exists $travel_example_dir/.git
}

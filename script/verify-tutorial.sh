#!/usr/bin/env bash

set -o errexit
# set -o xtrace # for debug uncomment

projectDir=$(cd "$(dirname "${0}")/.." && pwd)
# shellcheck source=script/util.sh
source "${projectDir}/script/util.sh" || source ./util.sh

verify () {
  tutorial="${1}"
  branches="${2}"

  println "VERIFYING ${tutorial}..."

  tutorialPath="tmp/${tutorial}"
  if [ ! -d "${tutorialPath}" ]; then
    git clone "https://github.com/EngoEngine/${tutorial}.git" "${tutorialPath}"
  fi
  cd "${tutorialPath}"

  for branch in $branches
  do
      println "VERIFYING ${branch}..."
      git checkout "${branch}"
      rm -f go.mod
      rm -f go.sum
      go mod init github.com/EngoEngine/TrafficManager
      if [ "${OS_FAMILY}" == "windows" ]; then
        go mod edit -replace="github.com/EngoEngine/engo=D:$(printf "%s" "${projectDir:2}" | tr / \\)"
      else
        go mod edit -replace="github.com/EngoEngine/engo=${projectDir}"
      fi
      "${projectDir}/script/go-build.sh"
      go clean
  done
}

verify "TrafficManager" "01-hello-world 02-first-system 03-camera-movement 04-hud 05-tilemaps 06-spritesheets-and-automated-city-building 07-hud-text"

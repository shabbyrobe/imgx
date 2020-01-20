#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
script_abspath="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mods=(
    "rgba"
    "termpalette"
    "testimg"
)

cmd-test() {
    for mod in "${mods[@]}"; do
        pushd "$script_abspath/$mod"
            go test
        popd > /dev/null
    done
}

"cmd-$1" "${@:2}"

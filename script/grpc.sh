#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "grpc protoc helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -l, --lang          generate grpc code for language (default \"go\") [go]"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o hl: --long help,lang: -n "$Command" -- "$@"`
eval set -- "${ARGS}"
lang="go"
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -l|--lang)
            lang="$2"
            shift 2
        ;;
        --)
            shift
            break
        ;;
        *)
            echo Error: unknown flag "$1" for "$Command"
            echo "Run '$Command --help' for usage."
            exit 1
        ;;
    esac
done

function buildGo(){
    cd "$Dir"
    local output=protocol
    if [[ -d "$output" ]];then
        rm "$output/*" -rf
    else
        mkdir "$output"
    fi
    local document="static/document/api"
    if [[ -d "$document" ]];then
        rm "$document/*" -rf
    else
        mkdir "$document"
    fi

    local command=(
        protoc -I "pb" -I "third_party/googleapis" \
    )
    # generate grpc code
    local opts=(
        --go_out="'$output'" --go_opt=paths=source_relative
        --go-grpc_out="'$output'" --go-grpc_opt=paths=source_relative
        --grpc-gateway_out="'$output'" --grpc-gateway_opt=paths=source_relative
        --openapiv2_out="'$document'" --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true --openapiv2_opt allow_merge=true
    )
    local exec="${command[@]} ${opts[@]} ${Protos[@]}"
    echo $exec
    eval "$exec"
}
case "$lang" in
    go)
        buildGo
    ;;
    *)
        echo Error: unknown language "$1" for "$Command"
        echo "Run '$Command --help' for usage."
        exit 1
    ;;
esac
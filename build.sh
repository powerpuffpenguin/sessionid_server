#!/usr/bin/env bash
Target=server
function help(){
    echo "build script"
    echo
    echo "Usage:"
    echo "  $0 [flags]"
    echo "  $0 [command]"
    echo
    echo "Available Commands:"
    echo "  help              help for $0"
    echo "  go                go build helper"
    echo
    echo "Flags:"
    echo "  -h, --help          help for $0"
}
function helpGo(){
    echo "go build helper"
    echo
    echo "Usage:"
    echo "  $0 go [flags]"
    echo
    echo "Flags:"
    echo "      --arch          GOARCH (default \"$(go env GOHOSTARCH)\")"
    echo "      --os            GOOS (default \"$(go env GOHOSTOS)\")"
    echo "  -c, --clear         clear output"
    echo "  -d, --debug         build debug mode"
    echo "  -l, --list          list all supported platforms"
    echo "  -p,--pack           pack to compressed package [7z gz bz2 xz zip]"
    echo "  -h, --help          help for $0"
}
function buildGo(){
    local command="$0 go"
    ARGS=`getopt -o hldp:c --long help,list,os:,arch:,debug,pack:,clear -n "$command" -- "$@"`
    if [ $? != 0 ]; then
        exit 1
    fi
    eval set -- "${ARGS}"
    local list=0
    local debug=0
    local clear=0
    local os="$(go env GOHOSTOS)"
    local arch="$(go env GOHOSTARCH)"
    local pack
    while true
    do
        case "$1" in
            -h|--help)
                helpGo
                exit 0
            ;;
            -l|--list)
                list=1
                shift
            ;;
            -d|--debug)
                debug=1
                shift
            ;;
            -c|--clear)
                clear=1
                shift
            ;;
            --os)
                os="$2"
                shift 2
            ;;
            --arch)
                arch="$2"
                shift 2
            ;;
            -p|--pack)
                case "$2" in
                    7z|gz|bz2|xz|zip)
                        pack="$2"
                    ;;
                    *)
                        echo Error: unknown pack "$2" for "$command --pack"
                        echo "Supported: 7z gz bz2 xz zip"
                        exit 1
                    ;;
                esac
                shift 2
            ;;
            --)
                shift
                break
            ;;
            *)
                echo Error: unknown flag "$1" for "$command"
                echo "Run '$command --help' for usage."
                exit 1
            ;;
        esac
    done

    if [[ "$list" == 1 ]];then
        go tool dist list
        exit $?
    fi
    if [[ "$clear" == 1 ]];then
        cd "$(dirname $BASH_SOURCE)"/bin
        rm -f *.gz *.bz *.bz2 *.xz *.7z *.zip *.db \
            server serverd
        exit $?
    fi
    
    export GOOS="$os"
    export GOARCH="$arch"
    local target="$Target"
    if [[ "$debug" == 1 ]];then
        target="$target"d
    fi
    if [[ "$os" == "windows" ]];then
        target="$target.exe"
    fi
    if [[ "$debug" == 1 ]];then
        local args=(
            go build 
            -o "bin/$target"
        )
    else
        local args=(
            go build 
            -ldflags '"-s -w"'
            -o "bin/$target"
        )
    fi

    set -e
    
    # build
    cd $(dirname $BASH_SOURCE)
    echo "build for \"$GOOS/$GOARCH\""
    local exec="${args[@]}"
    echo $exec
    eval "$exec"

    # pack
    args=""
    if [[ "$debug" == 1 ]];then
        local name="${Target}d_${GOOS}_$GOARCH"
    else
        local name="${Target}_${GOOS}_$GOARCH"
    fi
    case "$pack" in
        7z)
            name="$name.7z"
            args=(7z a "$name")
        ;;
        zip)
            name="$name.zip"
            args=(zip -r "$name")
        ;;
        gz)
            name="$name.tar.gz"
            args=(tar -zcvf "$name")
        ;;
        bz2)
            name="$name.tar.gz"
            args=(tar -jcvf "$name")
        ;;
        xz)
            name="$name.tar.gz"
            args=(tar -Jcvf "$name")
        ;;
    esac
    if [[ "$args" == "" ]];then
        return 0
    fi
    cd bin
    if [[ -f "$name" ]];then
        rm "$name"
    fi
    local source=(
        "$target"
        server.jsonnet cnf
    )
    local exec="${args[@]} ${source[@]} "
    echo $exec
    eval "$exec"
}

case "$1" in
    help|-h|--help)
        help
    ;;
    go)
        shift
        buildGo "$@"
    ;;
    *)
        if [[ "$1" == "" ]];then
            help
        elif [[ "$1" == -* ]];then
            echo Error: unknown flag "$1" for "$0"
            echo "Run '$0 --help' for usage."
        else
            echo Error: unknown command "$1" for "$0"
            echo "Run '$0 --help' for usage."
        fi        
        exit 1
    ;;
esac
Target="server"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v1.0.0"
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/mips
    linux/mips64
    linux/ppc64le
    linux/s390x
    linux/amd64
)
Protos=(
    manager/manager.proto
)
#!/bin/bash

set -e
function msg { out "$*" >&2 ;}
function err { local x=$? ; msg "$*" ; return $(( $x == 0 ? 1 : $x )) ;}
function out { printf '%s\n' "$*" ;}

function globals {
    this="$(cd "$(dirname "$0")" && pwd -P)"
    OutputPath="$this/../output/"
    ProgramName="crond"
    ConfigPath="$this/../config/config.yaml"
    TargetPlatform="Linux"
    DoFunction=""
    dryRun=0
    DebugMode=0
}; globals

function parseCmds {
  while [[ $# -gt 0 ]]
  do
    case "$1" in                                      # Munging globals, beware
      -t)                           TargetPlatform="$2"         ; shift 2 ;;
      -o)                           OutputPath="$2"             ; shift 2 ;;
      -e)                           DoFunction="$2"             ; shift 2 ;;
      -d)                           DebugMode=1                 ; shift 1;;
      --dry-run)                    dryRun=true                 ; shift 1 ;;
      *)                            err 'Argument error. Please see help.' ;;
    esac
  done
  : "${OutputPath:=$this/../output/}"
}

do_prepare() {
    if [ -d "$OutputPath" ]; then
        rm -r "$OutputPath"
    fi
    mkdir -p $OutputPath
    # ./gen_proto.sh
}

do_build() {
    if [ "$TargetPlatform" = "Linux" ];then
        export CGO_ENABLED=0
        export GOOS=linux
        export GOARCH=amd64
    elif [ "$TargetPlatform" = "Linux_x86" ];then
        export CGO_ENABLED=0
        export GOOS=linux
        export GOARCH=386
        export CC=x86_64-linux-musl-gcc
        export CXX=x86_64-linux-musl-g++
    fi
    go build -o "$OutputPath/$ProgramName" "$this/../cmd/crontab/crond.go"
}

do_pack() {
    cd "$OutputPath"
    tar -czvf "$ProgramName.tar.gz" *
    cd - &>/dev/null
}

do_postAction() {
    cp "$ConfigPath" "$OutputPath"
    do_pack
}

main() {
    if [ $DebugMode -ne 0 ];then
        set -xe
    fi
    export PATH=$PATH:/usr/local/go/bin:$this
    do_prepare
    do_build
    # do_postAction
}

parseCmds "$@"

# Test single function
if [ ! -z "$DoFunction" ]; then
	do_$DoFunction "$@"
else
	main
fi

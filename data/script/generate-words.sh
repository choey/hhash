#!/usr/bin/env bash

print_usage() {
  echo "Usage: $0 [options] dict_file"
  echo "  -m (int) minimum word length"
  echo "  -M (int) maximum world length"
  echo "  -f (file) word dictionary (required)"
  echo "  -h help display"
  echo "(example: ./generate-words.sh -m4 -M9 ../data/data.adj | shuf | tail -n 1000)"
}

MIN_LEN=3
MAX_LEN=8

while getopts 'm:M:f:' flag; do
  case "${flag}" in
    m) MIN_LEN="${OPTARG}" ;;
    M) MAX_LEN="${OPTARG}" ;;
    f) DICT="${OPTARG}" ;;
    *) print_usage
       exit 1 ;;
  esac
done

if [[ "$DICT" = "" ]]; then
    echo "$0: the first positional argument must be the word dictionary" 1>&2
    print_usage
    exit 1
elif [[ ! -f "$DICT" ]]; then
    echo "$0: dictionary not found at $DICT" 1>&2
    print_usage
    exit 1
fi

cut -d ' ' -f 5 $DICT | egrep "^[a-zA-Z]{${MIN_LEN},${MAX_LEN}}$" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2)) }' | uniq
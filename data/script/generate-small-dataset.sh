#!/usr/bin/env bash

print_usage() {
  echo "Usage: $0 options"
  echo "  -d (file) wordnet data directory"
  echo "  -t (file) target directory in which to create dataset files"
  echo "(example: ./generate-default-datasets.sh -d ../data -t ../datasets.go.template)"
}


while getopts 'd:t:' flag; do
  case "${flag}" in
    d) DATA_DIR="${OPTARG}" ;;
    t) TARGET_DIR="${OPTARG}" ;;
    *) print_usage
       exit 1 ;;
  esac
done

if [[ ! -d "$DATA_DIR" ]]; then
    echo "$0: -d must point to wordnet data directory containing data.adj, data.noun, ..."
    exit 1
fi

if [[ ! -d "$TARGET_DIR" ]]; then
    echo "$0: invalid target directory: $TARGET_DIR"
    exit 1
fi

GENERATE_WORDS_SH="$(dirname $0)/generate-words.sh"
${GENERATE_WORDS_SH} -m 3 -M 6 -f ${DATA_DIR}/data.adj | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/adj-sample.csv
${GENERATE_WORDS_SH} -m 3 -M 9 -f ${DATA_DIR}/data.adv | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/adv-sample.csv
${GENERATE_WORDS_SH} -m 3 -M 9 -f ${DATA_DIR}/data.verb | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/verb-sample.csv
${GENERATE_WORDS_SH} -m 3 -M 9 -f ${DATA_DIR}/data.noun | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/noun-sample.csv

MIN_LEN=3
MAX_LEN=8
egrep "^[a-zA-Z]{$MIN_LEN,$MAX_LEN} [a-zA-Z]{$MIN_LEN,$MAX_LEN}" ${DATA_DIR}/verb.exc | cut -d ' ' -f1 | egrep -v ".*ing$" | egrep -v ".*s$" | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/verb-past-sample.csv
egrep "^[a-zA-Z]{$MIN_LEN,$MAX_LEN} [a-zA-Z]{$MIN_LEN,$MAX_LEN}" ${DATA_DIR}/verb.exc | cut -d ' ' -f1 | egrep ".*ing$" | sort | uniq -i | shuf -n 5000 > ${TARGET_DIR}/verb-gerund-sample.csv
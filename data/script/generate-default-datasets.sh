#!/usr/bin/env bash

print_usage() {
  echo "Usage: $0 options"
  echo "  -d (file) dataset directory"
  echo "  -t (file) path to datasets.go.template"
  echo "(example: ./generate-default-datasets.sh -d ../data -t ../datasets.go.template)"
}

while getopts 'd:t:' flag; do
  case "${flag}" in
    d) DATA_DIR="${OPTARG}" ;;
    t) DATASETS_GO_TEMPLATE="${OPTARG}" ;;
    *) print_usage
       exit 1 ;;
  esac
done

if [[ ! -d "$DATA_DIR" ]]; then
    echo "$0: -d must point to datasets directory containing *-sample.csv"
    exit 1
fi

if [[ ! -f "$DATASETS_GO_TEMPLATE" ]]; then
    echo "$0: -f must point to datasets.go.template"
    exit 1
fi

QUOTE_REGEX="s/^\(.*\)$/\"\1\"/g"
DEFAULT_ADJECTIVES=$(sed ${QUOTE_REGEX} ${DATA_DIR}/adj-sample.csv | paste -sd "," -)
DEFAULT_ADVERBS=$(sed ${QUOTE_REGEX} ${DATA_DIR}/adv-sample.csv | paste -sd "," -)
DEFAULT_NOUNS=$(sed ${QUOTE_REGEX} ${DATA_DIR}/noun-sample.csv | paste -sd "," -)
DEFAULT_VERBS=$(sed ${QUOTE_REGEX} ${DATA_DIR}/verb-sample.csv | paste -sd "," -)
DEFAULT_VERBS_GERUND=$(sed ${QUOTE_REGEX} ${DATA_DIR}/verb-gerund-sample.csv | paste -sd "," -)
DEFAULT_VERBS_PAST=$(sed ${QUOTE_REGEX} ${DATA_DIR}/verb-past-sample.csv | paste -sd "," -)

DATASETS_GO="$(dirname ${DATASETS_GO_TEMPLATE})/$(filename -r ${DATASETS_GO_TEMPLATE})"
cp ${DATASETS_GO_TEMPLATE} ${DATASETS_GO}

sed -i "s/DEFAULT_ADJECTIVES/$DEFAULT_ADJECTIVES/" ${DATASETS_GO}
sed -i "s/DEFAULT_ADVERBS/$DEFAULT_ADVERBS/" ${DATASETS_GO}
sed -i "s/DEFAULT_NOUNS/$DEFAULT_NOUNS/" ${DATASETS_GO}
sed -i "s/DEFAULT_VERBS_GERUND/$DEFAULT_VERBS_GERUND/" ${DATASETS_GO}
sed -i "s/DEFAULT_VERBS_PAST/$DEFAULT_VERBS_PAST/" ${DATASETS_GO}
sed -i "s/DEFAULT_VERBS/$DEFAULT_VERBS/" ${DATASETS_GO}
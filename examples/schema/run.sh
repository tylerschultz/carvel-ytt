#!/bin/bash

set -e

./ytt -f examples/schema/config.yml -f examples/schema/schema.yml \
 --enable-experiment-schema

echo '***'

./ytt -f examples/schema/config.yml -f examples/schema/schema.yml \
  --data-value-yaml nothing="a new string" \
  --data-value-yaml string=str \
  --data-value-yaml bool=true \
  --data-value-yaml int=123 \
  --data-value-yaml float=123.123 \
  --data-value-yaml any=[1,2,4] \
  --enable-experiment-schema

echo '***'

export YAML_VAL_nothing="a new string"
export YAML_VAL_string=str
export YAML_VAL_bool=true
export YAML_VAL_int=123
export YAML_VAL_float=123.123
export YAML_VAL_any=[1,2,4]

./ytt -f examples/schema/config.yml -f examples/schema/schema.yml \
  --data-values-env-yaml YAML_VAL \
  --enable-experiment-schema

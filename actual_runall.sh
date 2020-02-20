rm -f input/*.OUT

find ./input -type f -name '*.txt' -exec runall.sh {} ';'

#!/bin/sh

find . -mindepth 2 -maxdepth 2 -type d | \
  grep -E './*/[0-9]{4}-[0-9]{2}-[0-9]{2}_*' | \
  awk -F'/' '{
    if (match($3, /^[0-9]{4}-[0-9]{2}-[0-9]{2}/, arr)) {
      print arr[0] "\t" $0
    }
  }' > /tmp/index_dirs_with_date.tmp

sort -r /tmp/index_dirs_with_date.tmp | cut -f2- > /tmp/index_dirs.tmp
count=$(wc -l < /tmp/index_dirs.tmp)

printf '# Today I Learned (TIL)\n\n'
printf '## Index (%s entries, newest first)\n\n' "$count"

while IFS= read -r entry; do
  entry_name=${entry#./}
  printf -- '- [%s](./%s)\n' "$entry_name" "$entry_name"
done < /tmp/index_dirs.tmp

rm -f /tmp/index_dirs_with_date.tmp /tmp/index_dirs.tmp


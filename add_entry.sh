#!/bin/sh

valid_re='^[a-zA-Z0-9-]+$'

read -p "Enter entry category (/$valid_re/): " category
echo "$category" | grep -qE $valid_re || { echo "Invalid category."; exit 1; }

read -p "Enter entry title (/$valid_re/): " title
echo "$title" | grep -qE $valid_re || { echo "Invalid title."; exit 1; }

today=$(date +%Y-%m-%d)
entry_dir="${category}/${today}_${title}"
mkdir -p "$entry_dir"
echo "Generated $entry_dir"

echo "# $title" > "$entry_dir/README.md"
echo "New entry created at $entry_dir/README.md"


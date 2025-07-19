#!/bin/sh
read -p "Enter entry category (alphanumeric/hyphen): " category
echo "$category" | grep -qE '^[a-zA-Z0-9-]+$' || { echo "Invalid category."; exit 1; }

read -p "Enter entry title (alphanumeric/hyphen): " title
echo "$title" | grep -qE '^[a-zA-Z0-9-]+$' || { echo "Invalid title."; exit 1; }

today=$(date +%Y-%m-%d)
entry_dir="${category}/${today}_${title}"
mkdir -p "$entry_dir"
echo "Generated $entry_dir"

echo "# $title" > "$entry_dir/README.md"
echo "New entry created at $entry_dir/README.md"


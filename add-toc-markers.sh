#!/bin/bash

# Add TOC markers to all markdown files that don't have them
find . -name "*.md" -not -path "./.git/*" -not -path "./node_modules/*" | while read file; do
  echo "Processing $file"
  
  # Check if file already has TOC markers
  if ! grep -q "<!-- toc -->" "$file"; then
    echo "  Adding TOC markers to $file"
    
    # Create a temporary file
    temp_file=$(mktemp)
    
    # Add TOC markers after the first line (title)
    {
      head -1 "$file"
      echo ""
      echo "<!-- toc -->"
      echo ""
      echo "<!-- tocstop -->"
      echo ""
      tail -n +2 "$file"
    } > "$temp_file"
    
    # Replace the original file
    mv "$temp_file" "$file"
  else
    echo "  $file already has TOC markers"
  fi
done
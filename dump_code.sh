#!/bin/bash

OUTPUT_FILE="codebase_dump.txt"

> "$OUTPUT_FILE"

EXCLUDE_PATTERNS=(
  ".git" "node_modules" "vendor" "dist" "build" "__pycache__" "__debug_bin*"
  ".venv" "venv" "env" "*.pyc" ".idea" ".vscode" "*.pem" "*.lock" "*.pb.go" "*.pb" "$OUTPUT_FILE"
)

EXCLUDE_DIR_PATTERNS=(
 "web/gen"
)

EXCLUDE_ARGS=()
for pattern in "${EXCLUDE_PATTERNS[@]}"; do
  EXCLUDE_ARGS+=(-name "$pattern" -prune -o)
done

EXCLUDE_DIR_ARGS=()
for pattern in "${EXCLUDE_DIR_PATTERNS[@]}"; do
  EXCLUDE_DIR_ARGS+=(-wholename "$pattern" -prune -o)
done

echo "Gathering files and generating $OUTPUT_FILE..."

SEARCH_DIR="${1:-.}"

find "$SEARCH_DIR" "${EXCLUDE_ARGS[@]}" "${EXCLUDE_DIR_ARGS[@]}"  -type f -print0 | while IFS= read -r -d '' file; do
  
  filepath=${file#./}
  
  echo -e "FILE: $filepath" >> "$OUTPUT_FILE"
  
  cat "$file" >> "$OUTPUT_FILE"
  
  echo "" >> "$OUTPUT_FILE"

done

echo "Done! Open $OUTPUT_FILE and paste its contents to the AI."
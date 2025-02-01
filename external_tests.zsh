#!/usr/bin/env zsh

# Check if tag argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

TAG=$1

# Function to handle errors
error_exit() {
    echo "$1" >&2
    exit 1
}

# Run tests
echo "🛠 Running tests..."
if ! go test ./... -v; then
    error_exit "❌ Tests failed. Exiting."
fi

# Commit changes
echo "💾 Committing changes..."
git add . || error_exit "❌ Git add failed"
git commit -m "Release $TAG" || error_exit "❌ Commit failed"

# Create and push tag
echo "🏷 Tagging with $TAG..."
git tag "$TAG" || error_exit "❌ Tagging failed"
git push origin HEAD || error_exit "❌ Push failed"
git push origin "$TAG" || error_exit "❌ Tag push failed"

# Change directory
echo "📂 Changing directory..."
cd ~/Desktop/Code/tiny-ai/tiny-ai-test || error_exit "❌ Directory change failed"

# Update dependency and run
echo "🔄 Updating dependency..."
go get github.com/matwate/sometinyai@latest || error_exit "❌ Dependency update failed"

echo "🚀 Running the code..."
go run . || error_exit "❌ Failed to run code"

echo "✅ All tasks completed successfully!"
d

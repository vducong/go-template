#!/bin/sh

# Set git hooks path to local directory
git config core.hooksPath scripts/githooks

# Make hooks executable
chmod +x scripts/githooks/pre-commit
chmod +x scripts/githooks/commit-msg

echo "✅ Git hooks configured successfully!"

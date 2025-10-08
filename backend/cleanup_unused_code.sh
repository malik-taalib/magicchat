#!/bin/bash

# Magic Chat - Cleanup Unused Code Script
# This script removes old code that has been replaced by vertical slice architecture

set -e  # Exit on error

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     Magic Chat - Unused Code Cleanup Script                   ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "📂 Current directory: $PWD"
echo ""

# Function to safely remove directory
remove_dir() {
    local dir=$1
    if [ -d "$dir" ]; then
        echo "🗑️  Removing: $dir"
        rm -rf "$dir"
        echo "   ✅ Deleted"
    else
        echo "⚠️  Not found: $dir (already deleted?)"
    fi
}

echo "Starting cleanup of unused code..."
echo ""

# Remove old internal structure
echo "1️⃣  Removing internal/ directory (old horizontal structure)"
remove_dir "internal"
echo ""

# Remove old API router
echo "2️⃣  Removing api/ directory (old router)"
remove_dir "api"
echo ""

# Remove old gateway entry point
echo "3️⃣  Removing cmd/gateway/ (old entry point)"
remove_dir "cmd/gateway"
echo ""

# Remove old auth service
echo "4️⃣  Removing cmd/auth-service/ (duplicate auth logic)"
remove_dir "cmd/auth-service"
echo ""

# Remove old microservice entry points
echo "5️⃣  Removing cmd/video-service/ (abandoned microservice)"
remove_dir "cmd/video-service"
echo ""

echo "6️⃣  Removing cmd/feed-service/ (abandoned microservice)"
remove_dir "cmd/feed-service"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Cleanup complete!"
echo ""

# Verify what's left in cmd/
echo "📁 Contents of cmd/ directory:"
ls -la cmd/ 2>/dev/null || echo "   cmd/ directory structure:"
find cmd -type d -maxdepth 1 2>/dev/null | sort
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔨 Verifying build..."
echo ""

# Test if the code still builds
if go build -o /tmp/magicchat-test cmd/server/main.go; then
    echo "✅ Build successful! The application still compiles."
    rm /tmp/magicchat-test
else
    echo "❌ Build failed! Something went wrong."
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Summary:"
echo ""
echo "  Removed directories:"
echo "    • internal/"
echo "    • api/"
echo "    • cmd/gateway/"
echo "    • cmd/auth-service/"
echo "    • cmd/video-service/"
echo "    • cmd/feed-service/"
echo ""
echo "  Active code structure:"
echo "    ✅ cmd/server/        - Main entry point"
echo "    ✅ slices/            - 7 vertical slices"
echo "    ✅ pkg/               - Shared infrastructure"
echo "    ✅ migrations/        - Database setup"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 All unused code has been removed!"
echo ""
echo "Next steps:"
echo "  1. Review the changes: git status"
echo "  2. Test the application: make run"
echo "  3. Commit the cleanup: git add . && git commit -m 'Remove unused code'"
echo ""

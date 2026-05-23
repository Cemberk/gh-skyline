#!/bin/bash

# Build script for GitHub Skyline WASM module
# This script compiles the Go code to WebAssembly using TinyGo

set -e

echo "Building GitHub Skyline WASM module..."

# Check if TinyGo is installed
if ! command -v tinygo &> /dev/null; then
    echo "TinyGo is not installed. Installing..."
    echo "Please install TinyGo from: https://tinygo.org/getting-started/install/"
    echo "Or use: brew install tinygo (on macOS)"
    exit 1
fi

# Navigate to the wasm directory
cd "$(dirname "$0")/wasm"

# Build the WASM module
echo "Compiling Go to WebAssembly..."
tinygo build -o ../build/skyline.wasm -target wasm ./

# Copy the wasm_exec.js file from TinyGo
TINYGO_ROOT=$(tinygo env TINYGOROOT)
cp "$TINYGO_ROOT/targets/wasm_exec.js" ../static/js/

echo "Build complete!"
echo "WASM module: web/build/skyline.wasm"
echo "WASM exec: web/static/js/wasm_exec.js"
echo "Size: $(du -h ../build/skyline.wasm | cut -f1)"

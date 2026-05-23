# Build Directory

This directory contains the compiled WASM binary that powers the web interface.

## Contents

- `skyline.wasm` - The WebAssembly binary compiled from Go code

## Note

The WASM file is generated during the build process and is not committed to the repository. It is automatically built and deployed by GitHub Actions when changes are pushed to the main branch.

To build locally, run:
```bash
cd web
./build.sh
```

The generated file will be approximately 500KB in size.

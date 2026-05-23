# GitHub Skyline Web Interface

A browser-based version of GitHub Skyline that runs entirely in your browser using WebAssembly. No installation required!

## Features

- 🌐 **No Installation Required**: Generate your GitHub contribution skyline directly in your browser
- ⚡ **Fast & Secure**: All processing happens client-side using WebAssembly
- 📥 **Direct Download**: Download your STL file immediately
- 🎨 **ASCII Preview**: See a text-based preview of your contributions
- 🔒 **Privacy-Focused**: Your data never leaves your browser

## Usage

1. Visit the [GitHub Skyline Web App](https://cemberk.github.io/gh-skyline/)
2. Enter your GitHub username
3. Select the year you want to generate
4. (Optional) Enter a GitHub Personal Access Token for higher API rate limits
5. Click "Generate Skyline"
6. Preview the ASCII art
7. Download your STL file

## GitHub Personal Access Token

While optional, providing a GitHub Personal Access Token increases API rate limits from 60 to 5000 requests per hour.

To create a token:
1. Go to [GitHub Settings > Personal Access Tokens](https://github.com/settings/tokens/new?scopes=&description=GitHub%20Skyline)
2. No scopes are required (leave all checkboxes unchecked)
3. Click "Generate token"
4. Copy and paste the token into the web interface

**Note**: The token is only used for API requests and is never stored or sent anywhere except to GitHub's API.

## Technology

- **WebAssembly (WASM)**: Compiled from Go using TinyGo
- **GitHub GraphQL API**: Fetches contribution data
- **Vanilla JavaScript**: No framework dependencies
- **GitHub Pages**: Static hosting

## Development

### Building the WASM Module

Requirements:
- [TinyGo](https://tinygo.org/getting-started/install/)

Build the WASM module:

```bash
cd web
./build.sh
```

This will:
1. Compile the Go code to WebAssembly
2. Copy the WASM binary to `web/build/skyline.wasm`
3. Copy the WASM exec helper to `web/static/js/wasm_exec.js`

### Local Development

To test locally:

```bash
# Start a local web server
cd docs
python3 -m http.server 8000
```

Then open http://localhost:8000 in your browser.

### Deployment

The web interface is automatically deployed to GitHub Pages when changes are pushed to the main branch. The deployment is handled by GitHub Actions (see `.github/workflows/deploy-pages.yml`).

## Architecture

```
web/
├── wasm/                    # Go source code for WASM
│   ├── main.go             # WASM entry point and JS bindings
│   ├── geometry.go         # 3D geometry calculations
│   ├── stl.go              # STL file generation
│   ├── ascii.go            # ASCII art generation
│   └── types.go            # Data structures
├── build/                   # Build output
│   └── skyline.wasm        # Compiled WASM binary
└── build.sh                 # Build script

docs/                        # GitHub Pages source
├── index.html              # Main page
├── css/
│   └── styles.css          # Styling
├── js/
│   ├── app.js              # Main application logic
│   ├── github.js           # GitHub API integration
│   ├── wasm-loader.js      # WASM loading and initialization
│   └── wasm_exec.js        # TinyGo WASM runtime (generated)
└── build/
    └── skyline.wasm        # WASM binary (copied from web/build)
```

## Browser Support

The web interface requires a modern browser with WebAssembly support:

- ✅ Chrome/Edge 57+
- ✅ Firefox 52+
- ✅ Safari 11+
- ✅ Opera 44+

## Limitations

- **API Rate Limits**: Without a token, GitHub API is limited to 60 requests/hour
- **Single Year Only**: The web interface currently only supports single-year skylines
- **Text Rendering**: The web version does not include embossed text or logos (kept simple for smaller WASM size)

## Comparison with CLI Version

| Feature | CLI Extension | Web Interface |
|---------|--------------|---------------|
| Installation | Requires `gh` CLI | None |
| Multi-year support | ✅ Yes | ❌ No (planned) |
| Embossed text | ✅ Yes | ❌ No |
| Logo image | ✅ Yes | ❌ No |
| Offline use | ✅ Yes | ❌ No |
| Cross-platform | ✅ Yes | ✅ Yes |
| File size | ~10MB | ~500KB WASM |

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

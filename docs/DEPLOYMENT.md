# GitHub Skyline - Web Deployment Guide

This guide explains how the web deployment works and how to build and deploy the GitHub Skyline web interface.

## Overview

The web interface is a browser-based version of GitHub Skyline that:
- Runs entirely in the browser using WebAssembly (WASM)
- Fetches contribution data directly from GitHub's API
- Generates STL files client-side
- Requires no server-side processing

## Architecture

### Components

1. **WASM Module** (`web/wasm/`)
   - Go code compiled to WebAssembly using TinyGo
   - Handles all geometry calculations and STL generation
   - Exposes functions to JavaScript

2. **Web Frontend** (`docs/`)
   - HTML/CSS/JavaScript interface
   - GitHub API integration
   - WASM loader and file download functionality

3. **GitHub Actions** (`.github/workflows/deploy-pages.yml`)
   - Automated build and deployment to GitHub Pages
   - Runs on every push to main branch

## Building Locally

### Prerequisites

1. **TinyGo**: Install from https://tinygo.org/getting-started/install/
   ```bash
   # macOS
   brew install tinygo
   
   # Ubuntu/Debian
   wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.2/tinygo_0.31.2_amd64.deb
   sudo dpkg -i tinygo_0.31.2_amd64.deb
   ```

2. **Go**: Version 1.25.0 or later (already required for the main project)

### Build Steps

1. Navigate to the web directory:
   ```bash
   cd web
   ```

2. Run the build script:
   ```bash
   ./build.sh
   ```

3. The script will:
   - Compile Go code to WASM
   - Copy the WASM binary to `web/build/skyline.wasm`
   - Copy TinyGo's `wasm_exec.js` to `web/static/js/`

### Testing Locally

1. Start a local web server:
   ```bash
   cd docs
   python3 -m http.server 8000
   ```

2. Open http://localhost:8000 in your browser

3. Test the functionality:
   - Enter a GitHub username
   - Select a year
   - Click "Generate Skyline"
   - Download the STL file

## Deployment Process

### Automatic Deployment

The web interface is automatically deployed to GitHub Pages when changes are pushed to the main branch.

The workflow (`.github/workflows/deploy-pages.yml`) performs these steps:

1. **Checkout code**: Gets the latest code from the repository
2. **Set up Go**: Installs Go based on go.mod
3. **Install TinyGo**: Downloads and installs TinyGo
4. **Build WASM**: Compiles the Go code to WebAssembly
5. **Copy files**: Copies WASM binary and helper files to docs/
6. **Deploy**: Publishes to GitHub Pages

### Manual Deployment

If you need to deploy manually:

1. Build the WASM module:
   ```bash
   cd web
   ./build.sh
   ```

2. Copy build artifacts to docs:
   ```bash
   mkdir -p docs/build
   cp web/build/skyline.wasm docs/build/
   cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js docs/js/
   ```

3. Commit and push:
   ```bash
   git add docs/
   git commit -m "Update web deployment"
   git push
   ```

## Configuration

### GitHub Pages Settings

To enable GitHub Pages for your fork:

1. Go to repository Settings
2. Navigate to Pages (left sidebar)
3. Under "Source", select "GitHub Actions"
4. The site will be available at `https://[username].github.io/gh-skyline/`

### Custom Domain (Optional)

To use a custom domain:

1. Add a `CNAME` file to the `docs/` directory with your domain
2. Configure DNS settings at your domain provider
3. Enable "Enforce HTTPS" in GitHub Pages settings

## Troubleshooting

### WASM Module Not Loading

**Problem**: "Failed to load WASM module" error

**Solutions**:
- Check browser console for specific errors
- Ensure WASM file is accessible (check network tab)
- Verify `wasm_exec.js` is in `docs/js/`
- Clear browser cache and reload

### Build Fails

**Problem**: Build script fails with TinyGo errors

**Solutions**:
- Verify TinyGo is installed: `tinygo version`
- Ensure Go version matches go.mod
- Check for syntax errors in Go files
- Try cleaning and rebuilding: `rm -rf web/build && ./web/build.sh`

### GitHub API Rate Limits

**Problem**: "API rate limit exceeded" error

**Solutions**:
- Wait for rate limit to reset (check console for reset time)
- Provide a GitHub Personal Access Token
- Use authenticated requests (creates token with no scopes)

### STL File Issues

**Problem**: STL file is corrupted or won't open in 3D software

**Solutions**:
- Check browser console for generation errors
- Verify contribution data was fetched successfully
- Try a different year or username
- Check file size (should be > 1KB)

## Performance Considerations

### WASM File Size

The WASM binary should be kept under 5MB for good user experience:
- Current size: ~500KB (compressed)
- TinyGo produces smaller binaries than standard Go
- Avoid importing large libraries

### Load Time

Typical load times:
- WASM module: 1-2 seconds
- Contribution data fetch: 1-2 seconds
- STL generation: < 1 second
- Total: 2-5 seconds

### Optimization Tips

1. **Compress WASM**: Serve with gzip/brotli compression
2. **Cache WASM**: Set appropriate cache headers
3. **Minimize JS**: Bundle and minify JavaScript files
4. **CDN**: Use GitHub Pages CDN for global distribution

## Security

### GitHub API Tokens

- Tokens are only used for API authentication
- Never stored or sent to any server except GitHub
- Stored only in memory during the session
- No scopes required (read-only public data)

### WASM Sandboxing

- WASM runs in a sandboxed environment
- No file system access
- No network access (except through JavaScript)
- Cannot access user's system

### Content Security Policy

Consider adding CSP headers for production:
```
Content-Security-Policy: default-src 'self'; script-src 'self' 'wasm-unsafe-eval'; connect-src https://api.github.com
```

## Monitoring

### Analytics (Optional)

To add analytics:

1. Add Google Analytics or similar to `docs/index.html`
2. Track key events:
   - Page loads
   - STL generations
   - Downloads
   - Errors

### Error Tracking

Consider adding error tracking:
- Sentry
- Rollbar
- LogRocket

## Future Enhancements

Potential improvements:

- [ ] Multi-year support
- [ ] 3D preview using Three.js
- [ ] Text embossing
- [ ] Logo support
- [ ] Export to other formats (OBJ, PLY)
- [ ] Share functionality
- [ ] Dark/light theme toggle
- [ ] Mobile app (PWA)

## Contributing

To contribute to the web interface:

1. Make changes to `web/` or `docs/`
2. Test locally
3. Update documentation
4. Submit pull request
5. Ensure CI passes

See [CONTRIBUTING.md](../CONTRIBUTING.md) for general guidelines.

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [WebAssembly MDN](https://developer.mozilla.org/en-US/docs/WebAssembly)
- [GitHub GraphQL API](https://docs.github.com/en/graphql)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)

## License

Same as the main project: MIT License

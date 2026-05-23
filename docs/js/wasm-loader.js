// WASM Loader module
// Handles loading and initializing the WebAssembly module

let wasmModule = null;
let wasmInstance = null;

/**
 * Loads and initializes the WASM module
 * @param {string} wasmPath - Path to the .wasm file
 * @returns {Promise<void>}
 */
export async function loadWASM(wasmPath) {
    if (wasmModule && wasmInstance) {
        console.log('WASM already loaded');
        return;
    }

    try {
        const go = new Go();
        let result;
        
        // Try instantiateStreaming first (more efficient)
        try {
            result = await WebAssembly.instantiateStreaming(fetch(wasmPath), go.importObject);
        } catch (streamError) {
            console.warn('instantiateStreaming failed, falling back to fetch + instantiate:', streamError);
            
            // Fallback: fetch as ArrayBuffer then instantiate
            // This works even if the MIME type is incorrect
            const response = await fetch(wasmPath);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const wasmBytes = await response.arrayBuffer();
            result = await WebAssembly.instantiate(wasmBytes, go.importObject);
        }
        
        wasmInstance = result.instance;
        wasmModule = result.module;
        
        // Run the Go program
        go.run(wasmInstance);

        // Wait for WASM to be ready
        await waitForWASMReady();

        console.log('WASM module loaded successfully');
    } catch (error) {
        console.error('Error loading WASM:', error);
        throw new Error(`Failed to load WASM module: ${error.message}`);
    }
}

/**
 * Waits for the WASM module to signal it's ready
 * @returns {Promise<void>}
 */
function waitForWASMReady() {
    return new Promise((resolve, reject) => {
        const checkReady = () => {
            if (window.wasmReady) {
                resolve();
            } else {
                setTimeout(checkReady, 100);
            }
        };
        
        checkReady();

        // Timeout after 10 seconds
        setTimeout(() => {
            if (!window.wasmReady) {
                reject(new Error('WASM initialization timeout'));
            }
        }, 10000);
    });
}

/**
 * Generates STL binary data using the WASM module
 * @param {Array} contributions - 3D array of contribution data [year][week][day]
 * @param {string} username - GitHub username
 * @param {number} startYear - Start year
 * @param {number} endYear - End year
 * @returns {Uint8Array} STL binary data
 */
export function generateSTL(contributions, username, startYear, endYear) {
    if (!window.generateSTL) {
        throw new Error('WASM module not loaded or generateSTL function not available');
    }

    try {
        const contributionsJSON = JSON.stringify(contributions);
        const stlData = window.generateSTL(contributionsJSON, username, startYear, endYear);
        return stlData;
    } catch (error) {
        console.error('Error generating STL:', error);
        throw error;
    }
}

/**
 * Generates ASCII art using the WASM module
 * @param {Array} contributions - 2D array of contribution data [week][day]
 * @returns {string} ASCII art string
 */
export function generateASCII(contributions) {
    if (!window.generateASCII) {
        throw new Error('WASM module not loaded or generateASCII function not available');
    }

    try {
        const contributionsJSON = JSON.stringify(contributions);
        const ascii = window.generateASCII(contributionsJSON);
        return ascii;
    } catch (error) {
        console.error('Error generating ASCII:', error);
        throw error;
    }
}

/**
 * Checks if WASM module is loaded and ready
 * @returns {boolean}
 */
export function isWASMReady() {
    return window.wasmReady === true;
}

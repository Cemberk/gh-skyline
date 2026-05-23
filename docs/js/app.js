// Main application logic for GitHub Skyline web interface
// This module coordinates the UI, WASM, and GitHub API interactions

import { fetchContributions, formatContributionsForWASM, formatRateLimit } from './github.js';
import { loadWASM, generateSTL, generateASCII, isWASMReady } from './wasm-loader.js';

// State
let currentSTLData = null;
let currentYear = new Date().getFullYear();
let currentUsername = '';

// DOM Elements
const elements = {
    usernameInput: document.getElementById('username'),
    yearInput: document.getElementById('year'),
    tokenInput: document.getElementById('token'),
    generateBtn: document.getElementById('generate-btn'),
    downloadBtn: document.getElementById('download-btn'),
    loading: document.getElementById('loading'),
    loadingMessage: document.getElementById('loading-message'),
    error: document.getElementById('error'),
    errorMessage: document.getElementById('error-message'),
    asciiSection: document.getElementById('ascii-section'),
    asciiOutput: document.getElementById('ascii-output'),
    downloadSection: document.getElementById('download-section'),
    stats: document.getElementById('stats'),
    rateLimit: document.getElementById('rate-limit')
};

/**
 * Initialize the application
 */
async function init() {
    console.log('Initializing GitHub Skyline web app...');
    
    // Set up event listeners
    elements.generateBtn.addEventListener('click', handleGenerate);
    elements.downloadBtn.addEventListener('click', handleDownload);
    
    // Set current year as default
    elements.yearInput.value = currentYear;
    
    // Load WASM module
    try {
        showLoading('Loading WASM module...');
        await loadWASM('../build/skyline.wasm');
        hideLoading();
        console.log('WASM module ready');
    } catch (error) {
        showError(`Failed to load WASM module: ${error.message}`);
        console.error(error);
    }
}

/**
 * Handle generate button click
 */
async function handleGenerate() {
    hideError();
    hideResults();
    
    // Get input values
    const username = elements.usernameInput.value.trim();
    const year = parseInt(elements.yearInput.value);
    const token = elements.tokenInput.value.trim() || null;
    
    // Validate inputs
    if (!username) {
        showError('Please enter a GitHub username');
        return;
    }
    
    if (year < 2008 || year > new Date().getFullYear()) {
        showError(`Year must be between 2008 and ${new Date().getFullYear()}`);
        return;
    }
    
    // Check if WASM is ready
    if (!isWASMReady()) {
        showError('WASM module is not ready. Please refresh the page and try again.');
        return;
    }
    
    currentUsername = username;
    currentYear = year;
    
    try {
        // Disable generate button
        elements.generateBtn.disabled = true;
        
        // Fetch contributions from GitHub API
        showLoading(`Fetching contributions for ${username} (${year})...`);
        const { contributions, rateLimit } = await fetchContributions(username, year, token);
        
        // Update rate limit display
        elements.rateLimit.textContent = formatRateLimit(rateLimit);
        
        // Format data for WASM
        const formattedContributions = formatContributionsForWASM(contributions);
        
        // Generate ASCII art
        showLoading('Generating ASCII preview...');
        const ascii = generateASCII(formattedContributions);
        elements.asciiOutput.textContent = ascii;
        elements.asciiSection.classList.remove('hidden');
        
        // Generate STL
        showLoading('Generating 3D model...');
        // Wrap in array for single year (WASM expects [year][week][day])
        const contributionsArray = [formattedContributions];
        currentSTLData = generateSTL(contributionsArray, username, year, year);
        
        // Show download section
        const fileSize = (currentSTLData.byteLength / 1024).toFixed(2);
        elements.stats.textContent = `Total contributions: ${contributions.totalContributions} | STL file size: ${fileSize} KB`;
        elements.downloadSection.classList.remove('hidden');
        
        hideLoading();
        
    } catch (error) {
        showError(`Error: ${error.message}`);
        console.error('Error during generation:', error);
    } finally {
        elements.generateBtn.disabled = false;
    }
}

/**
 * Handle download button click
 */
function handleDownload() {
    if (!currentSTLData) {
        showError('No STL data available. Please generate a skyline first.');
        return;
    }
    
    try {
        // Create blob from STL data
        const blob = new Blob([currentSTLData], { type: 'application/octet-stream' });
        
        // Create download link
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${currentUsername}-${currentYear}-github-skyline.stl`;
        
        // Trigger download
        document.body.appendChild(a);
        a.click();
        
        // Cleanup
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        
        console.log('STL file downloaded successfully');
    } catch (error) {
        showError(`Error downloading file: ${error.message}`);
        console.error('Download error:', error);
    }
}

/**
 * Show loading state with optional message
 */
function showLoading(message = 'Loading...') {
    elements.loadingMessage.textContent = message;
    elements.loading.classList.remove('hidden');
}

/**
 * Hide loading state
 */
function hideLoading() {
    elements.loading.classList.add('hidden');
}

/**
 * Show error message
 */
function showError(message) {
    elements.errorMessage.textContent = message;
    elements.error.classList.remove('hidden');
}

/**
 * Hide error message
 */
function hideError() {
    elements.error.classList.add('hidden');
}

/**
 * Hide result sections
 */
function hideResults() {
    elements.asciiSection.classList.add('hidden');
    elements.downloadSection.classList.add('hidden');
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}

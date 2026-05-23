// Package main provides the WebAssembly interface for gh-skyline.
// This package exposes Go functions to JavaScript for generating GitHub contribution skylines.
package main

import (
	"encoding/json"
	"syscall/js"
)

// main is the entry point for the WASM module.
// It registers JavaScript callable functions and keeps the module running.
func main() {
	// Register functions that can be called from JavaScript
	js.Global().Set("generateSTL", js.FuncOf(generateSTLWrapper))
	js.Global().Set("generateASCII", js.FuncOf(generateASCIIWrapper))
	js.Global().Set("wasmReady", js.ValueOf(true))

	// Log that WASM is ready
	js.Global().Get("console").Call("log", "GitHub Skyline WASM module loaded successfully")

	// Keep the program running
	<-make(chan struct{})
}

// generateSTLWrapper is the JavaScript-callable wrapper for STL generation.
// It takes contribution data, username, and year range as arguments and returns STL binary data.
//
// JavaScript signature:
// generateSTL(contributionsJSON: string, username: string, startYear: number, endYear: number) => Uint8Array
func generateSTLWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 {
		return jsError("generateSTL requires 4 arguments: contributionsJSON, username, startYear, endYear")
	}

	contributionsJSON := args[0].String()
	username := args[1].String()
	startYear := args[2].Int()
	endYear := args[3].Int()

	// Parse contributions JSON
	var contributions [][][]ContributionDay
	if err := json.Unmarshal([]byte(contributionsJSON), &contributions); err != nil {
		return jsError("Failed to parse contributions JSON: " + err.Error())
	}

	// Generate STL binary
	stlData, err := GenerateSTLBytes(contributions, username, startYear, endYear)
	if err != nil {
		return jsError("Failed to generate STL: " + err.Error())
	}

	// Convert Go byte slice to JavaScript Uint8Array
	uint8Array := js.Global().Get("Uint8Array").New(len(stlData))
	js.CopyBytesToJS(uint8Array, stlData)

	return uint8Array
}

// generateASCIIWrapper is the JavaScript-callable wrapper for ASCII art generation.
// It takes contribution data and returns ASCII art as a string.
//
// JavaScript signature:
// generateASCII(contributionsJSON: string) => string
func generateASCIIWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("generateASCII requires 1 argument: contributionsJSON")
	}

	contributionsJSON := args[0].String()

	// Parse contributions JSON
	var contributions [][]ContributionDay
	if err := json.Unmarshal([]byte(contributionsJSON), &contributions); err != nil {
		return jsError("Failed to parse contributions JSON: " + err.Error())
	}

	// Generate ASCII art
	ascii, err := GenerateASCII(contributions)
	if err != nil {
		return jsError("Failed to generate ASCII art: " + err.Error())
	}

	return js.ValueOf(ascii)
}

// jsError creates a JavaScript Error object
func jsError(message string) js.Value {
	errorConstructor := js.Global().Get("Error")
	return errorConstructor.New(message)
}

{
    "name": "test-zig",
    "version": "0.1.0",
    "description": "Zig tests for wazero-small-cache",
    "authors": ["pantopic"],
    "target": "wasm32-wasi",
    "dependencies": {
        "small_cache": {
            "path": "../sdk-zig"
        }
    },
    "buildOptions": {
        "target": "wasm32-wasi",
        "optimize": "ReleaseSmall"
    }
}
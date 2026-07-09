const std = @import("std");

// WASM memory management
var id: u64 = 0;
var key_cap: u32 = 256;
var key_len: u32 = 0;
var key: [256]u8 = undefined;
var val_cap: u32 = 4 << 10; // 4 KiB
var val_len: u32 = 0;
var val: [4096]u8 = undefined;
var meta: [8]u32 = undefined;

// Exported function to get the metadata pointer
// This is equivalent to the Go __small_cache function
export fn __small_cache() u32 {
    // Fill the metadata with pointers to our variables
    meta[0] = @intCast(@intFromPtr(&id));
    meta[1] = @intCast(@intFromPtr(&key_cap));
    meta[2] = @intCast(@intFromPtr(&key_len));
    meta[3] = @intCast(@intFromPtr(&key));
    meta[4] = @intCast(@intFromPtr(&val_cap));
    meta[5] = @intCast(@intFromPtr(&val_len));
    meta[6] = @intCast(@intFromPtr(&val));

    // Return the pointer to metadata
    return @intCast(@intFromPtr(&meta));
}

// WASM module exports for cache operations
extern fn __small_cache_put() void;
extern fn __small_cache_get() void;
extern fn __small_cache_del() void;
extern fn __small_cache_min() void;

// Local cache struct
pub const Local = struct {
    id: u64,

    // Get method
    pub fn get(self: *Local, k: []u8) []u8 {
        id = self.id;
        key_len = @intCast(k.len);

        // Copy key data to our global buffer
        const copy_len = if (k.len < key_cap) k.len else key_cap;
        @memcpy(key[0..copy_len], k[0..copy_len]);

        __small_cache_get();

        return val[0..val_len];
    }

    // Put method
    pub fn put(self: *Local, k: []u8, v: []u8) void {
        id = self.id;
        key_len = @intCast(k.len);
        val_len = @intCast(v.len);

        // Copy key data to our global buffer
        const key_copy_len = if (k.len < key_cap) k.len else key_cap;
        @memcpy(key[0..key_copy_len], k[0..key_copy_len]);

        // Copy value data to our global buffer
        const val_copy_len = if (v.len < val_cap) v.len else val_cap;
        @memcpy(val[0..val_copy_len], v[0..val_copy_len]);

        __small_cache_put();
    }

    // Del method
    pub fn del(self: *Local, k: []u8) void {
        id = self.id;
        key_len = @intCast(k.len);

        // Copy key data to our global buffer
        const copy_len = if (k.len < key_cap) k.len else key_cap;
        @memcpy(key[0..copy_len], k[0..copy_len]);

        __small_cache_del();
    }

    // Min method
    pub fn min(self: *Local) []u8 {
        id = self.id;
        __small_cache_min();
        return key[0..key_len];
    }
};

// Constructor for Local
pub fn newLocal(local_id: u64) *Local {
    return &Local{ .id = local_id };
}

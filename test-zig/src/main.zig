const std = @import("std");
const builtin = @import("builtin");

// Import the small_cache library
const small_cache = @import("small_cache");

const SMALL_CACHE_ID_TEST_1 = 0;
const SMALL_CACHE_ID_TEST_2 = 1;

var testLocalCache1: *small_cache.Local = undefined;
var testLocalCache2: *small_cache.Local = undefined;

// Initialize the caches
fn init() void {
    testLocalCache1 = small_cache.newLocal(SMALL_CACHE_ID_TEST_1);
    testLocalCache2 = small_cache.newLocal(SMALL_CACHE_ID_TEST_2);
}

// Helper function to convert u64 to byte array (Little Endian)
fn u64ToBytes(value: u64) [8]u8 {
    var result: [8]u8 = undefined;
    result[0] = @intCast(value);
    result[1] = @intCast(value >> 8);
    result[2] = @intCast(value >> 16);
    result[3] = @intCast(value >> 24);
    result[4] = @intCast(value >> 32);
    result[5] = @intCast(value >> 40);
    result[6] = @intCast(value >> 48);
    result[7] = @intCast(value >> 56);
    return result;
}

// Helper function to convert byte array to u64 (Little Endian)
fn bytesToU64(bytes: [8]u8) u64 {
    return @as(u64, bytes[0]) |
        (@as(u64, bytes[1]) << 8) |
        (@as(u64, bytes[2]) << 16) |
        (@as(u64, bytes[3]) << 24) |
        (@as(u64, bytes[4]) << 32) |
        (@as(u64, bytes[5]) << 40) |
        (@as(u64, bytes[6]) << 48) |
        (@as(u64, bytes[7]) << 56);
}

// Exported functions for testing
export fn testLocalPut(k: u64, v: u64) void {
    const key_bytes = u64ToBytes(k);
    const val_bytes = u64ToBytes(v);
    testLocalCache1.put(&key_bytes, &val_bytes);
}

export fn testLocalGet(k: u64) u64 {
    const key_bytes = u64ToBytes(k);
    const result = testLocalCache1.get(&key_bytes);

    // Check if we got 8 bytes back
    if (result.len != 8) {
        return 0;
    }

    return bytesToU64(@as([8]u8, @bitCast(result)));
}

export fn testLocalDel(k: u64) void {
    const key_bytes = u64ToBytes(k);
    testLocalCache1.del(&key_bytes);
}

export fn testLocalMin() u64 {
    const result = testLocalCache1.min();

    // Check if we got 8 bytes back
    if (result.len != 8) {
        return 0;
    }

    return bytesToU64(@as([8]u8, @bitCast(result)));
}

export fn testLocalPut2(k: u64, v: u64) void {
    const key_bytes = u64ToBytes(k);
    const val_bytes = u64ToBytes(v);
    testLocalCache2.put(&key_bytes, &val_bytes);
}

export fn testLocalGet2(k: u64) u64 {
    const key_bytes = u64ToBytes(k);
    const result = testLocalCache2.get(&key_bytes);

    // Check if we got 8 bytes back
    if (result.len != 8) {
        return 0;
    }

    return bytesToU64(@as([8]u8, @bitCast(result)));
}

export fn testLocalDel2(k: u64) void {
    const key_bytes = u64ToBytes(k);
    testLocalCache2.del(&key_bytes);
}

export fn testLocalMin2() u64 {
    const result = testLocalCache2.min();

    // Check if we got 8 bytes back
    if (result.len != 8) {
        return 0;
    }

    return bytesToU64(@as([8]u8, @bitCast(result)));
}

// Entry point
pub fn main() void {
    init();
}

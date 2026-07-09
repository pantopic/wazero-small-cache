const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Add the small_cache library as a dependency - reference the SDK directory directly
    // const small_cache = b.dependency("small_cache", .{
    //     .path = "../sdk-zig",
    //     .target = target,
    //     .optimize = optimize,
    // });

    const exe = b.addExecutable(.{
        .name = "test_zig",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .target = target,
            .optimize = optimize,
        }),
    });
    const small_cache = b.addModule("small_cache", .{
        .root_source_file = b.path("../sdk-zig/src/lib.zig"),
    });

    // Import the small_cache library
    // exe.root_module.addImport("small_cache", small_cache.module("small_cache"));
    exe.root_module.addImport("small_cache", small_cache);

    b.installArtifact(exe);

    const run_cmd = b.addRunArtifact(exe);
    run_cmd.step.dependOn(b.getInstallStep());

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);

    const lib_unit_tests = b.addTest(.{
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .target = target,
            .optimize = optimize,
        }),
    });

    // Add the small_cache dependency to the test as well
    // lib_unit_tests.root_module.addImport("small_cache", small_cache.module("small_cache"));
    lib_unit_tests.root_module.addImport("small_cache", small_cache);

    const run_lib_unit_tests = b.addRunArtifact(lib_unit_tests);

    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_lib_unit_tests.step);
}

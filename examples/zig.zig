const std = @import("std");

pub fn main() !void {
    var buf: [1024]u8 = undefined;

    while (true) {
        const n = try std.fs.File.stdin().read(&buf);
        if (n == 0) break;
        try std.fs.File.stdout().writeAll(buf[0..n]);
    }
}

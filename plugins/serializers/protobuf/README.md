# Protocol Buffers:

Protocol buffers are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data.

https://developers.google.com/protocol-buffers/

### Message definitions:

```protobuf
syntax="proto3";

package protobuf;

message field_value {
    oneof value {
        string string_value = 1;
        double float_value = 2;
        int64 int_value = 3;
        bool bool_value = 4;
    }
}

message metric {
    string name = 1;
    map<string, string> tags = 2;
    map<string, field_value> fields = 3;
    int64 timestamp = 4;
}
```

* `timestamp` is a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC.

### Protocol Buffers Configuration:

```toml
[[outputs.file]]
  ## Files to write to, "stdout" is a specially handled file.
  files = ["stdout", "/tmp/metrics.out"]

  ## Data format to output.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "protobuf"

  prepend_length = true
```

Output formats like `file` requires a mechanism to identify the beginning and the end point of each message. In Protocol Buffers, there is no way to do it by the format itself. To solve this, you can set the optional parameter `prepend_length` to `true`. If this parameter is set, the length of each message will be prepended in front of real message.

```
--------------------------------------------------------------------------------------------
|       N       |          message         |       N2      |          message         |  ...
--------------------------------------------------------------------------------------------
     4 bytes               N bytes              4 bytes               N2 bytes
```

Lengths are little endian uint32(4 bytes).
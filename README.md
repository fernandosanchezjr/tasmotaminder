# tasmotaminder

Manage Tasmota smart plugs connected to an MQTT server with a simple yaml configuration.

## Supported devices

* [EZPlug+ by TH3D](https://www.th3dstudio.com/product/ezplug-open-source-wifi-smart-plug/)

## Compilation

From the root of the git repo:

```shell
go build ./cmd/tasmotaminder
```

This will result in a `tasmotaminder` binary in the root folder. Copy the file to a convenient location, such as `/usr/local/bin/tasmotaminder`.

As a basic sense-check, run the program now without configuration. You should see output similar to the following:

```shell
$ tasmotaminder
2024/04/10 11:39:24 error opening rule config file: open /etc/tasmotaminder/rules.yaml: no such file or directory
```

# tasmotaminder

Manage Tasmota smart plugs connected to an MQTT server with a simple yaml configuration.

## Supported devices

* [EZPlug+ by TH3D](https://www.th3dstudio.com/product/ezplug-open-source-wifi-smart-plug/)

## Supported MQTT brokers

* [Eclipse Mosquitto](https://mosquitto.org/)

## Requirements

* Running MQTT broker: tasmotaminder connects to an MQTT server, listening to Tasmota topics for plugs specified in the configuration file.
* Smart plugs already connected to MQTT broker

## Compilation

From the root of the git repo, run the following command:

```shell
go build ./cmd/tasmotaminder
```

This will result in a `tasmotaminder` binary in the root folder. Copy the file to a convenient location, such as `/usr/local/bin/tasmotaminder`.

As a basic sense-check, run the program now without configuration. You should see output similar to the following:

```shell
$ tasmotaminder
2024/04/10 11:39:24 error opening rule config file: open /etc/tasmotaminder/rules.yaml: no such file or directory
```

## Environment Variables

`tasmotaminder` relies on environment variables to resolve the MQTT broker to connect to:

* `BROKER_HOST`: hostname of the MQTT broker. If not provided, defaults to `localhost`.
* `BROKER_PORT`: port of the MQTT broker. If not provided, defaults to `1883`.
* `CLIENT_ID`: client id provided to the MQTT broker. If not provided, defaults to `tasmotaminder`.
* `BROKER_USERNAME`: username to use when connecting to the MQTT broker.
* `BROKER_PASSWORD`: password to use when connecting to the MQTT broker.
* `RULE_CONFIG_YAML`: location of yaml config file containing plug rules.

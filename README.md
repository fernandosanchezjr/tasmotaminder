# tasmotaminder

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoffee.com/fernandosanchezjr)

Manage Tasmota smart plugs connected to an MQTT broker with a simple yaml configuration.

```yaml
# bike plug: let this plug run until the device is consuming less than 8 Watts, when the 
# light goes green on the charger
- deviceId: EZPlug_6E6729
  powerTimer:
    # device consumes 8 Watts or less
    power: 8
    # power comparison options: lessThan (<), greaterThan (>=), equalTo (==)
    # defaults to greaterThan.
    powerComparison: lessThan
    action: off

# outside xmas lights plug - turn on at 5 PM, then turn off at 10 PM
- deviceId: EZPlug_8B4EB2
  # no limit on the number of schedules allowed for a plug - can be combined with 
  # powerTimer entries
  powerSchedules:
    # cron spec for 17:00 every day
    - cron: 0 17 * * *
      action: on
    # cron spec for 22:00 every day
    - cron: 0 22 * * *
      action: off
```

## Features

* Turn smart plugs off/on/reset after consuming greaterThan/lessThan/equalTo watts for a specified runtime
* Turn smart plugs on/off on schedules defined by simple cron notation
* No additional system dependencies required
* Simple YAML-based configuration
* Small footprint: 8 MiB RAM usage on average. CPU use negligible. Perfect for Raspberry Pis and other constrained devices.
* Private: only needs access to your MQTT broker and nothing else

## Supported devices

* [EZPlug+ by TH3D](https://www.th3dstudio.com/product/ezplug-open-source-wifi-smart-plug/)

## Supported MQTT brokers

* [Eclipse Mosquitto](https://mosquitto.org/)

## Requirements

* Running MQTT broker: tasmotaminder connects to an MQTT server, listening to Tasmota topics for plugs specified in the configuration file.
* Smart plugs already connected to MQTT broker
* Optional
  * System user: if running as a systemd unit, create a system user as follows `sudo adduser --system --no-create-home --group --disabled-login tasmotaminder`. Replace `tasmotaminder` with your desired username. This will create a group with the same name, which is convenient for admin users.
  * Config folder: if running as a system service, instead of compartmentalized in a docker container, it is useful to set up reasonable permissions on the config folder: `sudo mkdir -m 770 -p /etc/tasmotaminder && sudo chown -R tasmotaminder:tasmotaminder /etc/tasmotaminder`

## Install pre-compiled binaries

* Download the correct package for your architecture from the [Releases](https://github.com/fernandosanchezjr/tasmotaminder/releases) page
* Extract the .tar.gz package
* Copy the `tasmotaminder` binary to an appropriate location, such as `/usr/local/bin/tasmotaminder` on linux
* As a basic sense-check, run the program now without configuration. You should see output similar to the following:
  ```shell
  $ tasmotaminder
  2024/04/10 11:39:24 error opening rule config file: open /etc/tasmotaminder/rules.yaml: no such file or directory
  ```

## Compilation

* Clone this repo
* From the root of the git repo, run the following command:
  ```shell
  go build ./cmd/tasmotaminder
  ```
* You should now have a `tasmotaminder` binary in the root folder. Copy the file to a convenient location, such as `/usr/local/bin/tasmotaminder`.
* As a basic sense-check, run the program now without configuration. You should see output similar to the following:
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
* `RULE_CONFIG_YAML`: location of yaml config file containing plug rules. If not provided, defaults to `/etc/tasmotaminder/rules.yaml`.

## Example Systemd unit

The following example systemd unit can be used to run tasmotaminder. Place the following file in `/etc/systemd/system/tasmotaminder.service`:

```
[Unit]
Description=Tasmota Minder
Documentation=https://github.com/fernandosanchezjr/tasmotaminder
After=network-online.target

[Service]
# user and group created using instructions from Requirements
User=tasmotaminder
Group=tasmotaminder

# config folder created using instructions from Requirements
EnvironmentFile=/etc/tasmotaminder/settings.sh

Restart=on-failure
ExecStart=/usr/local/bin/tasmotaminder
KillSignal=SIGINT

[Install]
WantedBy=multi-user.target
```

### Example Systemd environment file

The following example environment file can be stored in `/etc/tasmotaminder/settings.sh` to use with the Systemd unit.

```shell
BROKER_HOST=127.0.0.1
BROKER_USERNAME=someusername
BROKER_PASSWORD=hopefullynotthispassword
```

## Rules file

Here is an example config file stored in `/etc/tasmotaminder/rules.yaml`:

```yaml
# coffee maker plug: if it consumes 1 Watt or more, let it run for 1 hour, then reset
- deviceId: EZPlug_10BFE3
  # the coffeemaker will forget about the coffee cycle if powered off for a few seconds
  resetDurationSeconds: 300
  powerTimer:
    # device consumes more than 1 Watt
    power: 1
    # let the device run for 3600 seconds
    runtimeSeconds: 3600
    # reset after runtime
    action: reset

# drumset plug: my daughter forgets to turn this off, so let her play for an hour
- deviceId: EZPlug_8B4E91
  powerTimer:
    # device consumes 1 Watt or more
    power: 1
    # let the device run for 3600 seconds
    runtimeSeconds: 3600
    # power off after runtime
    action: off

# bike plug: let this plug run until the device is consuming less than 8 Watts, when the 
# light goes green on the charger
- deviceId: EZPlug_6E6729
  powerTimer:
    # device consumes 8 Watts or less
    power: 8
    # power comparison options: lessThan (<), greaterThan (>=), equalTo (==)
    # defaults to greaterThan
    powerComparison: lessThan
    action: off

# outside xmas lights plug - turn on at 5 PM, then turn off at 10 PM
- deviceId: EZPlug_8B4EB2
  # no limit on the number of schedules allowed for a plug - can be combined with 
  # powerTimer entries
  powerSchedules:
    # cron spec for 17:00 every day
    - cron: 0 17 * * *
      action: on
    # cron spec for 22:00 every day
    - cron: 0 22 * * *
      action: off

# front window xmas lights plug - turn on at 7 PM, then turn off at 9 PM
- deviceId: EZPlug_8B490A
  powerSchedules:
  # cron spec for 19:00 every day
  - cron: 0 19 * * *
    action: on
  # cron spec for 21:00 every day
  - cron: 0 21 * * *
    action: off
```

### Loading and reloading settings

Settings and environment variables are only read at start time. If you change any rule or environment variable you must restart the service for the changes to take effect. If using the suggested Systemd unit, restarting is possible with `systemctl restart tasmotaminder.service`.

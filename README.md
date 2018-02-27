# IPFIX Translator

IPFIX Translator is translator piece of code for IPFIX messages as received on Kafka Bus.
IPFIX messages from different devices are received by [vFlow](https://github.com/VerizonDigital/vflow) module and it then sends them on Kafka bus.
Once IPFIX Translator receives the messages on Kafka bus, it parses the data and pushes to different databases based on configuration. Currently it pushes the data only to Appformix Data Manager (DM).

### Build
```
make build
```
It creates ```ipfix-translator``` binary which provides Command Line Interface (CLI) utility, it can be invoked as below
```
$ ./ipfix-translator --help
NAME:
   Kafka Consumer CLI - A new cli application

USAGE:
   ipfix-translator [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     kafka-consumer  Kafka Consumer
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

To run the IPFIX Translator:
```
./ipfix-translator kafka-consumer --config-file /etc/ipfix-translator-config.yml
```
# Configuration Parameters
The configuration file (ipfix-translator-config.yml) is an yml file with the below possible configurations
```
verbose: False
log-file: "/var/log/ipfix-translator.log"
kafka-broker-list: "127.0.0.1:9092"

data-manager-ip: "127.0.0.1"
data-manager-port: "9000"
```
```verbose:``` Boolean, if Verbose mode is on or off
```log-file:``` The file path where the log file should be created
```kafka-broker-list:``` Kafka Broker List
```data-manager-ip:``` IP Address where Data Manager is running
```data-manager-port:``` Port Data Manager is listening on


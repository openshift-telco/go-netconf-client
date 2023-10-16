## Testing

In order to test the NETCONF client, we are using a NETCONF simulator.

The NETCONF simulator is a Java project provided by [Pantheon Technologies](https://github.com/PANTHEONtech)

It requires Java and Maven to be built.

Here are the steps to get the simulator, and emulate 200 devices.

**Clone the repository**
~~~
git clone https://github.com/PANTHEONtech/lighty-netconf-simulator.git
~~~
**Compile the repository**
~~~
mvn clean install
~~~
**Run the simulator**
~~~
// RPC tests
java -jar lighty-netconf-simulator/examples/devices/lighty-toaster-multiple-devices/target/lighty-toaster-multiple-devices-15.0.1-SNAPSHOT.jar --starting-port 20000 --device-count 200 --thread-pool-size 200

// Notification tests
java -jar lighty-netconf-simulator/examples/devices/lighty-notifications-device/target/lighty-notifications-device-15.0.1-SNAPSHOT.jar 12346
~~~
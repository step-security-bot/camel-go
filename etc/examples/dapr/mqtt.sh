mosquitto_pub -h test.mosquitto.org -p 1883 -t camel-iot -m '{ "source": "sensor-1", "data": "foo" }'
mosquitto_pub -h test.mosquitto.org -p 1883 -t camel-iot -m '{ "source": "sensor-1", "data": "foo" }' | jq .
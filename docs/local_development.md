
# Local development


This guide configures an MQTT broker and simulates Zuul reports using the mosquitto_pub client, which sends messages in place of an actual Zuul instance. This approach eliminates the need to have a real Zuul instance running.

1. Clone the repository
```bash
git clone https://github.com/SovereignCloudStack/zuul-mqtt-matrix-bridge.git
```

2. Deploy MQTT broker container only
```bash
cd ./zuul-mqtt-matrix-bridge/docker
docker-compose up -d mqtt-broker
```

3. Run zuul-mqtt-matrix-bridge. Set up at least the `matrix-token` and `matrix-room-id` variables to guarantee that you will see reports in the specific Matrix room of your choice. If you don't wish to send messages to the Matrix chat, you can simply enable debug mode to verify that the bridge correctly consumes MQTT messages and attempts to relay them to the Matrix chat
```bash
cd ..
GO_LOG=debug go run . -matrix-token <token> -matrix-room-id <room-id> -mqtt-broker tcp://127.0.0.1:1883 -mqtt-user mqtt-user -mqtt-pass secret
```

4. Intall mosquitto_pub client, e.g.:

```bash
apt-get install mosquitto-clients
```

5. Publish a set of Zuul messages and observe how the zuul-mqtt-matrix-bridge processes them

<details open>
  <summary>Build start message</summary>

  ```bash
  mosquitto_pub -h 127.0.0.1 -u mqtt-user -P secret -q 1 -d -t zuul/test -m '{"action":"start","tenant":"openstack.org","pipeline":"check","project":"sf-jobs","branch":"master","change_url":"https://gerrit.example.com/r/3","message":"Starting check jobs.","trigger_time":1524801056.2545865,"enqueue_time":1524801093.5689456,"change":"3","patchset":"1","commit_id":"2db20c7fb26adf9ac9936a9e750ced9b4854a964","owner":"username","ref":"refs/changes/03/3/1","zuul_ref":"Zf8b3d7cd34f54cb396b488226589db8f","buildset":{"uuid":"f8b3d7cd34f54cb396b488226589db8f","builds":[{"job_name":"linters","voting":true}]}}'
  ```
</details>

<details open>
  <summary>Build succeeded message</summary>

  ```bash
  mosquitto_pub -h 127.0.0.1 -u mqtt-user -P secret -q 1 -d -t zuul/test -m '{"action":"success","tenant":"openstack.org","pipeline":"check","project":"sf-jobs","branch":"master","change_url":"https://gerrit.example.com/r/3","message":"Build succeeded.","trigger_time":1524801056.2545864,"enqueue_time":1524801093.5689457,"change":"3","patchset":"1","commit_id":"2db20c7fb26adf9ac9936a9e750ced9b4854a964","owner":"username","ref":"refs/changes/03/3/1","zuul_ref":"Zf8b3d7cd34f54cb396b488226589db8f","buildset":{"uuid":"f8b3d7cd34f54cb396b488226589db8f","builds":[{"job_name":"linters","voting":true,"uuid":"16e3e55aca984c6c9a50cc3c5b21bb83","execute_time":1524801120.7563295,"start_time":1524801179.8557224,"end_time":1524801208.928095,"log_url":"https://logs.example.com/logs/3/3/1/check/linters/16e3e55/","web_url":"https://tenant.example.com/t/tenant-one/build/16e3e55aca984c6c9a50cc3c5b21bb83/","result":"SUCCESS","dependencies":[],"artifacts":[]}]}}'
  ```
</details>

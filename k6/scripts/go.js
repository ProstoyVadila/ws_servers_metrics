import ws from "k6/ws";
import { sleep, check } from "k6";

const SLEEP_TIME = 0.4;
const GO_APP_URL = __ENV.GO_APP_URL
  ? `${__ENV.GO_APP_URL}`
  : "ws://localhost:8000/";

export const options = {
  scenarios: {
    constVUWith100Messages: {
      executor: "constant-vus",
      exec: "constVUWithDirectMessages",
      vus: 80,
      startTime: "0",
      duration: "45s",
      tags: { test_type: "3_direct_messages" },
    },
    rampUp50With100Messages: {
      executor: "ramping-vus",
      exec: "rampUp50With100Messages",
      startVUs: 5,
      startTime: "2ms",
      tags: { test_type: "100_messages" },
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
        { target: 100, duration: "30s" },
        { target: 150, duration: "30s" },
        { target: 200, duration: "30s" }, // 3m30s
      ],
    },
    rampUp50WithSingleMessage: {
      executor: "ramping-vus",
      exec: "rampUp50WithSingleMessage",
      startVUs: 5,
      startTime: "6m",
      tags: { test_type: "signle_message" },
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
        { target: 100, duration: "30s" },
        { target: 150, duration: "30s" },
        { target: 200, duration: "30s" },
        { target: 250, duration: "30s" },
        { target: 300, duration: "30s" },
        { target: 350, duration: "30s" },
        { target: 400, duration: "30s" },
        { target: 450, duration: "30s" },
        { target: 500, duration: "30s" },
        { target: 550, duration: "30s" },
        { target: 600, duration: "30s" },
        { target: 650, duration: "30s" },
        { target: 700, duration: "30s" },
        { target: 750, duration: "30s" }, // 9m30s
      ],
    },
  },
};

export function rampUp50WithSingleMessage() {
  sleep(SLEEP_TIME);
  const data = {
    action_type: "broadcast",
    body: `${Date.now()}`,
  };

  // const url = getRandomAppURL();

  const app = ws.connect(GO_APP_URL, (socket) => {
    socket.on("open", function open() {
      console.debug(`VU ${__VU}: connected`);

      socket.send(JSON.stringify(data));
      sleep(SLEEP_TIME);

      socket.on("message", (msg) => {
        console.debug(`VU ${__VU} recieved msg: ${msg}`);
      });

      socket.on("close", () => {
        console.debug(`VU ${__VU} disconnected`);
      });

      socket.on("error", (err) => {
        console.debug(`VU ${__VU} got error: ${err}`);
      });
      socket.close();
    });
  });

  check(app, {
    "Connected successfully": (r) => r && r.status === 101,
  });
}

export function rampUp50With100Messages() {
  sleep(SLEEP_TIME);
  const data = {
    action_type: "broadcast",
    body: `${Date.now()}`,
  };

  const url = getRandomAppURL();

  const app = ws.connect(url, (socket) => {
    socket.on("open", function open() {
      console.debug(`VU ${__VU}: connected`);

      socket.send(JSON.stringify(data));
      sleep(SLEEP_TIME);

      for (let i = 0; i < 100; i++) {
        sleep(SLEEP_TIME);
        socket.send(JSON.stringify(data));
      }
      socket.send(JSON.stringify(data));

      socket.on("message", (msg) => {
        console.debug(`VU ${__VU} recieved msg: ${msg}`);
      });

      socket.on("close", () => {
        console.debug(`VU ${__VU} disconnected`);
      });

      socket.on("error", (err) => {
        console.debug(`VU ${__VU} got error: ${err}`);
      });
      socket.close();
    });
  });

  check(app, {
    "Connected successfully": (r) => r && r.status === 101,
  });
}

export function constVUWithDirectMessages() {
  sleep(SLEEP_TIME);
  const data = {
    action_type: "direct",
    body: `${Date.now()}`,
  };

  const url = getRandomAppURL();

  const app = ws.connect(url, (socket) => {
    socket.on("open", function open() {
      console.debug(`VU ${__VU}: connected`);

      socket.send(JSON.stringify(data));
      sleep(SLEEP_TIME);
      let flag = true;

      socket.on("message", (msg) => {
        console.debug(`VU ${__VU} recieved msg: ${msg}`);
        sleep(SLEEP_TIME);
        if (flag) {
          socket.send(JSON.stringify(data));
        } else {
          socket.close();
        }
      });

      socket.on("close", () => {
        console.debug(`VU ${__VU} disconnected`);
      });

      socket.on("error", (err) => {
        console.debug(`VU ${__VU} got error: ${err}`);
      });
      socket.close();
    });
  });

  check(app, {
    "Connected successfully": (r) => r && r.status === 101,
  });
}

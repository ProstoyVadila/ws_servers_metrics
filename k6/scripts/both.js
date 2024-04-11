import ws from "k6/ws";
import { sleep, check } from "k6";

const SLEEP_TIME = 0.4;
const RUST_APP_URL = "ws://localhost:8000/ws";
const PYTHON_APP_URL = "ws://localhost:8001/ws";

const randomIntBetween = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1) + min);
};
const getRandomAppURL = () => {
  const oneOrZero = Math.random() >= 0.5 ? 1 : 0;
  console.debug(
    oneOrZero === 1 ? "connect to rust app" : "connect to python app"
  );

  const url = oneOrZero === 1 ? RUST_APP_URL : PYTHON_APP_URL;
  return url;
};
//
// const sessionDuration = randomIntBetween(10000, 60000); // 10s ~ 60s

export const options = {
  scenarios: {
    constVUWith100Messages: {
      executor: "constant-vus",
      exec: "constVUWithDirectMessages",
      vus: 80,
      startTime: "0",
      duration: "45s",
    },
    rampUp50With100Messages: {
      executor: "ramping-vus",
      exec: "rampUp50With100Messages",
      startVUs: 5,
      startTime: "1ms",
      tags: { test_type: "100_messages" },
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
        { target: 100, duration: "30s" },
        { target: 150, duration: "30s" },
        { target: 200, duration: "30s" },
        { target: 250, duration: "30s" }, // 4m
      ],
    },
    rampUp50WithSingleMessage: {
      executor: "ramping-vus",
      exec: "rampUp50WithSingleMessage",
      startVUs: 5,
      startTime: "4m20s",
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

  const url = getRandomAppURL();

  const app = ws.connect(url, (socket) => {
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

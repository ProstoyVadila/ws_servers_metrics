import ws from "k6/ws";
import { sleep, check } from "k6";

const randomIntBetween = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1) + min);
};

const RUST_APP_URL = "ws://localhost:8000/ws";
const PYTHON_APP_URL = "ws://localhost:8001/ws";

const sessionDuration = randomIntBetween(10000, 60000); // 10s ~ 60s
const SLEEP_TIME = 0.5;

export const options = {
  scenarios: {
    single_broadcast_msg_per_user: {
      executor: "ramping-vus",
      startVUs: 3,
      // iteration: 20,
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
        { target: 100, duration: "30s" },
        { target: 150, duration: "30s" },
        { target: 200, duration: "30s" }, // it's gonna kill my mac
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
        { target: 750, duration: "30s" },
      ],
    },
  },
};

export default function () {
  sleep(SLEEP_TIME);
  const data = {
    action_type: "broadcast",
    body: `${Date.now()}`,
  };

  const oneOrZero = Math.random() >= 0.5 ? 1 : 0;
  console.debug(
    oneOrZero === 1 ? "connect to rust app" : "connect to python app"
  );

  const url = oneOrZero === 1 ? RUST_APP_URL : PYTHON_APP_URL;

  const app = ws.connect(url, (socket) => {
    socket.on("open", function open() {
      console.debug(`VU ${__VU}: connected`);

      socket.send(JSON.stringify(data));
      sleep(SLEEP_TIME);

      // for (let i = 0; i < 100; i++) {
      //   sleep(SLEEP_TIME);
      //   socket.send(JSON.stringify(data));
      // }
      // socket.send(JSON.stringify(data))

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

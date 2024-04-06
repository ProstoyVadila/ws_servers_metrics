import ws from "k6/ws";
import { sleep, check } from "k6";

const randomIntBetween = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1) + min);
};

// const RUST_APP_URL = "ws://rust_ws_app:8000/ws";
// const PYTHON_APP_URL = "ws://python_ws_app:8001/ws";
const URL = "ws://localhost:8001/ws";

const sessionDuration = randomIntBetween(10000, 60000); // 10s ~ 60s
const SLEEP_TIME = 0.01;

export const options = {
  scenarios: {
    contacts: {
      executor: "ramping-vus",
      startVUs: 3,
      // iteration: 20,
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
      ],
    },
  },
};

export default function () {
  const data = {
    action_type: "broadcast",
    body: `${Date.now()}`,
  };

  const app = ws.connect(URL, (socket) => {
    socket.on("open", function open() {
      console.log(`VU ${__VU}: connected`);
      socket.send(JSON.stringify(data));

      // for (let i = 0; i < 100; i++) {
      //   sleep(SLEEP_TIME);
      //   socket.send(JSON.stringify(data));
      // }
      // socket.send(JSON.stringify(data))

      socket.on("message", (msg) => {
        console.log(`VU ${__VU} recieved msg: ${msg}`);
      });

      socket.on("close", () => {
        console.log(`VU ${__VU} disconnected`);
      });

      socket.on("error", (err) => {
        console.log(`VU ${__VU} got error: ${err}`);
      });
      socket.close();
    });
  });

  check(app, {
    "Connected successfully": (r) => r && r.status === 101,
  });
}

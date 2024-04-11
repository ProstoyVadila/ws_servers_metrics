import ws from "k6/ws";
import { sleep, check } from "k6";

const randomIntBetween = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1) + min);
};

const PYTHON_APP_URL = "ws://localhost:8001/ws";

const sessionDuration = randomIntBetween(10000, 60000); // 10s ~ 60s
const SLEEP_TIME = 0.1;

export const options = {
  scenarios: {
    single_broadcast_msg_per_user: {
      executor: "ramping-vus",
      startVUs: 3,
      stages: [
        { target: 20, duration: "30s" },
        { target: 40, duration: "30s" },
        { target: 60, duration: "30s" },
        { target: 80, duration: "30s" },
        { target: 100, duration: "30s" },
        { target: 150, duration: "30s" },
        { target: 200, duration: "30s" },
        // { target: 250, duration: "30s" },
        // { target: 300, duration: "30s" },
        // { target: 350, duration: "30s" },
        // { target: 400, duration: "30s" },
        // { target: 450, duration: "30s" },
        // { target: 500, duration: "30s" },
      ],
    },
  },
};

export default function () {
  const data = {
    action_type: "broadcast",
    body: `${Date.now()}`,
  };

  const res_python_app = ws.connect(PYTHON_APP_URL, (socket) => {
    socket.on("open", function open() {
      console.debug(`VU ${__VU}: connected`);

      socket.send(JSON.stringify(data));

      for (let i = 0; i < 100; i++) {
        sleep(SLEEP_TIME);
        socket.send(JSON.stringify(data));
      }

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

  check(res_python_app, {
    "Connected successfully": (r) => r && r.status === 101,
  });
}

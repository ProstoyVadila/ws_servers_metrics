import asyncio
import logging
import signal
import os
# from collections import defaultdict
from functools import partial


import uvloop
from tornado import web, ioloop, locks
from prometheus_client import start_http_server

from routes import routes


shutdown_event = locks.Event()


class Application(web.Application):
    def __init__(self, routes: list[str]) -> None:
        settings = {
            "websocket_ping_interval": 1,
            "websocket_ping_timeout": 2,
        }
        self.users = set()
        super().__init__(routes, **settings)


async def on_shutdown(app: Application) -> None:
    logging.info("graceful shutdown")
    app.users.clear()
    shutdown_event.set()


def exit_handler(app: Application, *args, **kwargs) -> None:
    ioloop.IOLoop.instance().add_callback_from_signal(on_shutdown, app)


async def main() -> None:
    init_logging(logging.DEBUG)
    app = Application(routes=routes)

    port = os.environ.get('PORT', 8001)
    metrics_port = os.environ.get('METRICS_PORT', 8002)
    logging.info(f"serving on http://0.0.0.0:{port} and metrics on http://0.0.0.0:{metrics_port}")
    app.listen(
        port=int(port),
    )
    # prometheus metrics
    start_http_server(port=int(metrics_port))

    signal.signal(signal.SIGTERM, partial(exit_handler, app))
    signal.signal(signal.SIGINT, partial(exit_handler, app))
    await shutdown_event.wait()


def init_logging(level: int | str) -> None:
    logging.basicConfig(format="%(asctime)-15s | %(levelname).1s | %(name)8s |  %(message)s", level=level)
    logger = logging.getLogger("ws-server")
    logger.setLevel(level)


if __name__ == "__main__":
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    ioloop.IOLoop.current().run_sync(main)

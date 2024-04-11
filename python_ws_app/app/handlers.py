import sys
import logging
from uuid import uuid4

import orjson
from datetime import datetime, UTC
from tornado import websocket, web
from pydantic import ValidationError

from models import WsMessage, ActionType
from metrics import (
    WS_CONNECTIONS,
    ALL_WS_CONNECTIONS_TOTAL,
    WS_BROADCAST_DURATION_SECONDS,
    WS_MESSAGE_HANDLING_DURATION_SECONDS,
    WS_CONN_CLOSED_ERRORS_TOTAL
)


class WsHandler(websocket.WebSocketHandler):

    def check_origin(self, origin: str) -> bool:
        return True

    def select_subprotocol(self, subprotocols: websocket.List[str]) -> str | None:
        if subprotocols:
            return subprotocols[0]
        return None

    async def open(self) -> websocket.Awaitable[None] | None:
        ALL_WS_CONNECTIONS_TOTAL.inc()
        WS_CONNECTIONS.inc()
        logging.info("Adding connection")
        self.user_id = str(uuid4())
        self.application.users.add(self)
        await self.write_message(message=orjson.dumps({"accept": True}))


    def on_close(self) -> None:
        logging.info("removing connection")
        WS_CONNECTIONS.dec()
        self.application.users.discard(self)

    @WS_MESSAGE_HANDLING_DURATION_SECONDS.time()
    async def on_message(self, message: str | bytes) -> websocket.Awaitable[None] | None:
            if isinstance(message, bytes):
                await self.write_message(message=message)
            elif isinstance(message, str):
                await self.validate_action(message=message)
            else:
                logging.warning("unknown message type %s", message)

    async def validate_action(self, message: str) -> websocket.Awaitable[None] | None:
        try:
            msg = orjson.loads(message)
            ws_msg = WsMessage(**msg)
            ws_msg.user_id = ws_msg.user_id or self.user_id
            ws_msg.body = ws_msg.body or "successful"

            match ws_msg.action_type:
                case ActionType.DIRECT:
                    ws_msg.data = datetime.now(UTC).timestamp()
                    await self.write_message(orjson.dumps(ws_msg.dict()))
                case ActionType.BROADCAST:
                    ws_msg.data = str(sys.version)
                    await self.broadcast_message(ws_msg=ws_msg)
                case ActionType.PING:
                    ws_msg.data = "pong"
                    await self.send(ws_msg=ws_msg)
                case ActionType.PONG:
                    ws_msg.data = "ping"
                    await self.send(ws_msg=ws_msg)

        except websocket.WebSocketClosedError:
            WS_CONN_CLOSED_ERRORS_TOTAL.inc()
            logging.exception("Websocket connection was closed")

        except ValidationError:
            logging.exception("Cannot validate message")
            ws_msg = WsMessage(
                user_id=self.user_id or str(uuid4()),
                action_type=ActionType.BROADCAST,
                body="cannot validate a message",
                data="cannot validate a message"
            )
            await self.send(ws_msg=ws_msg)
        except Exception:
            logging.exception("Error validating message")
            ws_msg = WsMessage(
                user_id=self.user_id or str(uuid4()),
                action_type=ActionType.BROADCAST,
                body="cannot validate a message",
                data="cannot validate a message"
            )
            await self.send(ws_msg=ws_msg)

    @WS_BROADCAST_DURATION_SECONDS.time()
    async def broadcast_message(self, ws_msg: WsMessage) -> websocket.Awaitable[None] | None:
        for _con in self.application.users:
            try:
                _con.write_message(message=orjson.dumps(ws_msg.dict()))
            except Exception:
                logging.exception("Error handling message")

    async def send(self, ws_msg: WsMessage) -> websocket.Awaitable[None] | None:
        self.write_message(message=orjson.dumps(ws_msg.dict()))

class RootHandler(web.RequestHandler):
    async def get(self) -> web.Awaitable[None] | None:
        self.write("python_ws_app")

class HealthcheckHandler(web.RequestHandler):
    async def get(self) -> web.Awaitable[None] | None:
        self.write("ok")


class MetricHandler(web.RequestHandler):
    async def get(self) -> web.Awaitable[None] | None:
        connections = len(self.application.users)
        ws_connections_total = f"ws_connections_total {connections}" 
        metrics = [ws_connections_total]
        self.write("\n".join(metrics))

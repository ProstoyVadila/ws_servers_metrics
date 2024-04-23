import asyncio
import logging
import uuid
import sys
from datetime import UTC, datetime

import websockets
import orjson
import uvloop
from websockets.server import serve
from websockets import WebSocketServerProtocol

from models import WsMessage, ActionType

logger = logging.getLogger('websockets')
logger.setLevel(logging.DEBUG)
logger.addHandler(logging.StreamHandler())

CONNECTIONS = set()

async def handle_msg(ws: WebSocketServerProtocol, msg: str):
    try:
        msg = orjson.loads(msg)
        ws_msg = WsMessage(**msg)
        ws_msg.user_id = ws_msg.user_id or uuid.uuid4()
        ws_msg.body = ws_msg.body or "successful"
        match ws_msg.action_type:
            case ActionType.DIRECT:
                ws_msg.data = datetime.now(UTC).timestamp()
                ws.send(orjson.dumps(ws_msg.model_dump()))
            case ActionType.BROADCAST:
                ws_msg.data = str(sys.version)
                websockets.broadcast(CONNECTIONS, orjson.dumps(ws_msg.model_dump()))
    except Exception:
        logger.exception("cannot handle a message")


async def ws_handler(websocket: WebSocketServerProtocol):
    CONNECTIONS.add(websocket)
    try:
        message = await websocket.recv()
        if isinstance(message, bytes):
            websocket.send(message)
        elif isinstance(message, str):
            await handle_msg(websocket, message)
    except Exception:
        logger.exception("got error")
    finally:
        CONNECTIONS.remove(websocket)


async def main():
    print("ws server is running...")
    async with serve(ws_handler, "localhost", 8001):
        await asyncio.Future()

if __name__ == "__main__":   
    uvloop.run(main())

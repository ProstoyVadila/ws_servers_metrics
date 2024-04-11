from enum import StrEnum

from pydantic import BaseModel

class ActionType(StrEnum):
    DIRECT = "direct"
    BROADCAST = "broadcast"
    PING = "ping"
    PONG = "pong"


class WsMessage(BaseModel):
    user_id: str | None
    action_type: ActionType
    body: str
    data: str | None

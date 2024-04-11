use serde::{self, Deserialize, Serialize};
use serde_json::json;


#[derive(PartialEq, Eq, Clone, Serialize, Deserialize)]
pub enum ActionType {
    #[serde(alias="direct", alias="DIRECT")]
    Direct,
    #[serde(alias="broadcast", alias="BROADCAST")]
    Broadcast,
    #[serde(alias="ping", alias="ping")]
    Ping,
    #[serde(alias="pong", alias="pong")]
    Pong,
}


#[derive(Clone, Serialize, Deserialize)]
pub struct WsMessage {
    pub user_id: Option<usize>,
    pub action_type: ActionType,
    pub body: String,
    pub data: Option<String>,
}

impl WsMessage {
    pub fn new(
        user_id: usize, 
        action_type: ActionType, 
        body: String, 
        data: Option<String>
    ) -> WsMessage {
        WsMessage {
            user_id: Some(user_id), 
            action_type,
            body, 
            data,
        }
    }
}

impl ToString for WsMessage {
    fn to_string(&self) -> String {
        json!(self).to_string()
    }
}

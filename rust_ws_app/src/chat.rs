
use std::collections::HashMap;

use rocket::{
    futures::{stream::SplitSink, SinkExt}, tokio::sync::Mutex,
};
use rocket_ws::{Message, stream::DuplexStream};
use log;

use crate::models::{WsMessage, ActionType};
// use crate::metrics::WS_BROADCAST_DURATION_SECONDS;


pub struct ChatRoomConnection {
    pub user_id: usize,
    pub sink: SplitSink<DuplexStream, Message>,
}

impl ChatRoomConnection {
    fn new(user_id: usize, sink: SplitSink<DuplexStream, Message>) -> ChatRoomConnection {
        ChatRoomConnection {
            user_id,
            sink,
        }
    }
}

#[derive(Default)]
pub struct ChatRoom {
    pub connections: Mutex<HashMap<usize, ChatRoomConnection>>,
}

impl ChatRoom {
    pub async fn add(&self, user_id: usize, ws_sink: SplitSink<DuplexStream, Message>) {
        let mut conns = self.connections.lock().await;
        let connection = ChatRoomConnection::new(user_id.clone(), ws_sink);
        conns.insert(user_id, connection);
    }

    pub fn parse_message(&self, msg: String) -> Option<WsMessage> {
        let new_msg: WsMessage = match serde_json::from_str(msg.as_str()) {
            Ok(new_msg) => new_msg,
            Err(err) => {
                log::error!("Cannot deserialize json message: {}", err);
                return None;
            }
        };
        Some(new_msg)
    }

    pub async fn handle_message(&self, user_id: usize, msg: String) {
        if let Some(msg) = self.parse_message(msg.clone()) {
            match msg.action_type {
                ActionType::Direct => self.handle_direct(user_id, msg).await,
                ActionType::Broadcast => self.handle_broadcast(user_id, msg).await,
                ActionType::Ping => {
                    todo!()
                },
                ActionType::Pong => {
                    todo!()
                }
            }
        }
    }

    pub async fn handle_direct(&self, user_id: usize, msg: WsMessage) {
        // get timestamp for latency
        let now = chrono::Utc::now();
        let data = now.timestamp_nanos_opt().and_then(|now| Some(now.to_string()));
        let msg = WsMessage::new(user_id, ActionType::Direct, msg.body, data);

        self.send_direct_message(user_id, msg).await;
    }

    pub async fn send_direct_message(&self, user_id: usize, msg: WsMessage) {
        let mut conns = self.connections.lock().await;
        let conn = match conns.get_mut(&user_id) {
            Some(conn) => conn,
            None => {
                log::error!("cannot find connection for user {}", user_id);
                return;
            }
        };
        let _ = conn.sink.send(Message::Text(msg.to_string())).await;
    }

    pub async fn handle_broadcast(&self, user_id: usize, msg: WsMessage) {
        let data = Some("rust version".to_string()); // temp
        let user_id = match msg.user_id {
            Some(id) => id,
            _ => user_id,
        };
        let new_msg = WsMessage::new(user_id, ActionType::Broadcast, msg.body, data);
        self.broadcast_message(new_msg).await;
    }

    pub async fn broadcast_message(&self, msg: WsMessage) {
        // let timer = WS_BROADCAST_DURATION_SECONDS.start_timer();
        let mut conns = self.connections.lock().await;
        for (_id, conn) in conns.iter_mut() {
            let _ = conn.sink.send(Message::Text(msg.to_string())).await;
        }
        // timer.stop_and_record();
    }

    pub async fn flush(&self, user_id: usize) {
        log::debug!("removing user {}", user_id);
        let mut conns = self.connections.lock().await;
        let _ = match conns.remove(&user_id) {
            Some(conn) => conn,
            _ => {
                log::warn!("Cannot find a user {} to remove", user_id);
                return;
            }
        };
    }

}

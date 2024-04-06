use std::sync::atomic::{AtomicUsize, Ordering};

use rocket::{
    futures::StreamExt, 
    State
};
use rocket_ws::{Channel, Message, WebSocket};
use log;

use crate::chat::ChatRoom;
use crate::metrics::{ALL_WS_CONNECTIONS_TOTAL, WS_CONNECTIONS, WS_MESSAGE_HANDLING_DURATION_SECONDS};

static USER_ID_COUNTER: AtomicUsize = AtomicUsize::new(1);


#[rocket::get("/")]
pub fn chat<'r>(ws: WebSocket, state: &'r State<ChatRoom>) -> Channel<'r> {
    ws.channel(move |stream| Box::pin(async move {
        let user_id = USER_ID_COUNTER.fetch_add(1, Ordering::Relaxed);
        let (ws_sink, mut ws_stream) = stream.split();

        state.add(user_id, ws_sink).await;
        ALL_WS_CONNECTIONS_TOTAL.inc();
        WS_CONNECTIONS.inc();
    
        while let Some(msg) = ws_stream.next().await {
            if let Ok(msg_content) = msg {
                match msg_content {
                    Message::Text(json_msg) => {
                        let timer = WS_MESSAGE_HANDLING_DURATION_SECONDS.start_timer();
                        state.handle_message(user_id, json_msg).await;
                        timer.stop_and_record();
                    },
                    Message::Ping(_) => {},
                    Message::Pong(_) => {},
                    _ => {
                        // Unsupported
                        log::warn!("Unsupported message type {}", msg_content);

                    }
                }
            }
        }
        state.flush(user_id).await;
        WS_CONNECTIONS.dec();
    
        Ok(())
    }))
}

mod chat;
mod handlers;
mod metrics;
mod models;
mod routes;

#[rocket::main]
async fn main() {
    env_logger::init();
    let prom = metrics::get_prometheus();

    log::info!("Starting ws server on http://localhost:8082");
    let _ = rocket::build()
        .attach(prom.clone())
        .mount("/", routes::get_routes())
        .mount("/metrics", prom)
        .manage(chat::ChatRoom::default())
        .launch()
        .await;

    log::info!("Ws server is stopped")
}

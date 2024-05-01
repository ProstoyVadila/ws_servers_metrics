use crate::handlers;

pub fn get_routes() -> Vec<rocket::Route> {
    let routes = rocket::routes![
        handlers::get_chat,
        handlers::get_healthcheck,
        handlers::get_root
    ];
    routes
}

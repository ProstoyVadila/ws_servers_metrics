use rocket;

use crate::handlers;




pub fn get_routes() -> Vec<rocket::Route> {
    let routes = rocket::routes![
        handlers::chat,
    ];
    routes
}
from handlers import WsHandler, Handler, MetricHandler

routes = [
    (r"/", Handler),
    (r"/ws", WsHandler),
    # (r"/metrics", MetricHandler)
]
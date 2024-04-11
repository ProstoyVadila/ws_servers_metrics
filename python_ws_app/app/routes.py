from handlers import WsHandler, RootHandler, HealthcheckHandler, MetricHandler

routes = [
    (r"/", RootHandler),
    (r"/healthcheck", HealthcheckHandler),
    (r"/ws", WsHandler),
    # (r"/metrics", MetricHandler)
]
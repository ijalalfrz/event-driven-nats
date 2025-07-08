import tornado.web
import tornado.log
import tornado.options
import tornado.ioloop
import sqlite3
import logging
import json
import time
import signal
import nats

class App(tornado.web.Application):

    def __init__(self, handlers, **kwargs):
        super().__init__(handlers, **kwargs)

        # Initialising db connection
        self.db = sqlite3.connect("listings.db")
        self.db.row_factory = sqlite3.Row
        self.init_db()

    def init_nats(self, connection_url):
        self.nats_conn = nats.connect(connection_url)

    def init_db(self):
        cursor = self.db.cursor()

        # Create table
        cursor.execute(
            "CREATE TABLE IF NOT EXISTS 'listings' ("
            + "id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,"
            + "user_id INTEGER NOT NULL,"
            + "listing_type TEXT NOT NULL,"
            + "price INTEGER NOT NULL,"
            + "created_at INTEGER NOT NULL,"
            + "updated_at INTEGER NOT NULL"
            + ");"
        )
        self.db.commit()

    def close_db(self):
        """Close database connection gracefully"""
        if hasattr(self, 'db') and self.db:
            self.db.close()
            logging.info("Database connection closed")

class BaseHandler(tornado.web.RequestHandler):
    def write_json(self, obj, status_code=200):
        self.set_header("Content-Type", "application/json")
        self.set_status(status_code)
        self.write(json.dumps(obj))

# /listings
class ListingsHandler(BaseHandler):
    @tornado.gen.coroutine
    def get(self):
        # Parsing pagination params
        page_num = self.get_argument("page_num", 1)
        page_size = self.get_argument("page_size", 10)
        try:
            page_num = int(page_num)
        except:
            logging.exception("Error while parsing page_num: {}".format(page_num))
            self.write_json({"result": False, "errors": "invalid page_num"}, status_code=400)
            return

        try:
            page_size = int(page_size)
        except:
            logging.exception("Error while parsing page_size: {}".format(page_size))
            self.write_json({"result": False, "errors": "invalid page_size"}, status_code=400)
            return

        # Parsing user_id param
        user_id = self.get_argument("user_id", None)
        if user_id is not None:
            try:
                user_id = int(user_id)
            except:
                self.write_json({"result": False, "errors": "invalid user_id"}, status_code=400)
                return

        # Building select statement
        select_stmt = "SELECT * FROM listings"
        # Adding user_id filter clause if param is specified
        if user_id is not None:
            select_stmt += " WHERE user_id=?"
        # Order by and pagination
        limit = page_size
        offset = (page_num - 1) * page_size
        select_stmt += " ORDER BY created_at DESC LIMIT ? OFFSET ?"

        # Fetching listings from db
        if user_id is not None:
            args = (user_id, limit, offset)
        else:
            args = (limit, offset)
        cursor = self.application.db.cursor()
        results = cursor.execute(select_stmt, args)

        listings = []
        for row in results:
            fields = ["id", "user_id", "listing_type", "price", "created_at", "updated_at"]
            listing = {
                field: row[field] for field in fields
            }
            listings.append(listing)

        self.write_json({"result": True, "listings": listings})

    @tornado.gen.coroutine
    def post(self):
        # Collecting required params
        user_id = self.get_argument("user_id")
        listing_type = self.get_argument("listing_type")
        price = self.get_argument("price")

        # Validating inputs
        errors = []
        user_id_val = self._validate_user_id(user_id, errors)
        listing_type_val = self._validate_listing_type(listing_type, errors)
        price_val = self._validate_price(price, errors)
        time_now = int(time.time() * 1e6) # Converting current time to microseconds

        # End if we have any validation errors
        if len(errors) > 0:
            self.write_json({"result": False, "errors": errors}, status_code=400)
            return

        # Proceed to store the listing in our db
        cursor = self.application.db.cursor()
        cursor.execute(
            "INSERT INTO 'listings' "
            + "('user_id', 'listing_type', 'price', 'created_at', 'updated_at') "
            + "VALUES (?, ?, ?, ?, ?)",
            (user_id_val, listing_type_val, price_val, time_now, time_now)
        )
        self.application.db.commit()

        # Publishing the listing to NATS
        self.application.nats_conn.publish("listings.created", json.dumps({
            "user_id": user_id_val,
            "listing_type": listing_type_val,
            "price": price_val,
            "created_at": time_now,
            "updated_at": time_now
        }))

        # Error out if we fail to retrieve the newly created listing
        if cursor.lastrowid is None:
            self.write_json({"result": False, "errors": ["Error while adding listing to db"]}, status_code=500)
            return

        listing = dict(
            id=cursor.lastrowid,
            user_id=user_id_val,
            listing_type=listing_type_val,
            price=price_val,
            created_at=time_now,
            updated_at=time_now
        )

        self.write_json({"result": True, "listing": listing})

    def _validate_user_id(self, user_id, errors):
        try:
            user_id = int(user_id)
            return user_id
        except Exception as e:
            logging.exception("Error while converting user_id to int: {}".format(user_id))
            errors.append("invalid user_id")
            return None

    def _validate_listing_type(self, listing_type, errors):
        if listing_type not in {"rent", "sale"}:
            errors.append("invalid listing_type. Supported values: 'rent', 'sale'")
            return None
        else:
            return listing_type

    def _validate_price(self, price, errors):
        # Convert string to int
        try:
            price = int(price)
        except Exception as e:
            logging.exception("Error while converting price to int: {}".format(price))
            errors.append("invalid price. Must be an integer")
            return None

        if price < 1:
            errors.append("price must be greater than 0")
            return None
        else:
            return price

# /listings/ping
class PingHandler(tornado.web.RequestHandler):
    @tornado.gen.coroutine
    def get(self):
        self.write("pong!")

def signal_handler(sig, _frame, app, server):
    """Handle shutdown signals gracefully"""
    logging.info(f'Received signal {sig}, initiating graceful shutdown...')
    
    # Stop accepting new connections
    server.stop()
    
    # Close database connections
    app.close_db()
    
    # Stop the IOLoop after a short delay to allow pending operations to complete
    io_loop = tornado.ioloop.IOLoop.current()
    
    def shutdown():
        logging.info('Shutting down...')
        io_loop.stop()
    
    # Give it 1 second to finish pending requests
    io_loop.call_later(1, shutdown)

def make_app(options):
    app = App([
        (r"/listings/ping", PingHandler),
        (r"/listings", ListingsHandler),
    ], debug=options.debug)
    app.init_nats(options.nats_url)

    return app

if __name__ == "__main__":
    # Define settings/options for the web app
    # Specify the port number to start the web app on (default value is port 6000)
    tornado.options.define("port", default=6000)
    # Specify whether the app should run in debug mode
    # Debug mode restarts the app automatically on file changes
    tornado.options.define("debug", default=True)

    # Specify the NATS connection URL
    tornado.options.define("nats_url", default="nats://localhost:4222")

    # Read settings/options from command line
    tornado.options.parse_command_line()

    # Access the settings defined
    options = tornado.options.options

    # Create web app
    app = make_app(options)
    server = app.listen(options.port)
    
    # Register signal handlers for graceful shutdown
    signal.signal(signal.SIGTERM, lambda sig, _frame: signal_handler(sig, _frame, app, server))
    signal.signal(signal.SIGINT, lambda sig, _frame: signal_handler(sig, _frame, app, server))
    
    logging.info("Starting listing service. PORT: {}, DEBUG: {}".format(options.port, options.debug))
    logging.info("Press Ctrl+C to gracefully shutdown the server")

    try:
        # Start event loop
        tornado.ioloop.IOLoop.current().start()
    except KeyboardInterrupt:
        logging.info("Received KeyboardInterrupt, shutting down...")
    finally:
        # Ensure database is closed on exit
        app.close_db()
        logging.info("Listing service shutdown complete")
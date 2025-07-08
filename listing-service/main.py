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
        self.nats_conn = None
        self.nats_event = "listing.created"

    async def init_nats(self, connection_url):
        try:
            self.nats_conn = await nats.connect(connection_url)
            logging.info(f"Connected to NATS at {connection_url}")
        except Exception as e:
            logging.error(f"Failed to connect to NATS: {e}")
            raise

    async def close_nats(self):
        if self.nats_conn:
            await self.nats_conn.close()
            logging.info("NATS connection closed")

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

        if self.application.nats_conn:
            try:
                yield self.application.nats_conn.publish(self.application.nats_event, json.dumps({
                    "id": cursor.lastrowid,
                    "user_id": user_id_val,
                    "listing_type": listing_type_val,
                    "price": price_val,
                    "created_at": time_now,
                    "updated_at": time_now
                }).encode())
                logging.info("Published listing to NATS")
            except Exception as e:
                logging.error(f"Failed to publish to NATS: {e}")
        
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

async def make_app(options):
    app = App([
        (r"/listings/ping", PingHandler),
        (r"/listings", ListingsHandler),
    ], debug=options.debug)
    
    # Initialize NATS connection
    await app.init_nats(options.nats_url)
    
    return app

if __name__ == "__main__":
    # Define settings/options for the web app
    tornado.options.define("port", default=6000)
    tornado.options.define("debug", default=True)
    tornado.options.define("nats_url", default="nats://localhost:4222")

    # Read settings/options from command line
    tornado.options.parse_command_line()
    options = tornado.options.options

    async def setup_app():
        """Setup the app asynchronously"""
        return await make_app(options)
    
    try:
        app = tornado.ioloop.IOLoop.current().run_sync(setup_app)
        server = app.listen(options.port)
        
        # Register signal handlers for graceful shutdown
        def signal_handler(sig, _frame):
            async def shutdown():
                logging.info(f'Received signal {sig}, shutting down...')
                await app.close_nats()
                app.close_db()
                server.stop()
                tornado.ioloop.IOLoop.current().stop()
            
            tornado.ioloop.IOLoop.current().add_callback(shutdown)
        
        signal.signal(signal.SIGTERM, signal_handler)
        signal.signal(signal.SIGINT, signal_handler)
        
        logging.info("Starting listing service. PORT: {}, DEBUG: {}".format(options.port, options.debug))
        logging.info("Press Ctrl+C to gracefully shutdown the server")

        # Start event loop (this keeps the server running)
        tornado.ioloop.IOLoop.current().start()
        
    except KeyboardInterrupt:
        logging.info("Received KeyboardInterrupt, shutting down...")
    except Exception as e:
        logging.error(f"Error starting application: {e}")
        raise
    finally:
        # Cleanup
        if 'app' in locals() and app.nats_conn:
            tornado.ioloop.IOLoop.current().run_sync(app.close_nats)
        if 'app' in locals():
            app.close_db()
        logging.info("Listing service shutdown complete")
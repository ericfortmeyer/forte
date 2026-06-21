#!/usr/bin/env python3

from http.server import HTTPServer, BaseHTTPRequestHandler
import json

class HealthCheckHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/':
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response = json.dumps({'status':'ok','app':'forte-example-python'})
            self.wfile.write(response.encode())
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args):
        pass # Suppress default logging

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8000), HealthCheckHandler)
    print('Server listening on 0.0.0.0:8000')
    server.serve_forever()

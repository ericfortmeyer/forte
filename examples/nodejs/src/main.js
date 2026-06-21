#!/usr/bin/env node

const http = require('node:http');

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
        status: 'ok',
        app: 'example-nodejs-app',
    }));
});

server.listen(8000);

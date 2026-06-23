#!/usr/bin/env node

const http = require('node:http');
const process = require('node:process')

const server = http.createServer(async (req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
        status: 'ok',
        app: 'example-nodejs-app',
    }));
})

const shutdown = () => {
    console.log('Shutting down gracefully...');
    server.close();
    process.exit(0)
}

server.listen(8000);
process.on('SIGTERM', shutdown);

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;

final class SimpleServer {
    private static final int PORT = 8000;
    private static final String ROOT_PATH = "/";
    private static final String RESPONSE_BODY = "{\"status\":\"ok\",\"app\":\"example-java-app\"}";

    public static void main(String[] args) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress(PORT), 0);

        server.createContext(ROOT_PATH, SimpleServer::handle);
        server.setExecutor(null);
        server.start();

        System.out.printf("Server listening on port %d%n", PORT);
    }

    public static void handle(HttpExchange exchange) throws IOException {
        byte[] bytes = RESPONSE_BODY.getBytes(StandardCharsets.UTF_8);

        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(200, bytes.length);

        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}

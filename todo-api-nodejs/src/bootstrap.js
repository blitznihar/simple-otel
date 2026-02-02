import "dotenv/config";

// Start OTel first (must be first import after dotenv)
await import("./otel.js");

// Only after OTel is ready, load the app
await import("./index.js");
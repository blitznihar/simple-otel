import "dotenv/config";            // loads .env
import { startOtel } from "./otel.js";
import express from "express";
import { router } from "./routes/todos.js";

async function main() {
  await startOtel();

  const app = express();
  app.use(express.json());
  app.use(router);

  const port = process.env.PORT || "9100";

  const server = app.listen(Number(port), () => {
    console.log(`todo-api-nodejs listening on :${port}`);
  });

  const shutdown = async () => {
    console.log("shutting down...");
    server.close(() => process.exit(0));
  };

  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
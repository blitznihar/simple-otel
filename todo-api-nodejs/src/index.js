import express from "express";
import { router } from "./routes/todos.js";

const app = express();
app.use(express.json());
app.use(router);

const port = Number(process.env.PORT || "9100");

const server = app.listen(port, () => {
  console.log(`todo-api-nodejs listening on :${port}`);
});

const shutdown = () => {
  console.log("shutting down...");
  server.close(() => process.exit(0));
};

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
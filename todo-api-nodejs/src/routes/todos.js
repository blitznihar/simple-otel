import express from "express";
import {
  goGetTodos,
  goGetTodoById,
  goCreateTodo,
  goUpdateTodo,
} from "../clients/goTodosClient.js";
import { trace } from "@opentelemetry/api";

export const router = express.Router();

const tracer = trace.getTracer("manual");

router.get("/ping", (req, res) => {
  tracer.startActiveSpan("ping-span", (span) => {
    span.end();
    res.json({ ok: true });
  });
});

router.get("/health", (req, res) => {
  res.json({ status: "ok" });
});

router.get("/todos", async (req, res) => {
  try {
    const todos = await goGetTodos();
    res.json(todos);
  } catch (e) {
    res.status(502).json({ error: String(e.message || e) });
  }
});

router.get("/todos/:id", async (req, res) => {
  try {
    const todo = await goGetTodoById(req.params.id);
    res.json(todo);
  } catch (e) {
    res.status(502).json({ error: String(e.message || e) });
  }
});

router.post("/todos", async (req, res) => {
  try {
    const out = await goCreateTodo(req.body);
    res.status(201).json(out);
  } catch (e) {
    res.status(502).json({ error: String(e.message || e) });
  }
});

router.put("/todos/:id", async (req, res) => {
  try {
    const out = await goUpdateTodo(req.params.id, req.body);
    res.json(out);
  } catch (e) {
    res.status(502).json({ error: String(e.message || e) });
  }
});
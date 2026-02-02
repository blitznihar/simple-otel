import fetch from "node-fetch";

function baseUrl() {
  const v = process.env.TODO_GO_BASE_URL;
  if (!v) throw new Error("TODO_GO_BASE_URL is required");
  return v.replace(/\/$/, "");
}

export async function goGetTodos() {
  const res = await fetch(`${baseUrl()}/todos`);
  if (!res.ok) throw new Error(`go api error: ${res.status}`);
  return res.json();
}

export async function goGetTodoById(id) {
  const res = await fetch(`${baseUrl()}/todos/${id}`);
  if (!res.ok) throw new Error(`go api error: ${res.status}`);
  return res.json();
}

export async function goCreateTodo(payload) {
  const res = await fetch(`${baseUrl()}/todos`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw new Error(`go api error: ${res.status}`);
  return res.json();
}

export async function goUpdateTodo(id, payload) {
  const res = await fetch(`${baseUrl()}/todos/${id}`, {
    method: "PUT",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw new Error(`go api error: ${res.status}`);
  return res.json();
}
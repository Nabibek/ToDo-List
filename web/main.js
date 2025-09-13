// app.js - минимальный фронтенд для твоего Go API
// Предположения по эндпоинтам:
// GET  /todos?status=...&order=...   (получить список)
// POST /todo                          (создать задачу)
// GET  /todo/{id}                     (получить задачу по id)
// PUT  /todo/{id}                     (обновить задачу - передать весь объект)
// DELETE /todo/{id}                   (удалить задачу)

// Если у тебя другие пути — поправь API_BASE / ENDPOINTS ниже.
const API_BASE = "http://localhost:8000";

const selectors = {
  form: document.getElementById("todo-form"),
  inputTodo: document.getElementById("input-todo"),
  inputMessage: document.getElementById("input-message"),
  inputDeadline: document.getElementById("input-deadline"),
  list: document.getElementById("todos-list"),
  empty: document.getElementById("empty"),
  filterStatus: document.getElementById("filter-status"),
  filterOrder: document.getElementById("filter-order"),
  btnRefresh: document.getElementById("btn-refresh"),
};

async function fetchTodos(status = "", order = "") {
  const params = new URLSearchParams();
  if (status) params.set("status", status);
  if (order) params.set("order", order);
  const url = `${API_BASE}/todos${params.toString() ? "?" + params.toString() : ""}`;
  const res = await fetch(url);
  if (!res.ok) throw new Error("Fetch todos failed: " + res.status);
  return res.json();
}

async function fetchTodoById(id) {
  const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`);
  if (!res.ok) throw new Error("Fetch todo by id failed");
  return res.json();
}

async function createTodo(payload) {
  const res = await fetch(`${API_BASE}/todo`, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    const txt = await res.text();
    throw new Error("Create failed: " + res.status + " " + txt);
  }
  return res.json();
}

async function updateTodo(id, payload) {
  const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`, {
    method: "PUT",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    const txt = await res.text();
    throw new Error("Update failed: " + res.status + " " + txt);
  }
  return res;
}

async function deleteTodo(id) {
  const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    const txt = await res.text();
    throw new Error("Delete failed: " + res.status + " " + txt);
  }
  return res;
}

function isoToLocalInput(value) {
  if (!value) return "";
  const d = new Date(value);
  // produce something like 'YYYY-MM-DDTHH:MM' which fits datetime-local
  const pad = (n)=> String(n).padStart(2,"0");
  const y = d.getFullYear();
  const mo = pad(d.getMonth()+1);
  const day = pad(d.getDate());
  const h = pad(d.getHours());
  const m = pad(d.getMinutes());
  return `${y}-${mo}-${day}T${h}:${m}`;
}

function niceDate(value) {
  if (!value) return "";
  const d = new Date(value);
  return d.toLocaleString();
}

// Render list
function renderTodos(todos) {
  selectors.list.innerHTML = "";
  if (!todos || todos.length === 0) {
    selectors.empty.style.display = "block";
    return;
  }
  selectors.empty.style.display = "none";

  for (const t of todos) {
    const li = document.createElement("li");
    li.className = "todo-item";

    const left = document.createElement("div");
    left.className = "todo-left";

    const cb = document.createElement("input");
    cb.type = "checkbox";
    cb.checked = !!t.complete;
    cb.addEventListener("change", async () => {
      try {
        // fetch full todo to send full object (backend expects full object)
        const todo = await fetchTodoById(t.id);
        todo.complete = cb.checked;
        // if marking complete and completedAt empty, set now
        if (cb.checked && (!todo.completedAt || todo.completedAt === "0001-01-01T00:00:00Z")) {
          todo.completedAt = new Date().toISOString();
        }
        // if unchecking, clear completedAt
        if (!cb.checked) {
          todo.completedAt = null;
        }
        await updateTodo(t.id, todo);
        await reloadList();
      } catch (e) {
        alert("Ошибка при обновлении: " + e.message);
        console.error(e);
      }
    });

    const title = document.createElement("span");
    title.className = "title" + (t.complete ? " done" : "");
    title.textContent = t.todo || "(без названия)";

    const meta = document.createElement("div");
    meta.style.fontSize = "12px";
    meta.style.color = "#6b7280";
    const parts = [];
    if (t.deadline) parts.push("Дедлайн: " + niceDate(t.deadline));
    if (t.createdAt) parts.push("Создано: " + niceDate(t.createdAt));
    meta.textContent = parts.join(" • ");

    left.appendChild(cb);
    const textBlock = document.createElement("div");
    textBlock.appendChild(title);
    if (t.message) {
      const msg = document.createElement("div");
      msg.style.fontSize="13px";
      msg.style.color="#374151";
      msg.textContent = t.message;
      textBlock.appendChild(msg);
    }
    textBlock.appendChild(meta);
    left.appendChild(textBlock);

    const right = document.createElement("div");
    right.className = "todo-right";

    const btnDel = document.createElement("button");
    btnDel.className = "btn btn-danger btn-small";
    btnDel.textContent = "Удалить";
    btnDel.addEventListener("click", async () => {
      if (!confirm("Удалить задачу?")) return;
      try {
        await deleteTodo(t.id);
        await reloadList();
      } catch (e) {
        alert("Ошибка при удалении: " + e.message);
        console.error(e);
      }
    });

    right.appendChild(btnDel);

    li.appendChild(left);
    li.appendChild(right);

    selectors.list.appendChild(li);
  }
}

// Load & render according to filters
async function reloadList() {
  const status = selectors.filterStatus.value; // "", "complete", "incomplete"
  const order = selectors.filterOrder.value; // "asc"|"desc"
  try {
    const todos = await fetchTodos(status, order);
    renderTodos(todos);
  } catch (e) {
    console.error(e);
    alert("Не удалось загрузить список задач: " + e.message);
  }
}

selectors.form.addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const title = selectors.inputTodo.value.trim();
  if (!title) {
    alert("Введите название задачи");
    return;
  }
  const payload = {
    todo: title,
    message: selectors.inputMessage.value.trim(),
    deadline: selectors.inputDeadline.value ? new Date(selectors.inputDeadline.value).toISOString() : null,
  };
  try {
    await createTodo(payload);
    selectors.inputTodo.value = "";
    selectors.inputMessage.value = "";
    selectors.inputDeadline.value = "";
    await reloadList();
  } catch (e) {
    console.error(e);
    alert("Ошибка при создании: " + e.message);
  }
});

// filters & refresh
selectors.filterStatus.addEventListener("change", reloadList);
selectors.filterOrder.addEventListener("change", reloadList);
selectors.btnRefresh.addEventListener("click", reloadList);

// initial load
reloadList();

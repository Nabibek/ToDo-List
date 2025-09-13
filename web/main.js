const API = "http://localhost:8080"; // адрес бэка

async function loadTodos() {
  const status = document.getElementById("statusFilter").value;
  const order = document.getElementById("orderFilter").value;

  let url = `${API}/todos?order=${order}`;
  if (status) url += `&status=${status}`;

  const res = await fetch(url);
  const todos = await res.json();

  const list = document.getElementById("todoList");
  list.innerHTML = "";
  todos.forEach(todo => {
    const div = document.createElement("div");
    div.className = "todo-item " + (todo.complete ? "completed" : "");
    div.innerHTML = `
      <b>${todo.todo}</b> [${todo.priority}]<br>
      ${todo.message}<br>
      <small>Создано: ${new Date(todo.created_at).toLocaleString()}</small>
      <div class="controls">
        <button onclick="deleteTodo('${todo.id}')">Удалить</button>
        ${!todo.complete 
          ? `<button onclick="completeTodo('${todo.id}')">Выполнить</button>`
          : `<button onclick="uncompleteTodo('${todo.id}')">Вернуть</button>`}
      </div>
    `;
    list.appendChild(div);
  });
}

async function addTodo() {
  const todo = document.getElementById("todo").value;
  const message = document.getElementById("message").value;
  const priority = document.getElementById("priority").value;

  if (!todo || !message) {
    alert("Заполни название и описание");
    return;
  }

  await fetch(`${API}/todo`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ todo, message, priority })
  });

  document.getElementById("todo").value = "";
  document.getElementById("message").value = "";
  loadTodos();
}

async function deleteTodo(id) {
  await fetch(`${API}/todo/${id}`, { method: "DELETE" });
  loadTodos();
}

async function completeTodo(id) {
  await fetch(`${API}/todo/complete/${id}`, { method: "POST" });
  loadTodos();
}

// опционально: "развыполнить" → можно просто апдейтить через PUT
async function uncompleteTodo(id) {
  await fetch(`${API}/todo/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ complete: false })
  });
  loadTodos();
}

loadTodos();

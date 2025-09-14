// app.js - Полный фронтенд для ToDo Desktop приложения
const API_BASE = "http://localhost:8000/api";

const selectors = {
  form: document.getElementById("todo-form"),
  inputTodo: document.getElementById("input-todo"),
  inputMessage: document.getElementById("input-message"),
  inputDeadline: document.getElementById("input-deadline"),
  inputPriority: document.getElementById("input-priority"),
  list: document.getElementById("todos-list"),
  tasksActive: document.getElementById("todos-active"),
  tasksCompleted: document.getElementById("todos-completed"),
  empty: document.getElementById("empty"),
  filterStatus: document.getElementById("filter-status"),
  filterOrder: document.getElementById("filter-order"),
  filterPeriod: document.getElementById("filter-period"),
  btnRefresh: document.getElementById("btn-refresh"),
  themeToggle: document.getElementById("theme-toggle"),
  tasksCount: document.getElementById("tasks-count"),
  tasksStats: document.getElementById("tasks-stats"),
  deleteModal: document.getElementById("delete-modal"),
  modalCancel: document.getElementById("modal-cancel"),
  modalConfirm: document.getElementById("modal-confirm"),
  tasksContainer: document.getElementById("tasks-container")
};

// Глобальные переменные
let todos = [];
let todoToDelete = null;

// Инициализация приложения
function initApp() {
  initTheme();
  initEventListeners();
  loadTodos();
}

// Инициализация темы
function initTheme() {
  const savedTheme = localStorage.getItem('theme') || 'light';
  document.documentElement.setAttribute('data-theme', savedTheme);
  updateThemeIcon(savedTheme);
}

function toggleTheme() {
  const currentTheme = document.documentElement.getAttribute('data-theme');
  const newTheme = currentTheme === 'light' ? 'dark' : 'light';
  
  document.documentElement.setAttribute('data-theme', newTheme);
  localStorage.setItem('theme', newTheme);
  updateThemeIcon(newTheme);
}

function updateThemeIcon(theme) {
  const icon = selectors.themeToggle.querySelector('i');
  icon.className = theme === 'light' ? 'fas fa-moon' : 'fas fa-sun';
}

// Инициализация обработчиков событий
function initEventListeners() {
  
  // Форма добавления задачи
  selectors.form.addEventListener("submit", handleCreateTodo);
  
  // Фильтры
  selectors.filterStatus.addEventListener("change", applyFilters);
  selectors.filterOrder.addEventListener("change", applyFilters);
  selectors.filterPeriod.addEventListener("change", applyFilters);
  selectors.btnRefresh.addEventListener("click", loadTodos);
  loadTodos();
  
  // Тема
  selectors.themeToggle.addEventListener("click", toggleTheme);
  
  // Модальное окно удаления
  selectors.modalCancel.addEventListener("click", hideDeleteModal);
  selectors.modalConfirm.addEventListener("click", confirmDelete);
  selectors.deleteModal.addEventListener("click", (e) => {
    if (e.target === selectors.deleteModal) hideDeleteModal();
  });
}

// API функции
async function fetchTodos() {
  try {
    const status = selectors.filterStatus.value;
    const order = selectors.filterOrder.value;
    const period = selectors.filterPeriod.value;
    
    const params = new URLSearchParams();
    if (status) params.set("status", status);
    if (period) params.set("period", period);
    
    // Правильная обработка сортировки
    if (order === 'priority') {
      params.set("orderBy", "priority");
      params.set("orderDir", "asc"); // Для приоритета всегда asc (высокий -> низкий)
    } else {
      params.set("orderBy", "created_at");
      params.set("orderDir", order);
    }
    
    const url = `${API_BASE}/todos${params.toString() ? "?" + params.toString() : ""}`;
    console.log('Fetching URL:', url);
    
    const res = await fetch(url);
    
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
    
    return await res.json();
  } catch (error) {
    console.error("Fetch todos failed:", error);
    showError("Не удалось загрузить задачи");
    return [];
  }
}
async function fetchTodoById(id) {
  try {
    const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`);
    if (!res.ok) throw new Error("Fetch todo by id failed");
    return await res.json();
  } catch (error) {
    console.error("Fetch todo failed:", error);
    throw error;
  }
}

async function createTodo(payload) {
  try {
    const res = await fetch(`${API_BASE}/todo`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(payload),
    });
    
    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Create failed: ${res.status} ${errorText}`);
    }
    
    return await res.json();
  } catch (error) {
    console.error("Create todo failed:", error);
    throw error;
  }
}

async function updateTodo(id, payload) {
  try {
    const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`, {
      method: "PUT",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(payload),
    });
    
    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Update failed: ${res.status} ${errorText}`);
    }
    
    return res;
  } catch (error) {
    console.error("Update todo failed:", error);
    throw error;
  }
}

async function deleteTodo(id) {
  try {
    const res = await fetch(`${API_BASE}/todo/${encodeURIComponent(id)}`, {
      method: "DELETE",
    });
    
    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Delete failed: ${res.status} ${errorText}`);
    }
    
    return res;
  } catch (error) {
    console.error("Delete todo failed:", error);
    throw error;
  }
}

// Обработчики событий
async function handleCreateTodo(ev) {
  ev.preventDefault();
  
  const title = selectors.inputTodo.value.trim();
  if (!title) {
    showError("Введите название задачи");
    return;
  }
  
  const payload = {
    todo: title,
    message: selectors.inputMessage.value.trim(),
    deadline: selectors.inputDeadline.value ? new Date(selectors.inputDeadline.value).toISOString() : null,
    priority: selectors.inputPriority.value,
  };
  
  try {
    await createTodo(payload);
    
    // Сброс формы
    selectors.inputTodo.value = "";
    selectors.inputMessage.value = "";
    selectors.inputDeadline.value = "";
    
    // Перезагрузка списка
    await loadTodos();
    
    showSuccess("Задача успешно добавлена");
  } catch (error) {
    showError("Ошибка при создании задачи: " + error.message);
  }
}

async function handleToggleComplete(todo) {
  try {
    const updatedTodo = { ...todo };
    updatedTodo.complete = !updatedTodo.complete;
    
    if (updatedTodo.complete) {
      updatedTodo.completedAt = new Date().toISOString();
    } else {
      updatedTodo.completedAt = null;
    }
    
    await updateTodo(todo.id, updatedTodo);
    await loadTodos();
    
    showSuccess(`Задача отмечена как ${updatedTodo.complete ? 'выполненная' : 'активная'}`);
  } catch (error) {
    showError("Ошибка при обновлении задачи: " + error.message);
  }
}

async function handleDeleteClick(todoId) {
  showDeleteModal(todoId);
}

// Модальное окно удаления
function showDeleteModal(todoId) {
  todoToDelete = todoId;
  selectors.deleteModal.style.display = 'block';
}

function hideDeleteModal() {
  todoToDelete = null;
  selectors.deleteModal.style.display = 'none';
}

async function confirmDelete() {
  if (!todoToDelete) return;
  
  try {
    await deleteTodo(todoToDelete);
    await loadTodos();
    hideDeleteModal();
    showSuccess("Задача успешно удалена");
  } catch (error) {
    showError("Ошибка при удалении задачи: " + error.message);
    hideDeleteModal();
  }
}

// Загрузка и отображение задач
async function loadTodos() {
  try {
    selectors.btnRefresh.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Загрузка...';
    
    todos = await fetchTodos();
    renderTodos(todos);
    
    selectors.btnRefresh.innerHTML = '<i class="fas fa-sync"></i> Обновить';
  } catch (error) {
    selectors.btnRefresh.innerHTML = '<i class="fas fa-sync"></i> Обновить';
    console.error("Load todos failed:", error);
  }
}

function applyFilters() {
  const statusFilter = selectors.filterStatus.value;
  const periodFilter = selectors.filterPeriod.value;
  const orderFilter = selectors.filterOrder.value;
  
  // Если выбраны сложные фильтры, перезагружаем с сервера
  if (periodFilter || statusFilter === 'overdue' || 
      orderFilter === 'priority' || orderFilter === 'asc') {
    loadTodos();
  } else {
    // Иначе фильтруем локально (только по статусу и desc сортировка)
    renderTodos(todos);
  }
}
function renderTodos(todos) {
  // Фильтрация по статусу
  let filteredTodos = [...todos];
  const statusFilter = selectors.filterStatus.value;
  
  if (statusFilter === 'active') {
    filteredTodos = filteredTodos.filter(t => !t.complete);
  } else if (statusFilter === 'completed') {
    filteredTodos = filteredTodos.filter(t => t.complete);
  }
  
  // Фильтрация по периоду
  const periodFilter = selectors.filterPeriod.value;
  if (periodFilter === 'today') {
    const today = new Date().toDateString();
    filteredTodos = filteredTodos.filter(t => 
      new Date(t.createdAt).toDateString() === today
    );
  } else if (periodFilter === 'week') {
    const weekAgo = new Date();
    weekAgo.setDate(weekAgo.getDate() - 7);
    filteredTodos = filteredTodos.filter(t => 
      new Date(t.createdAt) >= weekAgo
    );
  } else if (periodFilter === 'overdue') {
    const now = new Date();
    filteredTodos = filteredTodos.filter(t => 
      t.deadline && new Date(t.deadline) < now && !t.complete
    );
  }
  
  // Сортировка
  const orderFilter = selectors.filterOrder.value;
  if (orderFilter === 'asc') {
    filteredTodos.sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt));
  } else if (orderFilter === 'desc') {
    filteredTodos.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
  } else if (orderFilter === 'priority') {
    const priorityOrder = { high: 3, medium: 2, low: 1 };
    filteredTodos.sort((a, b) => priorityOrder[b.priority] - priorityOrder[a.priority]);
  }
  
  // Разделение на активные и выполненные
  const activeTodos = filteredTodos.filter(t => !t.complete);
  const completedTodos = filteredTodos.filter(t => t.complete);
  
  renderTodosSection(selectors.tasksActive, activeTodos, 'active');
  renderTodosSection(selectors.tasksCompleted, completedTodos, 'completed');
  
  updateTasksStats(activeTodos.length, completedTodos.length);
  toggleEmptyState(todos.length === 0);
}

function renderTodosSection(container, todos, section) {
  container.innerHTML = '';
  
  if (todos.length === 0) {
    const emptyMsg = document.createElement('li');
    emptyMsg.className = 'muted';
    emptyMsg.innerHTML = section === 'active' 
      ? '<i class="fas fa-check-circle"></i><p>Нет активных задач</p>' 
      : '<i class="fas fa-inbox"></i><p>Нет выполненных задач</p>';
    container.appendChild(emptyMsg);
    return;
  }
  
  todos.forEach(todo => {
    const li = createTodoElement(todo);
    container.appendChild(li);
  });
}

function createTodoElement(t) {
  const li = document.createElement('li');
  li.className = 'todo-item';
  li.dataset.priority = t.priority;
  li.dataset.id = t.id;

  const left = document.createElement('div');
  left.className = 'todo-left';

  // Чекбокс выполнения
  const cb = document.createElement('input');
  cb.type = 'checkbox';
  cb.checked = !!t.complete;
  cb.addEventListener('change', () => handleToggleComplete(t));

  // Контент задачи
  const content = document.createElement('div');
  content.className = 'todo-content';

  // Заголовок
  const title = document.createElement('span');
  title.className = 'todo-title' + (t.complete ? ' done' : '');
  title.textContent = t.todo || '(без названия)';

  // Описание
  const description = document.createElement('div');
  description.className = 'todo-description';
  description.textContent = t.message || '';

  // Мета-информация
  const meta = document.createElement('div');
  meta.className = 'todo-meta';

  // Приоритет
  if (t.priority) {
    const priority = document.createElement('span');
    priority.className = `priority-badge priority-${t.priority}`;
    priority.textContent = getPriorityLabel(t.priority);
    meta.appendChild(priority);
  }

  // Дедлайн
  if (t.deadline) {
    const deadline = document.createElement('span');
    deadline.className = 'deadline';
    deadline.innerHTML = `<i class="fas fa-clock"></i> ${formatDeadline(t.deadline)}`;
    meta.appendChild(deadline);
  }

  // Дата создания
  if (t.createdAt) {
    const createdAt = document.createElement('span');
    createdAt.className = 'created-at';
    createdAt.innerHTML = `<i class="fas fa-calendar"></i> ${formatDate(t.createdAt)}`;
    meta.appendChild(createdAt);
  }

  content.appendChild(title);
  if (t.message) content.appendChild(description);
  content.appendChild(meta);

  left.appendChild(cb);
  left.appendChild(content);

  // Правая часть с кнопками
  const right = document.createElement('div');
  right.className = 'todo-right';

  // Кнопка удаления
  const btnDel = document.createElement('button');
  btnDel.className = 'btn btn-danger btn-small';
  btnDel.innerHTML = '<i class="fas fa-trash"></i>';
  btnDel.title = 'Удалить задачу';
  btnDel.addEventListener('click', () => handleDeleteClick(t.id));

  right.appendChild(btnDel);

  li.appendChild(left);
  li.appendChild(right);

  // Добавляем класс просроченности
  if (t.deadline && new Date(t.deadline) < new Date() && !t.complete) {
    li.classList.add('overdue');
  }

  return li;
}

function updateTasksStats(activeCount, completedCount) {
  const total = activeCount + completedCount;
  
  selectors.tasksCount.textContent = `Список задач (${total})`;
  
  let statsText = '';
  if (total > 0) {
    statsText = `Активных: ${activeCount} • Выполненных: ${completedCount}`;
    
    if (activeCount > 0) {
      const completionRate = Math.round((completedCount / total) * 100);
      statsText += ` • Прогресс: ${completionRate}%`;
    }
  }
  
  selectors.tasksStats.textContent = statsText;
}

function toggleEmptyState(isEmpty) {
  if (isEmpty) {
    selectors.tasksContainer.style.display = 'none';
    selectors.empty.style.display = 'block';
  } else {
    selectors.tasksContainer.style.display = 'block';
    selectors.empty.style.display = 'none';
  }
}

// Вспомогательные функции
function getPriorityLabel(priority) {
  const labels = {
    low: '🟢 Низкий',
    medium: '🟡 Средний', 
    high: '🔴 Высокий'
  };
  return labels[priority] || priority;
}

function formatDate(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  });
}

function formatDeadline(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  const now = new Date();
  
  if (date.toDateString() === now.toDateString()) {
    return `Сегодня, ${date.toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit'
    })}`;
  }
  
  return date.toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

function showError(message) {
  // Можно реализовать toast-уведомления
  alert(`❌ ${message}`);
}

function showSuccess(message) {
  // Можно реализовать toast-уведомления  
  alert(`✅ ${message}`);
}

// Утилиты для работы с датами
function isoToLocalInput(value) {
  if (!value) return "";
  const d = new Date(value);
  const pad = (n) => String(n).padStart(2, "0");
  const y = d.getFullYear();
  const mo = pad(d.getMonth() + 1);
  const day = pad(d.getDate());
  const h = pad(d.getHours());
  const m = pad(d.getMinutes());
  return `${y}-${mo}-${day}T${h}:${m}`;
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', initApp);

// Глобальные функции для отладки
window.debug = {
  getTodos: () => todos,
  clearFilters: () => {
    selectors.filterStatus.value = '';
    selectors.filterOrder.value = 'desc';
    selectors.filterPeriod.value = '';
    loadTodos();
  }
};
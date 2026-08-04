const API_BASE = '';

// Загрузка задач
async function loadTasks() {
    const taskList = document.getElementById('taskList');
    const status = document.getElementById('filterStatus').value;

    taskList.innerHTML = '<div class="loading">Загрузка...</div>';

    try {
        let url = `${API_BASE}/tasks`;
        if (status) {
            url += `?status=${status}`;
        }

        const response = await fetch(url);
        if (!response.ok) throw new Error('Ошибка загрузки');

        const tasks = await response.json();

        if (tasks.length === 0) {
            taskList.innerHTML = '<div class="empty-state"><p>Нет задач для отображения</p></div>';
            return;
        }

        taskList.innerHTML = tasks.map(task => `
            <div class="task-card">
                <div class="task-header">
                    <span class="task-title">${escapeHtml(task.title)}</span>
                    <span class="task-id">#${task.id}</span>
                </div>
                ${task.description ? `<p class="task-description">${escapeHtml(task.description)}</p>` : ''}
                <div class="task-meta">
                    <span class="badge badge-status-${task.status}">${getStatusText(task.status)}</span>
                    <span class="badge badge-priority-${task.priority}">${getPriorityText(task.priority)}</span>
                </div>
                <div class="task-actions">
                    <button onclick="openEditModal(${task.id})">Редактировать</button>
                    <button class="btn-danger" onclick="deleteTask(${task.id})">Удалить</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        taskList.innerHTML = `<div class="empty-state"><p>Ошибка: ${error.message}</p></div>`;
    }
}

// Создание задачи
document.getElementById('createForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const errorDiv = document.getElementById('createFormError');
    const successDiv = document.getElementById('createFormSuccess');

    const task = {
        title: document.getElementById('title').value.trim(),
        description: document.getElementById('description').value.trim(),
        status: document.getElementById('status').value,
        priority: document.getElementById('priority').value
    };

    try {
        const response = await fetch(`${API_BASE}/tasks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(task)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Ошибка создания');
        }

        successDiv.textContent = 'Задача успешно создана!';
        successDiv.style.display = 'block';
        errorDiv.style.display = 'none';

        document.getElementById('createForm').reset();
        loadTasks();

        setTimeout(() => successDiv.style.display = 'none', 3000);
    } catch (error) {
        errorDiv.textContent = error.message;
        errorDiv.style.display = 'block';
        successDiv.style.display = 'none';
    }
});

// Удаление задачи
async function deleteTask(id) {
    if (!confirm('Вы уверены, что хотите удалить эту задачу?')) return;

    try {
        const response = await fetch(`${API_BASE}/tasks/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Ошибка удаления');
        }

        loadTasks();
    } catch (error) {
        alert(`Ошибка: ${error.message}`);
    }
}

// Открытие модального окна редактирования
async function openEditModal(id) {
    try {
        const response = await fetch(`${API_BASE}/tasks/${id}`);
        if (!response.ok) throw new Error('Ошибка загрузки');

        const task = await response.json();

        document.getElementById('editId').value = task.id;
        document.getElementById('editTitle').value = task.title;
        document.getElementById('editDescription').value = task.description || '';
        document.getElementById('editStatus').value = task.status;
        document.getElementById('editPriority').value = task.priority;

        document.getElementById('editModal').classList.add('active');
        document.getElementById('editFormError').style.display = 'none';
    } catch (error) {
        alert(`Ошибка: ${error.message}`);
    }
}

// Закрытие модального окна
function closeModal() {
    document.getElementById('editModal').classList.remove('active');
}

// Сохранение изменений
document.getElementById('editForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const errorDiv = document.getElementById('editFormError');
    const id = document.getElementById('editId').value;

    const task = {
        title: document.getElementById('editTitle').value.trim(),
        description: document.getElementById('editDescription').value.trim(),
        status: document.getElementById('editStatus').value,
        priority: document.getElementById('editPriority').value
    };

    try {
        const response = await fetch(`${API_BASE}/tasks/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(task)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Ошибка обновления');
        }

        closeModal();
        loadTasks();
    } catch (error) {
        errorDiv.textContent = error.message;
        errorDiv.style.display = 'block';
    }
});

// Очистка фильтра
function clearFilter() {
    document.getElementById('filterStatus').value = '';
    loadTasks();
}

// Вспомогательные функции
function getStatusText(status) {
    const map = {
        'pending': 'Ожидает',
        'in_progress': 'В работе',
        'done': 'Завершено'
    };
    return map[status] || status;
}

function getPriorityText(priority) {
    const map = {
        'low': 'Низкий',
        'medium': 'Средний',
        'high': 'Высокий'
    };
    return map[priority] || priority;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Закрытие модального окна по клику вне его
document.getElementById('editModal').addEventListener('click', (e) => {
    if (e.target.id === 'editModal') closeModal();
});

// Загрузка задач при старте
loadTasks();
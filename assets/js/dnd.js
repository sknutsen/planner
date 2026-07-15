let draggedTaskId = null;
let sourceDate = null;

function clearDragOver() {
  document.querySelectorAll(".day-preview__tasks--drag-over").forEach((el) => {
    el.classList.remove("day-preview__tasks--drag-over");
  });
}

function taskFromDragHandle(target) {
  const handle = target.closest(".task-preview__drag-handle");
  if (!handle) {
    return null;
  }

  return handle.closest(".task-preview");
}

export function initTaskDragDrop() {
  document.addEventListener("dragstart", (e) => {
    const task = taskFromDragHandle(e.target);
    if (!task) {
      return;
    }

    draggedTaskId = task.dataset.taskId;
    sourceDate = task.dataset.date;
    task.classList.add("task-preview--dragging");
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", draggedTaskId);
  });

  document.addEventListener("dragend", () => {
    document.querySelectorAll(".task-preview--dragging").forEach((el) => {
      el.classList.remove("task-preview--dragging");
    });
    clearDragOver();
    draggedTaskId = null;
    sourceDate = null;
  });

  document.addEventListener("dragover", (e) => {
    const dropZone = e.target.closest(".day-preview__tasks");
    if (!dropZone || !draggedTaskId) {
      return;
    }

    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    clearDragOver();
    dropZone.classList.add("day-preview__tasks--drag-over");
  });

  document.addEventListener("dragleave", (e) => {
    const dropZone = e.target.closest(".day-preview__tasks");
    if (!dropZone) {
      return;
    }

    const related = e.relatedTarget;
    if (!related || !dropZone.contains(related)) {
      dropZone.classList.remove("day-preview__tasks--drag-over");
    }
  });

  document.addEventListener("drop", async (e) => {
    const dropZone = e.target.closest(".day-preview__tasks");
    if (!dropZone || !draggedTaskId) {
      return;
    }

    e.preventDefault();
    clearDragOver();

    const targetDate = dropZone.dataset.date;
    if (!targetDate || targetDate === sourceDate) {
      return;
    }

    try {
      const response = await fetch(`/task/${draggedTaskId}/reschedule`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ date: targetDate }),
      });

      if (response.ok && typeof htmx !== "undefined") {
        htmx.trigger(document.body, "updatedTask");
      }
    } catch (err) {
      console.error("Failed to reschedule task", err);
    }
  });
}

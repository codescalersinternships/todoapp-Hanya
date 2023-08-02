<template>
  <div class="container py-5 h-100">
    <div class="row d-flex justify-content-center align-items-center h-100">
      <div class="col col-lg-8 col-xl-6">
        <div class="card rounded-3">
          <div class="card-body p-4">
            <p class="mb-2">
              <span class="h2 me-2">Todo List</span>
            </p>
            <div class="input-group mb-3" style="flex-direction: row">
              <input
                type="text"
                class="form-control"
                aria-label="Default"
                aria-describedby="inputGroup-sizing-default"
                v-model="taskInput"
              />
              <div class="input-group-prepend">
                <button
                  @click="insertTask"
                  class="input-group-text btn btn-primary"
                  id="inputGroup-sizing-default"
                >
                  Add
                </button>
              </div>
            </div>
            <ul class="list-group rounded-0">
              <li
                class="list-group-item border-0.5 d-flex align-items-center ps-0"
                v-for="task in tasks"
                :key="task.id"
                :class="{ 'text-decoration-line-through': task.completed }"
              >
                <div class="d-flex col-1 justify-content-center">
                  <input
                    class="form-check-input me-3"
                    type="checkbox"
                    value=""
                    aria-label="..."
                    @click="toggleCompleted(task)"
                    :checked="task.completed"
                  />
                </div>
                <div
                  v-if="editBool && task.id === editTaskId"
                  class="col-12 d-flex justify-content-between"
                >
                  <input type="text" v-model="updatedTitle" />
                  <div>
                    <button
                      class="btn btn-primary btn-sm"
                      @click="saveTask(task)"
                    >
                      Save
                    </button>
                    <button
                      class="btn btn-secondary btn-sm ms-2"
                      @click="cancelEdit"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
                <div v-else class="col-11 d-flex justify-content-between">
                  <div class="d-flex col-12 justify-content-between">
                    <div>
                      <span>{{ task.task }}</span>
                    </div>
                    <div>
                      <button
                        class="btn btn-primary btn-sm"
                        @click="startEdit(task)"
                      >
                        <i class="bi bi-pencil-square"></i>
                      </button>
                      <button
                        class="btn btn-danger btn-sm"
                        @click="deleteTask(task.id)"
                      >
                        <i class="bi bi-trash3-fill"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap/dist/js/bootstrap.js";
import "bootstrap-icons/font/bootstrap-icons.css";
interface TodoItem {
  id: number;
  task: string;
  completed: boolean;
}
export default defineComponent({
  name: "TodoView",
  components: {},
  data() {
    return {
      tasks: [] as TodoItem[],
      taskInput: "",
      updatedTitle: "",
      editBool: false as boolean,
      editTaskId: null as number | null,
    };
  },
  beforeMount() {
    this.getAllTasks();
  },
  methods: {
    toggleCompleted(task: TodoItem) {
      task.completed = !task.completed;
      this.updateTask(task);
    },
    getAllTasks() {
      fetch("http://localhost:3000/")
        .then((response) => {
          if (!response.ok) {
            throw new Error("network response error");
          }
          return response.json();
        })
        .then((data) => {
          this.tasks = data;
        })
        .catch((error) => {
          console.error("error fetching tasks", error);
        });
    },
    insertTask() {
      const body = {
        task: this.taskInput,
      };
      if (this.taskInput == "") {
        throw new Error("task cannot be empty");
      }
      fetch("http://localhost:3000/todo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      }).then((response) => {
        if (response.ok) {
          this.getAllTasks();
          console.log("Data posted successfully!");
        } else {
          throw new Error("error posting data");
        }
      });
    },
    deleteTask(id: number) {
      fetch(`http://localhost:3000/todo/${id}`, {
        method: "DELETE",
      })
        .then((response) => {
          if (response.ok) {
            this.getAllTasks();
            console.log("task deleted successfully");
          } else {
            throw new Error("failed to delete task");
          }
        })
        .catch((error) => {
          console.error("error deleting task", error);
        });
    },
    updateTask(task: TodoItem) {
      fetch(`http://localhost:3000/todo/${task.id}`, {
        method: "PATCH",
        body: JSON.stringify({
          task: task.task,
          completed: task.completed,
        }),
      })
        .then((response) => {
          if (response.ok) {
            this.getAllTasks();
            console.log("task updated successfully!", response);
          } else {
            throw new Error("failed to update task");
          }
        })
        .catch((error) => {
          console.error("error updating task:", error);
        });
    },
    startEdit(task: TodoItem) {
      this.editBool = true;
      this.editTaskId = task.id;
      this.updatedTitle = task.task;
    },
    saveTask(task: TodoItem) {
      task.task = this.updatedTitle;
      this.updateTask(task);
      this.cancelEdit();
    },
    cancelEdit() {
      this.editBool = false;
      this.editTaskId = null;
      this.updatedTitle = "";
    },
  },
});
</script>

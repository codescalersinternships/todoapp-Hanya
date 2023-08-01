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
                  @click="insertTodo"
                  class="input-group-text btn btn-primary"
                  id="inputGroup-sizing-default"
                >
                  Add
                </button>
              </div>
            </div>
            <ul class="list-group rounded-0">
              <li
                class="list-group-item border-0 d-flex align-items-center ps-0"
                v-for="todo in todos"
                :key="todo.id"
              >
                <input
                  class="form-check-input me-3"
                  type="checkbox"
                  value=""
                  aria-label="..."
                  @click="toggleCompleted(todo)"
                />
                <div
                  class="d-flex align-items-center col-10 justify-content-between"
                >
                  <div>
                    <span>{{ todo.task }}</span>
                  </div>
                  <div>
                    <button
                      class="btn btn-primary btn-sm"
                      @click="updateTodo(todo)"
                    >
                      <i class="bi bi-pencil-square ms-auto"></i>
                    </button>
                    <button
                      class="btn btn-danger btn-sm"
                      @click="deleteTodo(todo.id)"
                    >
                      <i class="bi bi-trash3-fill ms-2"></i>
                    </button>
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
import Todo from "@/components/Todo.vue"; // @ is an alias to /src
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
      todos: [] as TodoItem[],
      taskInput: "",
    };
  },
  mounted() {
    this.getAllTodos();
  },
  methods: {
    toggleCompleted(todo: TodoItem) {
      todo.completed = !todo.completed;
      this.updateTodo(todo);
    },
    getAllTodos() {
      fetch("http://localhost:3000/")
        .then((response) => {
          if (!response.ok) {
            throw new Error("network response error");
          }
          return response.json();
        })
        .then((data) => {
          this.todos = data;
        })
        .catch((error) => {
          console.error("error fetching todos", error);
        });
    },
    insertTodo() {
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
          console.log("Data posted successfully!");
        } else {
          throw new Error("error posting data");
        }
      });
    },
    deleteTodo(id: number) {
      fetch(`http://localhost:3000/todo/${id}`, {
        method: "DELETE",
      })
        .then((response) => {
          if (response.ok) {
            console.log("todo deleted successfully");
          } else {
            throw new Error("failed to delete todo");
          }
        })
        .catch((error) => {
          console.error("error deleting todo", error);
        });
    },
    updateTodo(todo: TodoItem) {
      fetch(`http://localhost:3000/todo/${todo.id}`, {
        method: "PATCH",
        body: JSON.stringify({ task: todo.task, completed: todo.completed }),
      })
        .then((response) => {
          if (response.ok) {
            console.log("todo updated successfully!", response);
          } else {
            throw new Error("failed to update todo");
          }
        })
        .catch((error) => {
          console.error("error updating todo:", error);
        });
    },
  },
});
</script>

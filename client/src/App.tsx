import "./App.css";
import useSWR from "swr";
import { AddTodo } from "./components/AddTodo";
import { Box, List, ThemeIcon } from "@mantine/core";
import fetch from "unfetch";
import { CheckCircleFillIcon } from "@primer/octicons-react";

export interface Todo {
  id: number;
  title: string;
  body: string;
  isCompleted: boolean;
}

export const ENDPOINT: string = "http://localhost:4000";

const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((res) => {
    return res.json();
  });

function App() {
  const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);

  const toggleTodo = async (id: number) => {
    const updatedTodo = await fetch(`${ENDPOINT}/api/todos/${id}/toggle`, {
      method: "PATCH",
      headers: {
        "Content-type": "application/json",
      },
    }).then((res) => res.json());

    mutate(updatedTodo);
  };

  return (
    <Box
      sx={(theme) => ({
        padding: "2rem",
        width: "100%",
        maxWidth: "40rem",
        margin: "0 auto",
      })}
    >
      <List spacing="xs" size="sm" mb={12} center>
        {data?.map((todo) => {
          return (
            <List.Item
              onClick={() => toggleTodo(todo.id)}
              key={`todo__${todo.id}`}
              icon={
                todo.isCompleted ? (
                  <ThemeIcon color="teal" size={24} radius="xl">
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                ) : (
                  <ThemeIcon color="gray" size={24} radius="xl">
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                )
              }
            >
              <h1>{todo.title}</h1>
            </List.Item>
          );
        })}
      </List>
      <AddTodo mutate={mutate} />
    </Box>
  );
}

export default App;

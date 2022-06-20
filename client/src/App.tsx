import "./App.css";
import useSWR from "swr";
import { AddTodo } from "./components/AddTodo";
import { Box } from "@mantine/core";
import { ShowTodo } from "./components/ShowTodo";

export interface Todo {
  id: number;
  title: string;
  body: string;
  isCompleted: boolean;
}

export const ENDPOINT: string = "http://localhost:4000";

export const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((res) => {
    return res.json();
  });

function App() {
  const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);

  return (
    <Box
      sx={(theme) => ({
        padding: "2rem",
        width: "100%",
        maxWidth: "40rem",
        margin: "0 auto",
      })}
    >
      <ShowTodo data={data} mutate={mutate} />
      <AddTodo mutate={mutate} />
    </Box>
  );
}

export default App;

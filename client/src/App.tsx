import useSWR from "swr";
import { AddTodo } from "./components/AddTodo";
import { Box } from "@mantine/core";
import { ShowTodo } from "./components/ShowTodo";
import { Todo } from "./interfaces/todoInterface";
import { fetcher } from "./utilities/utilities";

export const App = () => {
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
};

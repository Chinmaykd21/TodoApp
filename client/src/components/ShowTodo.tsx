import { Grid, List, ThemeIcon } from "@mantine/core";
import { CheckCircleFillIcon } from "@primer/octicons-react";
import { KeyedMutator } from "swr";
import { Todo } from "../interfaces/todoInterface";
import { ENDPOINT } from "../utilities/utilities";
import { DeleteTodo } from "./DeleteTodo";
import { EditTodo } from "./EditTodo";

export const ShowTodo = ({
  data,
  mutate,
}: {
  data: Todo[] | undefined;
  mutate: KeyedMutator<Todo[]>;
}) => {
  const toggleTodo = async (id: number) => {
    const updatedTodo = await fetch(`${ENDPOINT}/api/todos/${id}/toggle`, {
      method: "PATCH",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((res) => res.json())
      .catch((reason) => console.log(reason));

    mutate(updatedTodo);
  };

  return (
    <List spacing="xs" size="sm" mb={12} center>
      {data &&
        data?.map((todo) => {
          return (
            <List.Item
              key={`todo__${todo.todoId}`}
              icon={
                todo.isCompleted ? (
                  <ThemeIcon
                    onClick={() => toggleTodo(todo.todoId)}
                    color="teal"
                    size={24}
                    radius="xl"
                  >
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                ) : (
                  <ThemeIcon
                    onClick={() => toggleTodo(todo.todoId)}
                    color="gray"
                    size={24}
                    radius="xl"
                  >
                    <CheckCircleFillIcon size={20} />
                  </ThemeIcon>
                )
              }
            >
              <Grid>
                <Grid.Col span={5}>
                  <h2>{todo.title}</h2>
                </Grid.Col>
                <Grid.Col span={2} mt={11}>
                  View Task
                </Grid.Col>
                <Grid.Col span={2} mt={11} ml={5}>
                  <EditTodo todo={todo} mutate={mutate} />
                </Grid.Col>
                <Grid.Col span={2} mt={11} ml={20}>
                  <DeleteTodo todo={todo} mutate={mutate} />
                </Grid.Col>
              </Grid>
            </List.Item>
          );
        })}
    </List>
  );
};

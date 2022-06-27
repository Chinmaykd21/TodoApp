import { List, ThemeIcon } from "@mantine/core";
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
    }).then((res) => res.json());

    mutate(updatedTodo);
  };

  return (
    <List spacing="xs" size="sm" mb={12} center>
      {data &&
        data?.map((todo) => {
          return (
            <List.Item
              onClick={() => toggleTodo(todo.todoId)}
              key={`todo__${todo.todoId}`}
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
              <EditTodo todo={todo} mutate={mutate} />
              {/* TODO: Create these components */}
              <DeleteTodo />
            </List.Item>
          );
        })}
    </List>
  );
};

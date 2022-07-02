import { Button, Group } from "@mantine/core";
import { KeyedMutator } from "swr";
import { Todo } from "../interfaces/todoInterface";
import { ENDPOINT } from "../utilities/utilities";

export const DeleteTodo = ({
  todo,
  mutate,
}: {
  todo: Todo;
  mutate: KeyedMutator<Todo[]>;
}) => {
  const handleDelete = async (todo: Todo) => {
    const deleteTodo = await fetch(
      `${ENDPOINT}/api/todos/${todo.todoId}/delete`,
      {
        method: "DELETE",
        headers: {
          "Content-type": "application/json",
        },
      }
    )
      .then((res) => res.json())
      .catch((reason) => console.log(reason));

    mutate(deleteTodo);
  };

  return (
    <Group position="center">
      <Button mb={10} onClick={() => handleDelete(todo)}>
        Delete
      </Button>
    </Group>
  );
};

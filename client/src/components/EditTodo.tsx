import { Button, Group, Modal, Textarea, TextInput } from "@mantine/core";
import { useForm } from "@mantine/hooks";
import { useState } from "react";
import { ENDPOINT } from "../utilities/utilities";
import { KeyedMutator } from "swr";
import { Todo } from "../interfaces/todoInterface";

export const EditTodo = ({
  todo,
  mutate,
}: {
  todo: Todo;
  mutate: KeyedMutator<Todo[]>;
}) => {
  // Initialize a form with pre-populated values
  const form = useForm({
    initialValues: {
      title: todo.title,
      body: todo.body,
    },
  });

  const [open, setOpen] = useState(false);

  // handlesubmit values
  const handleSubmit = async (values: { title: string; body: string }) => {
    const editedTodo = await fetch(`${ENDPOINT}/api/todos/edit`, {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify(values),
    }).then((res) => res.json());

    mutate(editedTodo);

    form.reset(); //  TODO: Test out whether to do this or not?
    setOpen(false);
  };

  return (
    <>
      <Modal opened={open} onClose={() => setOpen(false)} title="Create Todo">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            required
            mb={12}
            label="Title"
            placeholder="What do you want to do?"
            {...form.getInputProps("title")}
          />
          <Textarea
            required
            mb={12}
            label="Body"
            placeholder="Add a detailed description of the TODO"
            {...form.getInputProps("body")}
          />
          <Group position="right" mt="md">
            <Button fullWidth type="submit">
              Submit
            </Button>
          </Group>
        </form>
      </Modal>
      <Group position="center">
        <Button mb={12} onClick={() => setOpen(true)}>
          Edit
        </Button>
      </Group>
    </>
  );
};

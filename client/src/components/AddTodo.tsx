import { Button, Group, Modal, Textarea, TextInput } from "@mantine/core";
import { useForm } from "@mantine/hooks";
import { useState } from "react";
import { KeyedMutator } from "swr";
import { Todo } from "../interfaces/todoInterface";
import { ENDPOINT } from "../utilities/utilities";

export const AddTodo = ({ mutate }: { mutate: KeyedMutator<Todo[]> }) => {
  // Initialize an empty form
  const form = useForm({
    initialValues: {
      title: "",
      body: "",
    },
  });

  // handle state of the modal
  const [open, setOpen] = useState(false);

  // handlesubmit values
  const handleSubmit = async (values: { title: string; body: string }) => {
    const addedTodo = await fetch(`${ENDPOINT}/api/todos`, {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify(values),
    }).then((res) => res.json());

    mutate(addedTodo);

    form.reset();
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
              Add Todo
            </Button>
          </Group>
        </form>
      </Modal>
      <Group position="center">
        <Button fullWidth mb={12} onClick={() => setOpen(true)}>
          Add TODO
        </Button>
      </Group>
    </>
  );
};

export const ENDPOINT: string = "http://localhost:4000";

export const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((res) => {
    return res.json();
  });

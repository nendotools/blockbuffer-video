export const useFetch = async <T>(url: string, fetchOptions?: any) => {
  console.log(process.env);
  const response = await fetch(
    `/api${url}`,
  );

  if (!response.ok) {
    throw new Error(response.statusText);
  }
  return (await response.json()) as T;
};

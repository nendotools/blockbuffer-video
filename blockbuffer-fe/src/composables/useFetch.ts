export const useFetch = async <T>(url: string, fetchOptions?: any) => {
  const response = await fetch(
    `/api${url}`,
    fetchOptions
  );

  if (!response.ok) {
    throw new Error(response.statusText);
  }
  return (await response.json()) as T;
};

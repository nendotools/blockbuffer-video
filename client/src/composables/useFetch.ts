export const useFetch = async <T>(url: string, fetchOptions?: any) => {
  const response = await fetch(
    `/api${url}`,
    {
      method: fetchOptions?.method || "GET",
      body:
        fetchOptions?.data ? fetchOptions.data :
          fetchOptions?.body ? JSON.stringify(fetchOptions.body) : null,
    },
  );

  if (!response.ok) {
    throw new Error(response.statusText);
  }
  return (await response.json()) as T;
};

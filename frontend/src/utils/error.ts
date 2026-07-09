export function getErrorMessage(e: unknown): string {
  const data = (e as any)?.response?.data;
  if (!data) return "An error occurred";

  if (Array.isArray(data.errors) && data.errors.length > 0) {
    return data.errors.join(", ");
  }
  if (typeof data.error === "string") {
    return data.error;
  }
  if (typeof data.message === "string") {
    return data.message;
  }
  return "An error occurred";
}

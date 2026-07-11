function parseErrors(e: unknown): string[] {
  const data = (e as any)?.response?.data;
  if (!data) return [];

  if (Array.isArray(data.errors) && data.errors.length > 0) {
    return data.errors;
  }
  if (data.errors && typeof data.errors === "object") {
    const msgs: string[] = [];
    for (const key of Object.keys(data.errors)) {
      if (Array.isArray(data.errors[key])) {
        msgs.push(...data.errors[key]);
      }
    }
    if (msgs.length > 0) return msgs;
  }
  if (typeof data.error === "string") return [data.error];
  if (typeof data.message === "string") return [data.message];

  return [];
}

export function getErrorMessage(e: unknown): string {
  const msgs = parseErrors(e);
  return msgs.length > 0 ? msgs.join(", ") : "An error occurred";
}

export function getFirstErrorMessage(e: unknown): string {
  const msgs = parseErrors(e);
  return msgs.length > 0 ? msgs[0] : "An error occurred";
}

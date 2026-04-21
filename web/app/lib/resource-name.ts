// Helpers that mirror internal/platform/resourcename on the Go side.

export function formatUser(customId: string): string {
  return `users/${customId}`;
}

export function parseUser(name: string): string {
  const parts = name.split("/");
  if (parts.length !== 2 || parts[0] !== "users" || !parts[1]) {
    throw new Error(`invalid user resource name: ${name}`);
  }
  return parts[1];
}

export function formatCall(id: string): string {
  return `calls/${id}`;
}

export function parseCall(name: string): string {
  const parts = name.split("/");
  if (parts.length !== 2 || parts[0] !== "calls" || !parts[1]) {
    throw new Error(`invalid call resource name: ${name}`);
  }
  return parts[1];
}

export function parseCallParticipant(name: string): {
  callId: string;
  identity: string;
} {
  const parts = name.split("/");
  if (
    parts.length !== 4 ||
    parts[0] !== "calls" ||
    parts[2] !== "participants" ||
    !parts[1] ||
    !parts[3]
  ) {
    throw new Error(`invalid call participant resource name: ${name}`);
  }
  return { callId: parts[1], identity: parts[3] };
}

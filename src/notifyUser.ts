import { UserToNotify } from "./findUsersToNotify";
import { render as renderTemplate } from 'mustache';

enum NotificationError {
  AddressDoesNotExist = 404,
  ServerUnavailable = 503,
}

type NotificationResult = {
  error?: NotificationError
}


export async function notifyUser(user: UserToNotify): Promise<NotificationResult> {
  const { error } = await dummySendNewsletter(user);

  if (!error) {
    return {};
  }

  switch (error) {
    case NotificationError.AddressDoesNotExist: {
      throw new NonRetryableError("Address does not exist");
    } case NotificationError.ServerUnavailable: {
      throw new RetryableError("Server unavailable");
    }
    default:
      throw new Error("Unhandled error occurred");
  }
}

async function dummySendNewsletter(user: UserToNotify) {
  const shouldSimulateErrorCondition = Math.random() > 0.8;

  if (shouldSimulateErrorCondition) {
    return { error: Math.ceil(Math.random() * 10) % 2 == 0 ? NotificationError.ServerUnavailable : NotificationError.AddressDoesNotExist }
  }

  renderTemplate("<h1>Hello {{username}}!</h1>", { username: user.username });

  await new Promise((resolve, _) => setTimeout(resolve, 300));

  return {};
}

class NonRetryableError extends Error {
  constructor(msg: string) {
    super(msg);
  }
}

class RetryableError extends Error {
  constructor(msg: string) {
    super(msg);
    this.message = msg;
    this.name = "RetryableError";
  }
}
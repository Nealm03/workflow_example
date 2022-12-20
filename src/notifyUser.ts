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
  return await dummySendNewsletter(user);
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



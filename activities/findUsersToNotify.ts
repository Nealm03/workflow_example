type NotifyFilter = {
  olderThan: number;

}

export type UserToNotify = {
  username: string;
  email: string;
  age: number;
}

export async function findUsersToNotify(filter: NotifyFilter): Promise<UserToNotify[]> {
  return dummyUsers.filter(user => user.age > filter.olderThan);
}


const dummyUsers: Array<UserToNotify> = [
  {
    username: "tom.404",
    email: "tom@gmail.com",
    age: 30,
  },
  {
    username: "sally.429",
    email: "sally@gmail.com",
    age: 24
  },
  {
    username: "john.503",
    email: "john@gmail.com",
    age: 56
  }
];
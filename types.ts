export interface UserHours {
  notes: string;
  hours: string;
  release: {
    name: string;
    project: {
      name: string;
    };
  };
  hoursTag: {
    name: string;
  };
}

export interface ILoginPrompt {
  username: string;
  password: string;
}
export interface BaseDto {
  _id: string;
  name: string;
}

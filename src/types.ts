export interface CliOptions {
  r: number;
  relative: number;
  d: string;
  date: string;
  m: number;
  month: number;
  h: boolean;
  help: boolean;
}

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

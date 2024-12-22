export type UploadBody = {
  notes: string;
  hours: string;
  date: string;
  releaseId: string;
  hoursTagId: string;
  userId: string;
};

export interface UserHours {
  notes: string;
  date: string;
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

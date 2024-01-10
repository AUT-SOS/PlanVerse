import { ajax } from "rxjs/ajax";
import { END_POINT } from "../utils/consts";

export const API = {
  login(email: string, password: string) {
    return ajax.post(END_POINT, {
      email,
      password,
    });
  },
  signup(email: string, password: string, username: string){
    return ajax.post(END_POINT, {
      email,
      password,
      username
    })
  }
};

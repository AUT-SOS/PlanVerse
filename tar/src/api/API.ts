import { ajax } from "rxjs/ajax";
import { END_POINT, API_HEADERS } from "../utils/consts";

export const API = {
  login(email: string, password: string) {
    return ajax.post(
      `${END_POINT}/login`,
      {
        email,
        password,
      },
      API_HEADERS
    );
  },
  signup(email: string, password: string, username: string) {
    return ajax.post(
      `${END_POINT}/register`,
      {
        email: email,
        password: password,
        username: username,
      },
      API_HEADERS
    );
  },
  otpVerify(verify: string) {
    const at = getAccessToken() ?? "";
    return ajax.post(
      `${END_POINT}/verify`,
      {
        otp: verify,
      },
      { ...API_HEADERS, "Authorization" : getAccessToken()}
    );
  },
};

const getAccessToken = () => getCookie("access_token");

const getCookie = (cname : string) => {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');
  for(let i = 0; i <ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

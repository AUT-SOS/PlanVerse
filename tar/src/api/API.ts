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
    return ajax.post(
      `${END_POINT}/verify`,
      {
        otp: verify,
      },
      { ...API_HEADERS, Authorization: getAccessToken() }
    );
  },
  getMyId() {
    return ajax.get(`${END_POINT}/get-my-user`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
    });
  },
  refresh() {
    return ajax.post(
      `${END_POINT}/refresh`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  resendEmail() {
    return ajax.post(
      `${END_POINT}/resend-email`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  getUser(id: string) {
    return ajax.get(`${END_POINT}/get-user/${id}`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
    });
  },
  createProject(title: string, description: string, picture: string) {
    return ajax.post(
      `${END_POINT}/create-project`,
      {
        title,
        picture,
        description,
      },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  editProject(id: string, title: string, description: string, picture: string) {
    return ajax.post(
      `${END_POINT}/edit-project/${id}`,
      {
        title,
        picture,
        description,
      },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  getMyProjects() {
    return ajax.get(`${END_POINT}/list-project`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
    });
  },
  getFullProject(id: string) {
    return ajax.get(`${END_POINT}/get-project/${id}`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
    });
  },
  getProjectMembers(id: string) {
    return ajax.get(`${END_POINT}/get-project-members/${id}`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
    });
  },
  promote(projectId: string, userId: string) {
    return ajax.post(
      `${END_POINT}/promote/${projectId}/${userId}`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  demote(projectId: string, userId: string) {
    return ajax.post(
      `${END_POINT}/demote/${projectId}/${userId}`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  editUser(
    username: string,
    password: string,
    email: string,
    profile_pic: string
  ) {
    return ajax.post(
      `${END_POINT}/edit-profile`,
      {
        username,
        password,
        email,
        profile_pic,
      },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  showProject(join_link: string) {
    return ajax.post(
      `${END_POINT}/show-project`,
      {
        join_link,
      },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  joinProject(id: string) {
    return ajax.post(
      `${END_POINT}/join-project/${id}`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  shareLink(id: string, emails: string[]) {
    return ajax.post(
      `${END_POINT}/share-link/${id}`,
      {emails},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  }
};

const getAccessToken = () => getCookie("access_token");

const getCookie = (cname: string) => {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == " ") {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
};

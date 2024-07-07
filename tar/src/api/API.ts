import { ajax } from "rxjs/ajax";
import { END_POINT, API_HEADERS } from "../utils/consts";

export const API = {
  connectWS() {
    return ajax.get(`${END_POINT}/create-ws`, {
      ...API_HEADERS,
      Authorization: getAccessToken(),
      "Connection": "Upgrade",
      "Upgrade": "websocket",
    });
  },
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
  leaveProject(id: string) {
    return ajax.post(
      `${END_POINT}/leave-project/${id}`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
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
      { emails },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  deleteUser(id: string) {
    return ajax.post(
      `${END_POINT}/delete-account`,
      { user_id: id },
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  deleteProject(id: string) {
    return ajax.post(
      `${END_POINT}/delete-project/${id}`,
      {},
      {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      }
    );
  },
  Board: {
    getStates(id: string) {
      return ajax.get(`${END_POINT}/list-state/${id}`, {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      });
    },
    createState(
      id: string,
      title: string,
      back_ground_color: string,
      admin_access: boolean
    ) {
      return ajax.post(
        `${END_POINT}/create-state/${id}`,
        {
          title,
          back_ground_color,
          admin_access,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    getState(projId: string, stateId: string) {
      return ajax.get(`${END_POINT}/get-state/${projId}/${stateId}`, {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      });
    },
    createTask(
      id: string,
      state_id: string,
      title: string,
      back_ground_color: string,
      description: string,
      index: number,
      deadline?: string,
      estimated_time?: number,
      priority?: number
    ) {
      return ajax.post(
        `${END_POINT}/create-task/${id}/${state_id}`,
        {
          title,
          back_ground_color,
          description,
          index,
          deadline,
          estimated_time,
          priority,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    editState(
      id: string,
      state_id: string,
      title: string,
      back_ground_color: string,
      admin_access: boolean
    ) {
      return ajax.post(
        `${END_POINT}/edit-state/${id}/${state_id}`,
        {
          title,
          back_ground_color,
          admin_access,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    editTask(
      id: string,
      task_id: string,
      title: string,
      back_ground_color: string,
      description: string,
      index: number,
      deadline?: string,
      estimated_time?: number,
      priority?: number
    ) {
      return ajax.post(
        `${END_POINT}/edit-task/${id}/${task_id}`,
        {
          title,
          back_ground_color,
          description,
          index,
          deadline,
          estimated_time,
          priority,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    deleteState(id: string, state_id: string) {
      return ajax.post(
        `${END_POINT}/delete-state/${id}/${state_id}`,
        {},
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    changeState(id: string, task_id: string, state_id: string) {
      return ajax.post(
        `${END_POINT}/change-state/${id}/${task_id}`,
        {
          state_id,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    getTask(id: string) {
      return ajax.get(`${END_POINT}/get-task/${id}`, {
        ...API_HEADERS,
        Authorization: getAccessToken(),
      });
    },
    deleteTask(id: string, task_id: string) {
      return ajax.post(
        `${END_POINT}/delete-task/${id}/${task_id}`,
        {},
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    addAssign(id: string, task_id: string, performer_id: string) {
      return ajax.post(
        `${END_POINT}/add-performer/${id}/${task_id}`,
        {
          performer_id,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
    removeAssign(id: string, task_id: string, performer_id: string) {
      return ajax.post(
        `${END_POINT}/remove-performer/${id}/${task_id}`,
        {
          performer_id,
        },
        {
          ...API_HEADERS,
          Authorization: getAccessToken(),
        }
      );
    },
  },
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

import { END_POINT } from "../utils/consts";
import { Project, TaskColumn } from "../utils/types";

const getCulmons = (groupId: string) => {
  let group = new Promise(async (resolve, reject) => {
    await fetch(END_POINT)
      .then((response) => response.body)
      .then((body) => resolve(JSON.parse(JSON.stringify(body)) as TaskColumn[]))
      .catch((error) => reject(error));
  });
  return group;
};




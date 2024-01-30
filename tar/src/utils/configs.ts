import { THEMES } from "./consts";

export const setRootColors = (theme: "dark" | "light") => {
  const selectedTheme = theme === "dark" ? THEMES.dark : THEMES.light;
  for (let key of Object.keys(selectedTheme)) {
    setColor(key, selectedTheme[key]);
  }
};

const setColor = (key: string, val: string) => {
  document.documentElement.style.setProperty(key, val);
};

export const secToString = (secs: number) => {
  const min = Math.floor(secs / 60);
  const sec = secs - min * 60;
  return `${min < 10 ? 0 : ""}${min}:${sec < 10 ? 0 : ""}${sec}`

}

export const navigate = (dest: string) => {
  navigate(dest);
}

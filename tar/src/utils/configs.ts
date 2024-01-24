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

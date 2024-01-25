export const ICON_COLOR = "var(--color-neutrals-n-500)";
export const RANDOM_IMAGE = "https://picsum.photos/100/100";
export const END_POINT = "/api";
export const DOC_ADDRESS =
  "https://docs.google.com/document/d/1PxMYqd9ZbLHtC88sYyzR_MmyfN9CBZcgLuvFWx5X6pE/edit?usp=sharing";
export const API_HEADERS = {
  "Content-Type": "application/json",
  credentials: "include",
  withCredentials: "true"
};

export const THEMES : {light: Record<string, string>, dark: Record<string, string>} = {
  light: {
    "--color-neutrals-n-100": " #b3bac5",
    "--color-neutrals-n-100-rgb": " 179, 186, 197",

    "--color-neutrals-n-200": " #7a869a",
    "--color-neutrals-n-200-rgb": " 122, 134, 154",

    "--color-neutrals-n-300": " #5e6c84",
    "--color-neutrals-n-300-rgb": " 94, 108, 132",

    "--color-neutrals-n-400": " #42526e",
    "--color-neutrals-n-400-rgb": " 66, 82, 110",

    "--color-neutrals-n-500": " #091e42",
    "--color-neutrals-n-500-rgb": " 9, 30, 66",

    "--color-neutrals-n-600": " #0f0f0f",
    "--color-neutrals-n-600-rgb": " 15, 15, 15",

    "--color-neutrals-on-primary": " var(--white)",
    "--color-neutrals-on-primary-rgb": " 255, 255, 255",

    "--color-neutrals-n-00": " #fcfdfd",
    "--color-neutrals-n-00-rgb": " 252, 253, 253",

    "--color-neutrals-n-10": " #f4f5f7",
    "--color-neutrals-n-10-rgb": " 250, 251, 252",

    "--color-neutrals-n-20": " #f4f5f7",
    "--color-neutrals-n-20-rgb": " 244, 245, 247",

    "--color-neutrals-n-30": " #ebecf0",
    "--color-neutrals-n-30-rgb": " 235, 236, 240",

    "--color-neutrals-n-40": " #dfe1e6",
    "--color-neutrals-n-40-rgb": " 223, 225, 230",

    "--color-neutrals-n-50": " #c1c7d0",
    "--color-neutrals-n-50-rgb": " 193, 199, 208",

    "--color-background": " #f3e0e2",
   " --color-pink": "#fc5c7d"
  },
  dark: {
    "--color-neutrals-n-50": " #c5c5c5",
    "--color-neutrals-n-50-rgb": " 179, 186, 197",

    "--color-neutrals-n-40": " #919191",
    "--color-neutrals-n-40-rgb": " 122, 134, 154",

    "--color-neutrals-n-30": " #838383",
    "--color-neutrals-n-30-rgb": " 94, 108, 132",

    "--color-neutrals-n-20": " #636363",
    "--color-neutrals-n-20-rgb": " 66, 82, 110",

    "--color-neutrals-n-10": " #383838",
    "--color-neutrals-n-10-rgb": " 9, 30, 66",

    "--color-neutrals-n-00": " #0f0f0f",
    "--color-neutrals-n-00-rgb": " 15, 15, 15",

    "--color-neutrals-on-primary": " var(--white)",
    "--color-neutrals-on-primary-rgb": " 255, 255, 255",

    "--color-neutrals-n-600": " #fcfdfd",
    "--color-neutrals-n-600-rgb": " 252, 253, 253",

    "--color-neutrals-n-500": " #fafbfc",
    "--color-neutrals-n-500-rgb": " 250, 251, 252",

    "--color-neutrals-n-400": " #f4f5f7",
    "--color-neutrals-n-400-rgb": " 244, 245, 247",

    "--color-neutrals-n-300": " #ebecf0",
    "--color-neutrals-n-300-rgb": " 235, 236, 240",

    "--color-neutrals-n-200": " #dfe1e6",
    "--color-neutrals-n-200-rgb": " 223, 225, 230",

    "--color-neutrals-n-100": " #c1c7d0",
    "--color-neutrals-n-100-rgb": " 193, 199, 208",
    
    "--color-background": " #1A1B1C",
    " --color-pink": "#1b0a0d"

  },
};

export const validateEmail = (email?: string) => {
  return email
    ? email.match(
        /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      )
    : true;
};

export const validatePassword = (password?: string) => {
  return password
    ? password.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/)
    : true;
};

export const validateUsername = (username? : string) => {
  return username
  ? username.length >= 3
  : true;
}


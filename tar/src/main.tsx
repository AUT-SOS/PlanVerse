import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import "./palette.scss";
import { Provider } from "react-redux";
import { store } from "./redux/store.ts";
import { ToastContainer, toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.min.css';

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <ToastContainer/>
      <App />
    </Provider>
  </React.StrictMode>
);

export const showSuccessToastMessage = (message : string) => {
  toast.success(message, {
      position: toast.POSITION.BOTTOM_RIGHT,
      draggable : true,
      autoClose : 3000
  });
};

export const showFailToastMessage = (message : string) => {
  toast.error(message, {
      position: toast.POSITION.BOTTOM_RIGHT,
      draggable : true,
      autoClose : 3000
  });
};

export const showInfoToastMessage = (message : string) => {
  toast.info(message, {
      position: toast.POSITION.BOTTOM_RIGHT,
      draggable : true,
      autoClose : 3000
  });
};


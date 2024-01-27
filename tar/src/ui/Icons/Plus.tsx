import React from "react";
import { IconComponent } from "../../utils/types";

export const Plus: IconComponent = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={props.size}
      height={props.size}
      viewBox="0 0 10 10"
      fill="none"
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M5 0.578125C5.41421 0.578125 5.75 0.913911 5.75 1.32812V8.65448C5.75 9.06869 5.41421 9.40448 5 9.40448C4.58579 9.40448 4.25 9.06869 4.25 8.65448V1.32812C4.25 0.913911 4.58579 0.578125 5 0.578125Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M0.582031 4.99219C0.582031 4.57797 0.917818 4.24219 1.33203 4.24219H4.9987H8.66537C9.07958 4.24219 9.41537 4.57797 9.41537 4.99219C9.41537 5.4064 9.07958 5.74219 8.66537 5.74219H1.33203C0.917818 5.74219 0.582031 5.4064 0.582031 4.99219Z"
        fill={props.color}
      />
    </svg>
  );
};

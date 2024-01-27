import React from "react";
import { IconComponent } from "../../utils/types";

export const Logout: IconComponent = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={props.size}
      height={props.size}
      viewBox="0 0 21 20"
      fill="none"
      {...props}
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M9.319 20H4.433C1.989 20 0 18.011 0 15.565V4.436C0 1.99 1.989 0 4.433 0H9.308C11.754 0 13.744 1.99 13.744 4.436V5.368C13.744 5.782 13.408 6.118 12.994 6.118C12.58 6.118 12.244 5.782 12.244 5.368V4.436C12.244 2.816 10.927 1.5 9.308 1.5H4.433C2.816 1.5 1.5 2.816 1.5 4.436V15.565C1.5 17.184 2.816 18.5 4.433 18.5H9.319C10.931 18.5 12.244 17.188 12.244 15.576V14.633C12.244 14.219 12.58 13.883 12.994 13.883C13.408 13.883 13.744 14.219 13.744 14.633V15.576C13.744 18.016 11.758 20 9.319 20Z"
        fill={props.color}
      />
      <mask
        id="mask0_13497_14088"
        maskUnits="userSpaceOnUse"
        x="6"
        y="9"
        width="15"
        height="2"
      >
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M6.99609 9.25H20.537V10.75H6.99609V9.25Z"
          fill="white"
        />
      </mask>
      <g mask="url(#mask0_13497_14088)">
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M19.7871 10.75H7.74609C7.33209 10.75 6.99609 10.414 6.99609 10C6.99609 9.586 7.33209 9.25 7.74609 9.25H19.7871C20.2011 9.25 20.5371 9.586 20.5371 10C20.5371 10.414 20.2011 10.75 19.7871 10.75Z"
          fill={props.color}
        />
      </g>
      <mask
        id="mask1_13497_14088"
        maskUnits="userSpaceOnUse"
        x="16"
        y="6"
        width="5"
        height="8"
      >
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M16.1094 6.33594H20.5367V13.6667H16.1094V6.33594Z"
          fill="white"
        />
      </mask>
      <g mask="url(#mask1_13497_14088)">
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M16.8591 13.6667C16.6671 13.6667 16.4741 13.5937 16.3281 13.4457C16.0361 13.1517 16.0371 12.6777 16.3301 12.3857L18.7241 10.0007L16.3301 7.61669C16.0371 7.32469 16.0351 6.85069 16.3281 6.55669C16.6201 6.26269 17.0941 6.26269 17.3881 6.55469L20.3161 9.46969C20.4581 9.60969 20.5371 9.80169 20.5371 10.0007C20.5371 10.1997 20.4581 10.3917 20.3161 10.5317L17.3881 13.4477C17.2421 13.5937 17.0501 13.6667 16.8591 13.6667Z"
          fill={props.color}
        />
      </g>
    </svg>
  );
};

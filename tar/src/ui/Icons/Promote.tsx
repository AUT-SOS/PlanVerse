import React from "react";
import { IconComponent } from "../../utils/types";

export const Promote: IconComponent = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={props.size}
      height={props.size}
      viewBox="0 0 20 20"
      fill="none"
      {...props}
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M18.4771 18.3132C17.4701 19.3947 16.0343 20 14.334 20H5.665C3.96469 20 2.52906 19.3947 1.52228 18.3131C0.521154 17.2377 0 15.7523 0 14.084V5.916C0 4.24667 0.524275 2.76143 1.52641 1.68646C2.53397 0.605656 3.96915 0 5.665 0H14.334C16.0343 0 17.4701 0.605297 18.4771 1.68678C19.4785 2.76222 20 4.24761 20 5.916V14.084C20 15.7524 19.4785 17.2378 18.4771 18.3132ZM17.3794 17.291C18.0785 16.5402 18.5 15.4426 18.5 14.084V5.916C18.5 4.5574 18.0785 3.45978 17.3794 2.70897C16.6859 1.9642 15.6637 1.5 14.334 1.5H5.665C4.34085 1.5 3.31853 1.96385 2.62359 2.70929C1.92323 3.46057 1.5 4.55833 1.5 5.916V14.084C1.5 15.4427 1.92135 16.5403 2.62022 17.2911C3.31344 18.0358 4.3353 18.5 5.665 18.5H14.334C15.6637 18.5 16.6859 18.0358 17.3794 17.291Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M10 14.8359C9.58579 14.8359 9.25 14.5002 9.25 14.0859L9.25 5.91394C9.25 5.49972 9.58579 5.16394 10 5.16394C10.4142 5.16394 10.75 5.49972 10.75 5.91394L10.75 14.0859C10.75 14.5002 10.4142 14.8359 10 14.8359Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M14.2753 10.2072C13.9818 10.4995 13.5069 10.4985 13.2146 10.205L9.99809 6.97471L6.78155 10.205C6.48928 10.4985 6.01441 10.4995 5.72089 10.2072C5.42738 9.91497 5.42637 9.4401 5.71864 9.14658L9.46664 5.38258C9.60738 5.24124 9.79863 5.16178 9.99809 5.16178C10.1976 5.16178 10.3888 5.24124 10.5296 5.38258L14.2776 9.14658C14.5698 9.4401 14.5688 9.91497 14.2753 10.2072Z"
        fill={props.color}
      />
    </svg>
  );
};

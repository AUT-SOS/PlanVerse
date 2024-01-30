import { IconComponent } from "../../utils/types";

export const ColorPicker: IconComponent = (props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={props.size}
      height={props.size}
      viewBox="0 0 24 24"
      fill="none"
      {...props}
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M6.99971 10.75C6.86225 10.75 6.75 10.862 6.75 11C6.75 11.138 6.86225 11.25 6.99971 11.25C7.13774 11.25 7.25 11.1375 7.25 11C7.25 10.8625 7.13774 10.75 6.99971 10.75ZM5.25 11C5.25 10.0343 6.03311 9.25 6.99971 9.25C7.96577 9.25 8.75 10.0337 8.75 11C8.75 11.9663 7.96577 12.75 6.99971 12.75C6.03311 12.75 5.25 11.9657 5.25 11Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M10.9997 6.75C10.8622 6.75 10.75 6.86198 10.75 7C10.75 7.13802 10.8622 7.25 10.9997 7.25C11.1377 7.25 11.25 7.13747 11.25 7C11.25 6.86253 11.1377 6.75 10.9997 6.75ZM9.25 7C9.25 6.03426 10.0331 5.25 10.9997 5.25C11.9658 5.25 12.75 6.0337 12.75 7C12.75 7.9663 11.9658 8.75 10.9997 8.75C10.0331 8.75 9.25 7.96574 9.25 7Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M15.9997 8.75C15.8622 8.75 15.75 8.86198 15.75 9C15.75 9.13802 15.8622 9.25 15.9997 9.25C16.1377 9.25 16.25 9.13747 16.25 9C16.25 8.86253 16.1377 8.75 15.9997 8.75ZM14.25 9C14.25 8.03426 15.0331 7.25 15.9997 7.25C16.9658 7.25 17.75 8.0337 17.75 9C17.75 9.9663 16.9658 10.75 15.9997 10.75C15.0331 10.75 14.25 9.96574 14.25 9Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M12.1214 20.2236L12.111 20.2119C11.6126 19.6498 11.3167 18.9146 11.3167 18.1111C11.3167 16.3604 12.7327 14.9444 14.4833 14.9444H16.4444C18.6827 14.9444 20.5 13.1271 20.5 10.8889C20.5 6.96765 16.8651 3.5 12 3.5C7.30619 3.5 3.5 7.30619 3.5 12C3.5 16.6938 7.30619 20.5 12 20.5C12.0938 20.5 12.1667 20.4271 12.1667 20.3333C12.1667 20.3043 12.1616 20.2857 12.1569 20.2738C12.152 20.2613 12.1443 20.2483 12.132 20.235L12.1214 20.2236ZM2 12C2 6.47776 6.47776 2 12 2C17.5222 2 22 5.97776 22 10.8889C22 13.9555 19.5111 16.4444 16.4444 16.4444H14.4833C13.5611 16.4444 12.8167 17.1889 12.8167 18.1111C12.8167 18.5389 12.9722 18.9222 13.2333 19.2167C13.5056 19.5111 13.6667 19.9 13.6667 20.3333C13.6667 21.2556 12.9222 22 12 22C6.47776 22 2 17.5222 2 12Z"
        fill={props.color}
      />
    </svg>
  );
};

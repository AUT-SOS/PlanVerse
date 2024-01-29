import { IconComponent } from "../../utils/types";

export const BackIcon: IconComponent = (props) => {
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
        d="M3.5 12.0024C3.5 11.5744 3.84596 11.2274 4.27272 11.2274H19.7271C20.1538 11.2274 20.4998 11.5744 20.4998 12.0024C20.4998 12.4304 20.1538 12.7773 19.7271 12.7773H4.27272C3.84596 12.7773 3.5 12.4304 3.5 12.0024Z"
        fill={props.color}
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M12.9439 18.7719C12.6428 18.4686 12.6438 17.9779 12.9462 17.6759L18.6296 11.9995L12.9463 6.32413C12.6438 6.02214 12.6428 5.53146 12.9439 5.22816C13.245 4.92487 13.7343 4.92381 14.0367 5.2258L20.2699 11.4503C20.4156 11.5957 20.4974 11.7933 20.4974 11.9995C20.4975 12.2056 20.4156 12.4032 20.27 12.5486L14.0367 18.7742C13.7343 19.0762 13.2451 19.0752 12.9439 18.7719Z"
        fill={props.color}
      />
    </svg>
  );
};

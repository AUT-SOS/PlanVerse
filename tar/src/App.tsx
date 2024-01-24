import { useDispatch } from "react-redux";
import { RoutesFC } from "./Routes";
import { setRootColors } from "./utils/configs";
import { ConfigActions } from "./redux/slices/configs.slice";
function App() {
  const dispatch = useDispatch();
  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    setRootColors("dark");
  } else {
    setRootColors("light");
  }
  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (event) => {
      const newColorScheme = event.matches ? "dark" : "light";
      setRootColors(newColorScheme);
    });
  window.addEventListener("resize", () => {
    dispatch(ConfigActions.setBreakpoint(window.innerWidth));
  });

  return <RoutesFC />;
}

export default App;

import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Join } from "./components/Join/Join";
import { Board } from "./components/Boards/Board";
import { Auth } from "./components/Auth/Auth";
import { Home } from "./components/Home/Home";

type Props = React.HTMLProps<HTMLDivElement> 

export const RoutesFC: React.FC<Props> = (props) => {
  return (
    <Router>
      <Routes>
        <Route path="/*" element={<Auth/>} />
        <Route path="/projects/:id" element={<Board/> } />
        <Route path="/project/*" element={<Join />} />
        <Route path="/login" element={<Auth />} />
        <Route path="/signup" element={<Auth />} />
        <Route path="/home" element={<Home/>}/>
      </Routes>
      {props.children}  
    </Router>
  );
};

import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Join } from "./components/Join/Join";
import { Board } from "./components/Boards/Board";

type Props = React.HTMLProps<HTMLDivElement> 

export const RoutesFC: React.FC<Props> = (props) => {
  return (
    <Router>
      <Routes>
        <Route path="/home" />
        <Route path="/boards/:id" element={<Board/> } />
        <Route path="/join/:id" element={<Join />} />
      </Routes>
      {props.children}  
    </Router>
  );
};

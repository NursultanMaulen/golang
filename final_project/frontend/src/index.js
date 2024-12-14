import React from "react";
import ReactDOM from "react-dom";

import App from "./App";
import { BrowserRouter } from "react-router-dom";
import Pagestructure from "./Pagestructure";
import ExplorepageContext from "./Context/ExplorepageContext";

ReactDOM.render(
  <React.StrictMode>
    <ExplorepageContext>
      <BrowserRouter>
        <Pagestructure />
      </BrowserRouter>
    </ExplorepageContext>
  </React.StrictMode>,
  document.getElementById("root")
);

import React from "react";
import { Route, Routes } from "react-router-dom";
import "./App.css";
import Authentication from "./Components/Authentication/Authentication";
import Pagenotfound from "./Components/PageNotFound/Pagenotfound";
import Accountpage from "./Pages/Accountpage/Accountpage";
import Explorepage from "./Pages/Explorepage/Explorepage";
import Homepage from "./Pages/Homepage/Homepage";
import Cartpage from "./Pages/Cartpage/Cartpage";
import Loginpage from "./Pages/Loginpage/Loginpage";
import Signuppage from "./Pages/Singuppage/Signuppage";
import EditVideoPage from "./Pages/Videopage/Videopage";
import SearchResultsPage from "./Pages/SearchResultspage/SearchResultspage";

function App() {
  return (
    <div>
      <Routes>
        <Route exact path="/" element={<Homepage />} />
        <Route exact path="/explore" element={<Explorepage />} />
        <Route
          exact
          path="/likes"
          element={
            <Authentication>
              <Cartpage />
            </Authentication>
          }
        />
        <Route exact path="/search" element={<SearchResultsPage />} />
        <Route exact path="/login" element={<Loginpage />} />
        <Route exact path="/signup" element={<Signuppage />} />
        <Route
          exact
          path="/accounts"
          element={
            <Authentication>
              <Accountpage />
            </Authentication>
          }
        />
        <Route
          exact
          path="/video/:id/edit"
          element={
            <Authentication>
              <EditVideoPage />
            </Authentication>
          }
        />
        <Route exact path="*" element={<Pagenotfound />} />
      </Routes>
    </div>
  );
}

export default App;

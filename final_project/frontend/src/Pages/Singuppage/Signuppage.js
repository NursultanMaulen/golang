import React from "react";
import {
  Footer,
  Header,
  Sidebar,
  Signupinputs,
} from "../../Components/IndexAllComponents";
import { Layout } from "antd";

function Signuppage() {
  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout
        style={{
          backgroundColor: "#fff",
          height: "100%",
          marginBottom: "100px",
        }}
      >
        <Sidebar />
        <Signupinputs />
        <Footer />
      </Layout>
    </Layout>
  );
}

export default Signuppage;

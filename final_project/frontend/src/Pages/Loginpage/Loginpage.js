import React from "react";
import {
  Footer,
  Header,
  Logininputs,
  Sidebar,
} from "../../Components/IndexAllComponents";
import { Layout } from "antd";

function Loginpage() {
  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout>
        <Sidebar />
        <Logininputs />
        <Footer />
      </Layout>
    </Layout>
  );
}

export default Loginpage;

import React from "react";
import { Layout, Row, Col } from "antd";
import {
  Footer,
  Header,
  Hero,
  Sidebar,
} from "../../Components/IndexAllComponents";

const { Content } = Layout;

function Homepage() {
  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout>
        <Sidebar />
        <Hero />
        <Footer />
      </Layout>
    </Layout>
  );
}

export default Homepage;
